/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	oldmonkv1 "github.com/remotegarage/oldmonk/api/v1"
	x "github.com/remotegarage/oldmonk/pkg"
	"github.com/remotegarage/oldmonk/pkg/scalex"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)

// QueueAutoScalerReconciler reconciles a QueueAutoScaler object
type QueueAutoScalerReconciler struct {
	client.Client
	Log    logr.Logger
	mgr    ctrl.Manager
	Scheme *runtime.Scheme
}


var log = logf.Log.WithName("controller_queueAutoScaler")

const queueAutoScalerFinalizer = "finalizer.oldmonk.evalsocket.in"

// +kubebuilder:rbac:groups=oldmonk.remotegarage.club,resources=queueautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oldmonk.remotegarage.club,resources=queueautoscalers/status,verbs=get;update;patch

func (r *QueueAutoScalerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	reqLogger := r.Log.WithValues("queueautoscaler", req.NamespacedName)
	reqLogger.Info("Reconciling QueueAutoScaler")

	// Fetch the QueueAutoScaler instance
	queueAutoScaler := &oldmonkv1.QueueAutoScaler{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, queueAutoScaler)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("QueueAutoScaler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get QueueAutoScaler.")
		return ctrl.Result{}, err
	}

	if !queueAutoScaler.Spec.Autopilot {
		return ctrl.Result{}, nil
	}

	// Check if the AutoScaler  instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isAutoScalerMarkedToBeDeleted := queueAutoScaler.GetDeletionTimestamp() != nil
	if isAutoScalerMarkedToBeDeleted {
		if x.Contains(queueAutoScaler.GetFinalizers(), queueAutoScalerFinalizer) {
			// Run finalization logic for queueAutoScalerFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeAutoScaler(queueAutoScaler); err != nil {
				return ctrl.Result{}, err
			}

			// Remove queueAutoScalerFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			queueAutoScaler.SetFinalizers(x.Remove(queueAutoScaler.GetFinalizers(), queueAutoScalerFinalizer))
			err := r.Client.Update(context.TODO(), queueAutoScaler)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	} else {
		// Check if the Deployment already exists, if not create a new one
		deployment := &appsv1.Deployment{}
		dep := r.deploymentForDeployment(queueAutoScaler)
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: queueAutoScaler.Spec.Deployment, Namespace: queueAutoScaler.Namespace}, deployment)
		if err != nil && errors.IsNotFound(err) {
			// Define a new Deployment
			reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			err = r.Client.Create(context.TODO(), dep)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			// NOTE: that the requeue is made with the purpose to provide the deployment object for the next step to ensure the deployment size is the same as the spec.
			// Also, you could GET the deployment object again instead of requeue if you wish. See more over it here: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Deployment.")
			return ctrl.Result{}, err
		}

		reqLogger.Info("Updating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Update(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to update new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Ensure the deployment size is the same as the spec
		// TO-Do perform queue logic and then scale
		scaler := scalex.NewScalex(r.mgr, time.Second*60)
		deployment, delta, err := scaler.ExecuteScale(context.TODO(), queueAutoScaler)

		if err != nil {
			reqLogger.Error(err, "unable to perform desired replica, will retry:")
			return ctrl.Result{}, err
		}
		log.WithValues("Request.delta", delta, "Request.desired", *deployment.Spec.Replicas, "Request.available", deployment.Status.AvailableReplicas).Info("Updated deployment")

	}

	// Add finalizer for this CR
	if !x.Contains(queueAutoScaler.GetFinalizers(), queueAutoScalerFinalizer) {
		if err := r.addFinalizer(queueAutoScaler); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *QueueAutoScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.mgr = mgr
	return ctrl.NewControllerManagedBy(mgr).
		For(&oldmonkv1.QueueAutoScaler{}).
		Complete(r)
}

// deploymentForDeployment returns a deployment Deployment object
func (r *QueueAutoScalerReconciler) deploymentForDeployment(m *oldmonkv1.QueueAutoScaler) *appsv1.Deployment {
	// Attach volume if exist
	var volumes []corev1.Volume
	for _, v := range m.Spec.Volume {
		volumes = append(volumes, v)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.AppSpec.Name,
			Namespace: m.Namespace,
			Labels:    x.GetLabels(m).Spec.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &m.Spec.MinPods,
			Strategy: m.Spec.Strategy,
			Selector: &metav1.LabelSelector{
				MatchLabels: x.GetLabels(m).Spec.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: x.GetLabels(m).Spec.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: getContainer(m),
				},
			},
		},
	}
	if len(volumes) > 0 {
		dep.Spec.Template.Spec.Volumes = volumes
	}

	// Set Queue Based configmap and attach it to the container

	// Set Deployment instance as the owner of the Deployment.
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

func getContainer(m *oldmonkv1.QueueAutoScaler) []corev1.Container {

	// Attach Ports if exist
	var ports []corev1.ContainerPort
	for _, p := range m.Spec.AppSpec.Ports {
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: p.ContainerPort,
			Name:          p.Name,
		})
	}

	// Attach command if exist
	var command []string
	for _, c := range m.Spec.AppSpec.Command {
		command = append(command, c)
	}

	resource := corev1.ResourceRequirements(m.Spec.AppSpec.Resources)

	// Attach volumeMounts if exist
	var volumeMounts []corev1.VolumeMount
	for _, v := range m.Spec.AppSpec.VolumeMounts {
		volumeMounts = append(volumeMounts, v)
	}

	// Attach envFrom if exist
	var envFrom []corev1.EnvFromSource
	for _, v := range m.Spec.AppSpec.EnvFrom {
		envFrom = append(envFrom, v)
	}

	// Attach volumeMounts if exist
	var env []corev1.EnvVar
	for _, v := range m.Spec.AppSpec.Env {
		env = append(env, v)
	}

	var imagePullPolicy corev1.PullPolicy = "Always"
	if m.Spec.AppSpec.ImagePullPolicy != "" {
		imagePullPolicy = m.Spec.AppSpec.ImagePullPolicy
	}

	dep := []corev1.Container{{
		Image:           m.Spec.AppSpec.Image,
		Name:            m.Spec.AppSpec.Name,
		Command:         command,
		Ports:           ports,
		ImagePullPolicy: imagePullPolicy,
		Resources:       resource,
		LivenessProbe:   m.Spec.AppSpec.LivenessProbe,
		ReadinessProbe:  m.Spec.AppSpec.ReadinessProbe,
	}}
	if m.Spec.AppSpec.WorkingDir != "" {
		dep[0].WorkingDir = m.Spec.AppSpec.WorkingDir
	}
	if len(volumeMounts) > 0 {
		dep[0].VolumeMounts = volumeMounts
	}
	if len(envFrom) > 0 {
		dep[0].EnvFrom = envFrom
	}
	if len(env) > 0 {
		dep[0].Env = env
	}
	return dep
}

func (r *QueueAutoScalerReconciler) finalizeAutoScaler(m *oldmonkv1.QueueAutoScaler) error {
	// Delete Deployment
	// Check if the Deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: m.Spec.Deployment, Namespace: m.Namespace}, deployment)
	if err != nil {
		return err
	}
	err = r.client.Delete(context.TODO(), deployment)
	if err != nil {
		return err
	}
	return nil
}

func (r *QueueAutoScalerReconciler) addFinalizer(m *oldmonkv1.QueueAutoScaler) error {
	m.SetFinalizers(append(m.GetFinalizers(), queueAutoScalerFinalizer))
	// Update CR
	err := r.client.Update(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}