// Package scalex provides primitives for the scaling logic
package scalex

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"


	"github.com/prometheus/client_golang/prometheus"

	oldmonkv1 "github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1"
	"github.com/remotegarage/oldmonk/pkg/common"
	"github.com/remotegarage/oldmonk/pkg/common/queuex"
  
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

var (
	loopDurationSeconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "oldmonk",
			Subsystem: "controller",
			Name:      "loop_duration_seconds",
			Help:      "Number of seconds to complete the control loop succesfully, partitioned by oldmonk name and namespace",
		},
		[]string{"oldmonk", "namespace"},
	)

	loopCountSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "oldmonk",
			Subsystem: "controller",
			Name:      "loop_count_success",
			Help:      "How many times the control loop executed succesfully, partitioned by oldmonk name and namespace",
		},
		[]string{"oldmonk", "namespace"},
	)
)

// NewScalex create a new Scaler object
// It returns the Scaler
func NewScalex(mgr manager.Manager, interval time.Duration) *Scaler {
	prometheus.MustRegister(loopDurationSeconds)
	prometheus.MustRegister(loopCountSuccess)
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

// convertDesiredReplicasWithRules
func convertDesiredReplicasWithRules(desired int32, min int32, max int32) int32 {
	if desired > max {
		return max
	}
	if desired <= min {
		return min
	}
	return desired
}

// targetReplicas will check the desired replica for the scaling
// It returns the int32 and error
func (s Scaler) targetReplicas(size int32, scale *oldmonkv1.QueueAutoScaler, d *appsv1.Deployment) (int32, error) {
	replicas := d.Status.Replicas
	if len(scale.Spec.Policy) <= 0 {
		scale.Spec.Policy = "THRESOLD"
	}
	if scale.Spec.Policy == "THRESOLD" {
		if size > scale.Spec.ScaleUp.Threshold {
			desired := replicas + scale.Spec.ScaleUp.Amount
			return min(desired, scale.Spec.MaxPods), nil
		} else if size <= scale.Spec.ScaleDown.Threshold {
			desired := replicas - scale.Spec.ScaleDown.Amount
			return max(desired, scale.Spec.MinPods), nil
		}
		return replicas, nil
	} else if scale.Spec.Policy == "TARGET" {
		tolerance := 0.1
		usageRatio := float64(size) / float64(scale.Spec.TargetMessagesPerWorker)
		if size >= 0 {
			// return the current replicas if the change would be too small
			if size < scale.Spec.TargetMessagesPerWorker || math.Abs(1.0-usageRatio) <= tolerance {
				return replicas, nil
			}
			desiredWorkers := int32(math.Ceil(usageRatio))
			return convertDesiredReplicasWithRules(desiredWorkers, scale.Spec.MinPods, scale.Spec.MaxPods), nil
		}
		return replicas, nil
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

	c := queuex.NewQueueConnection(scale.Spec.Type, &scale.Spec.Option)
	if c == nil {
		return nil, 0, fmt.Errorf("error")
	}
	size := c.GetCount()
	_ = c.Close()
	fmt.Println("Queue :", scale.Spec.Type, "And Count : ", size)
	if size < 0 {
		return nil, 0, errors.Unwrap(fmt.Errorf("Something Goes wrong with queue drivers"))
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

	if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
		return nil, 0, nil
	}

	replicas, err := s.targetReplicas(size, scale, deployment)

	if err != nil {
		return nil, 0, err
	}
	delta := replicas - *deployment.Spec.Replicas

	deployment.Spec.Replicas = &replicas
	if err := s.client.Update(context.TODO(), deployment); err != nil {
		log.Error("unable to update deployment ")
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
		log.Error("Error", err)
	}

	var jobs chan oldmonkv1.QueueAutoScaler
	for  i := 0; i > 3; i++ {
		s.Worker(jobs)
	}
	for _, scaler := range instance.Items {
		jobs <- scaler
	}
}

func (s Scaler) Worker(jobs chan oldmonkv1.QueueAutoScaler) {
	for scaler := range jobs {
		// Run This logic in a goroutine
		ctx := context.TODO()
		now := time.Now()
		logger := log.WithFields(log.Fields{"delta": ""})
		op := func() error {
			deployment, delta, err := s.ExecuteScale(ctx, &scaler)
			if err != nil {
				logger.Warnf("unable to perform scale, will retry:")
				return err
			}
			if deployment != nil {
				logger.WithFields(log.Fields{"Delta": delta, "Desired": *deployment.Spec.Replicas, "Available": deployment.Status.AvailableReplicas, "Queue Type": scaler.Spec.Type, "Deployment ": scaler.Spec.Deployment, "Policy": scaler.Spec.Policy}).Info("Updated deployment")
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
		loopDurationSeconds.WithLabelValues(
			scaler.Spec.Deployment,
			scaler.Namespace,
		).Set(time.Since(now).Seconds())
		loopCountSuccess.WithLabelValues(
			scaler.Spec.Deployment,
			scaler.Namespace,
		).Inc()
	}
	

}
