// Package scalex provides primitives for the scaling logic
package scalex

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"
	"github.com/evalsocket/oldmonk/x"
	"github.com/evalsocket/oldmonk/x/queuex"
	log "github.com/sirupsen/logrus"
	"github.com/vmg/backoff"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// "k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Scaler hold interval and client for scaling logic
type Scaler struct {
	client   client.Client
	interval time.Duration
}

// NewScalex create a new Scaler object
// It returns the Scaler
func NewScalex(mgr manager.Manager, interval time.Duration) *Scaler {
	return &Scaler{
		client:   mgr.GetClient(),
		interval: interval,
	}
}

// Run create a ticker with a fixed interval
// It returns the errors
func (s Scaler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.do(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

// Run min return the min
// It returns the int32
func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Run max return the max
// It returns the int32
func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// targetReplicas will check the desired replica for the scaling
// It returns the int32 and error
func (s Scaler) targetReplicas(size int32, scale *oldmonkv1.QueueAutoScaler, d *appsv1.Deployment) (int32, error) {
	replicas := d.Status.Replicas

	if size > scale.Spec.ScaleUp.Threshold {
		desired := replicas + scale.Spec.ScaleUp.Amount
		return min(desired, scale.Spec.MaxPods), nil
	} else if size <= scale.Spec.ScaleDown.Threshold {
		desired := replicas - scale.Spec.ScaleDown.Amount
		return max(desired, scale.Spec.MinPods), nil
	}
	return replicas, nil
}

// ExecuteScale will check the scaling policy and the scale according to the logic
// It returns the updated deployment,delta and error
func (s Scaler) ExecuteScale(ctx context.Context, scale *oldmonkv1.QueueAutoScaler) (*appsv1.Deployment, int32, error) {

	// // Get Secrets
	secret := &corev1.Secret{}
	err := s.client.Get(context.TODO(), client.ObjectKey{
		Namespace: scale.ObjectMeta.Namespace,
		Name:      scale.Spec.Secrets,
	}, secret)
	if err != nil {
		return nil, 0, err
	}

	scale.Spec.Option.Uri = string(secret.Data["URI"])

	// To-Do Set secret to env variable and remove it from crd defination

	c := queuex.NewQueueConnection(scale.Spec.Type, &scale.Spec.Option)
	if c == nil {
		return nil, 0, fmt.Errorf("error")
	}
	size := c.GetCount()
	_ = c.Close()
	fmt.Println("Queue :", scale.Spec.Type, "And Count : ", size)
	if size <= 0 {
		return nil, 0, errors.Unwrap(fmt.Errorf("......"))
	}

	// Fetch the QueueAutoScaler instance
	deployment := &appsv1.Deployment{}
	err = s.client.Get(context.TODO(), client.ObjectKey{
		Namespace: scale.ObjectMeta.Namespace,
		Name:      scale.Spec.Deployment,
	}, deployment)
	if err != nil {
		return nil, 0, err
	}

	// if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
	// 	return nil, 0, fmt.Errorf("deployment available replicas not at target. won't adjust")
	// }

	replicas, err := s.targetReplicas(size, scale, deployment)
	if err != nil {
		return nil, 0, err
	}

	delta := replicas - *deployment.Spec.Replicas

	deployment.Spec.Replicas = &replicas
	if err := s.client.Update(context.TODO(), deployment); err != nil {
		fmt.Println("unable to update deployment ")
	}

	return deployment, delta, nil
}

const (
	ReasonScaleDeployment       = "ScaleSuccess"
	ReasonFailedScaleDeployment = "ScaleFail"
)

// do excute in a fixed interval and it will start scaling policy for all queue autoscale deployments
func (s Scaler) do(ctx context.Context) {

	instance := &oldmonkv1.QueueAutoScalerList{}
	err := s.client.List(ctx, instance)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, scaler := range instance.Items {
		logger := log.WithFields(log.Fields{"delta": ""})
		op := func() error {
			deployment, delta, err := s.ExecuteScale(ctx, &scaler)
			if err != nil {
				logger.Warnf("unable to perform scale, will retry:")
				return err
			}
			logger.WithFields(log.Fields{"delta": delta, "desired": *deployment.Spec.Replicas, "available": deployment.Status.AvailableReplicas}).Info("Updated deployment")
			return nil
		}
		strategy := backoff.NewExponentialBackOff()
		strategy.MaxInterval = time.Second
		strategy.MaxElapsedTime = time.Second * 5
		strategy.InitialInterval = time.Millisecond * 100

		err := backoff.Retry(op, strategy)
		if err != nil {

			msg := fmt.Sprintf("error scaling: %s", err)
			logger.Error(msg)

		}
	}

	for _, scaler := range instance.Items {
		logger := log.WithFields(log.Fields{"delta": "bnvh"})
		op := func() error {
			// Update the QueueAutoScaler status with the pod names
			// List the pods for this worker's deployment
			podList := &corev1.PodList{}
			listOpts := []client.ListOption{
				client.InNamespace(scaler.Namespace),
				client.MatchingLabels(x.GetLabels(&scaler).Spec.Labels),
			}
			err = s.client.List(context.TODO(), podList, listOpts...)
			if err != nil {
				logger.Error(err, "Failed to list pods.", "QueueAutoScaler.Namespace", scaler.Namespace, "QueueAutoScaler.Name", scaler.Name)
				return nil
			}
			podNames := x.GetPodNames(podList.Items)

			// Update status.Nodes if needed
			if !reflect.DeepEqual(podNames, scaler.Status.Nodes) {
				scaler.Status.Nodes = podNames
				err := s.client.Update(context.TODO(), &scaler)
				if err != nil {
					logger.Error(err, "Failed to update QueueAutoScaler status.")
					return nil
				}
			}
			return nil
		}
		strategy := backoff.NewExponentialBackOff()
		strategy.MaxInterval = time.Second
		strategy.MaxElapsedTime = time.Second * 5
		strategy.InitialInterval = time.Millisecond * 100

		err := backoff.Retry(op, strategy)
		if err != nil {

			msg := fmt.Sprintf("error scaling: %s", err)
			logger.Error(msg)

		}

	}
}
