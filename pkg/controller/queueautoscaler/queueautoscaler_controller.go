package queueautoscaler

import (
	"context"
	"reflect"
	"time"
	"os"
	"fmt"
	"strconv"

	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"

	"github.com/evalsocket/oldmonk/x"
	"github.com/evalsocket/oldmonk/x/scalex"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_queueAutoScaler")

const queueAutoScalerFinalizer = "finalizer.oldmonk.evalsocket.in"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new QueueAutoScaler Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileQueueAutoScaler{client: mgr.GetClient(), scheme: mgr.GetScheme(), mgr: mgr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("queueautoscaler-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource QueueAutoScaler
	err = c.Watch(&source.Kind{Type: &oldmonkv1.QueueAutoScaler{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner QueueAutoScaler
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &oldmonkv1.QueueAutoScaler{},
	})
	if err != nil {
		return err
	}
	poolDuration := time.Second*120
  if len(os.Getenv("POOL_DURATION")) != 0 {
		i, err := strconv.ParseInt(os.Getenv("POOL_DURATION"),10, 64);
		if  err != nil {
			 log.Error("Error in converting pool duration string to time")
		}
		poolDuration = time.Second*time.Duration(i)
	}
	checkState := func() {
		time.Sleep(2 * time.Second)
		scaler := scalex.NewScalex(mgr, poolDuration)
		scaler.Run(context.TODO())
	}
	go checkState()
	return nil
}

// blank assignment to verify that ReconcileQueueAutoScaler implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileQueueAutoScaler{}

// ReconcileQueueAutoScaler reconciles a QueueAutoScaler object
type ReconcileQueueAutoScaler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	mgr    manager.Manager
	scheme *runtime.Scheme
}

// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileQueueAutoScaler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling QueueAutoScaler")
	// Fetch the QueueAutoScaler instance
	queueAutoScaler := &oldmonkv1.QueueAutoScaler{}
	err := r.client.Get(context.TODO(), request.NamespacedName, queueAutoScaler)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("QueueAutoScaler resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get QueueAutoScaler.")
		return reconcile.Result{}, err
	}

	if !queueAutoScaler.Spec.Autopilot {
		return reconcile.Result{}, nil
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
				return reconcile.Result{}, err
			}

			// Remove queueAutoScalerFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			queueAutoScaler.SetFinalizers(x.Remove(queueAutoScaler.GetFinalizers(), queueAutoScalerFinalizer))
			err := r.client.Update(context.TODO(), queueAutoScaler)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	} else {
		// Check if the Deployment already exists, if not create a new one
		deployment := &appsv1.Deployment{}
		dep := r.deploymentForDeployment(queueAutoScaler)
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: queueAutoScaler.Spec.Deployment, Namespace: queueAutoScaler.Namespace}, deployment)
		if err != nil && errors.IsNotFound(err) {
			// Define a new Deployment
			reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			err = r.client.Create(context.TODO(), dep)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully - return and requeue
			// NOTE: that the requeue is made with the purpose to provide the deployment object for the next step to ensure the deployment size is the same as the spec.
			// Also, you could GET the deployment object again instead of requeue if you wish. See more over it here: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Deployment.")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Updating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Update(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to update new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Ensure the deployment size is the same as the spec
		// TO-Do perform queue logic and then scale
		scaler := scalex.NewScalex(r.mgr, time.Second*60)
		deployment, delta, err := scaler.ExecuteScale(context.TODO(), queueAutoScaler)

		if err != nil {
			reqLogger.Error(err, "unable to perform desired replica, will retry:")
			return reconcile.Result{}, err
		}
		log.WithValues("Request.delta", delta, "Request.desired", *deployment.Spec.Replicas, "Request.available", deployment.Status.AvailableReplicas).Info("Updated deployment")

		// Update the QueueAutoScaler status with the pod names
		// List the pods for this worker's deployment
		podList := &corev1.PodList{}
		listOpts := []client.ListOption{
			client.InNamespace(queueAutoScaler.Namespace),
			client.MatchingLabels(x.GetLabels(queueAutoScaler).Spec.Labels),
		}
		err = r.client.List(context.TODO(), podList, listOpts...)
		if err != nil {
			reqLogger.Error(err, "Failed to list pods.", "QueueAutoScaler.Namespace", queueAutoScaler.Namespace, "QueueAutoScaler.Name", queueAutoScaler.Name)
			return reconcile.Result{}, err
		}
		podNames := x.GetPodNames(podList.Items)

		// Update status.Nodes if needed
		if !reflect.DeepEqual(podNames, queueAutoScaler.Status.Nodes) {
			queueAutoScaler.Status.Nodes = podNames
			err := r.client.Status().Update(context.TODO(), queueAutoScaler)
			if err != nil {
				reqLogger.Error(err, "Failed to update QueueAutoScaler status.")
				return reconcile.Result{}, err
			}
		}
	}

	// Add finalizer for this CR
	if !x.Contains(queueAutoScaler.GetFinalizers(), queueAutoScalerFinalizer) {
		if err := r.addFinalizer(queueAutoScaler); err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

// deploymentForDeployment returns a deployment Deployment object
func (r *ReconcileQueueAutoScaler) deploymentForDeployment(m *oldmonkv1.QueueAutoScaler) *appsv1.Deployment {
	// Attach volume if exist
	var volumes []corev1.Volume
	for _, v := range m.Spec.Volume {
		volumes = append(volumes, v)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.AppSpec.Name,
			Namespace: m.Namespace,
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

func (r *ReconcileQueueAutoScaler) finalizeAutoScaler(m *oldmonkv1.QueueAutoScaler) error {
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

func (r *ReconcileQueueAutoScaler) addFinalizer(m *oldmonkv1.QueueAutoScaler) error {
	m.SetFinalizers(append(m.GetFinalizers(), queueAutoScalerFinalizer))
	// Update CR
	err := r.client.Update(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}
