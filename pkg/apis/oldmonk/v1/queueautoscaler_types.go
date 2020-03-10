package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AutoScalerSpec defines the desired state of QueueAutoScaler
// +k8s:openapi-gen=true
type QueueAutoScalerSpec struct {
	// Type contains the user-specified Queue type
	// Type can be RABITMQ/BEANSTALK/NATS/SQS
	Type string `json:"type"`

	// Option contains Queue Details
	Option ListOptions `json:"option"`

	// MinPods for deployment
	MinPods int32 `json:"minPods"`

	// MaxPods for deployment
	MaxPods int32 `json:"maxPods"`

	// targetMessagesPerWorker is the number used to find number of pod needed. It's a optional parameter and used in case of policy TARGET
	TargetMessagesPerWorker int32 `json:"targetMessagesPerWorker"`

	// ScaleUp contains scale up policy
	ScaleUp ScaleSpec `json:"scaleUp"`

	// ScaleDown contains scale down policy
	ScaleDown ScaleSpec `json:"scaleDown"`

	// Policy contains name of policy possible value is THRESOLD/TARGET
	Policy string `json:"policy"`

	// Secrets contains secret name. Pass secure connection uri using secrets
	Secrets string `json:"secrets"`

	// Deployment contains deployment name which user want to scale
	Deployment string `json:"deployment"`

	// AppSpec contains deployment specification used same as deployment file (optional)
	AppSpec corev1.Container `json:"appSpec,omitempty"`

	// Labels contains key value pair for deployment (optional)
	Labels map[string]string `json:"labels,omitempty"`

	// Strategy contains deployment strategy (optional)
	Strategy appsv1.DeploymentStrategy `json:"strategy,omitempty"`

	// Volume cotains a list of volume (optional)
	Volume []corev1.Volume `json:"volume,omitempty"`

	// Autopilot is a bool value. false means it only auto scale deployment and true means it will manage the entire life cycle of deployment
	Autopilot bool `json:"autopilot"`
}

// ListOptions defines the desired state of Queue
// +k8s:openapi-gen=true
type ListOptions struct {

	// Uri (Don't use it's depricated). Use secrets
	Uri string `json:"uri,omitempty"`

	// Region is a optional parameter and used in case of SQS
	Region string `json:"region,omitempty"`

	// Type is a optional parameter and used in case of Rabbitmq(exchange type)
	Type string `json:"type,omitempty"`

	// Queue is a optional parameter and used in case of Rabbitmq and nats
	Queue string `json:"queue,omitempty"`

	// Exchange is a optional parameter and used in case of Rabbitmq
	Exchange string `json:"exchange,omitempty"`

	// Tube is a optional parameter and used in case of Beanstalk
	Tube string `json:"tube,omitempty"`

	// Group is a optional parameter and used in case of Kafka
	Group string `json:"group,omitempty"`

	// Topic is a optional parameter and used in case of Kafka
	Topic string `json:"topic,omitempty"`
}

// ScaleSpec defines the desired state of Autoscaler
// +k8s:openapi-gen=true
type ScaleSpec struct {

	// Threshold is the amount of messages in queue
	Threshold int32 `json:"threshold"`

	// Amount is the number by which you want to scale
	Amount int32 `json:"amount"`
}

// QueueAutoScalerStatus defines the observed state of QueueAutoScaler
type QueueAutoScalerStatus struct {
	// Nodes are the names of the pods
	// +listType=set
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueueAutoScaler is the Schema for the queueautoscalers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=queueautoscalers,scope=Namespaced
type QueueAutoScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueueAutoScalerSpec   `json:"spec,omitempty"`
	Status QueueAutoScalerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueueAutoScalerList contains a list of QueueAutoScaler
type QueueAutoScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueueAutoScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QueueAutoScaler{}, &QueueAutoScalerList{})
}
