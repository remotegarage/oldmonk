package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AutoScalerSpec defines the desired state of QueueAutoScaler
// +k8s:openapi-gen=true
type QueueAutoScalerSpec struct {
	Type       string                    `json:"type"`
	Option     ListOptions               `json:"option"`
	MinPods    int32                     `json:"minPods"`
	MaxPods    int32                     `json:"maxPods"`
	ScaleUp    ScaleSpec                 `json:"scaleUp"`
	ScaleDown  ScaleSpec                 `json:"scaleDown"`
	Secrets    string                    `json:"secrets"`
	Deployment string                    `json:"deployment"`
	AppSpec    corev1.Container          `json:"appSpec,omitempty"`
	Labels     map[string]string         `json:"labels,omitempty"`
	Strategy   appsv1.DeploymentStrategy `json:"strategy,omitempty"`
	Volume     []corev1.Volume           `json:"volume,omitempty"`
	Autopilot  bool                      `json:"autopilot"`
}

// ListOptions defines the desired state of Queue
// +k8s:openapi-gen=true
type ListOptions struct {
	Uri      string `json:"uri,omitempty"`
	Region   string `json:"region,omitempty"`
	Type     string `json:"type,omitempty"`
	Queue    string `json:"queue,omitempty"`
	VsHost   string `json:"vshost,omitempty"`
	Key      string `json:"key,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	Tube     string `json:"tube,omitempty"`
	Group    string `json:"group,omitempty"`
	Topic    string `json:"topic,omitempty"`
}

// ScaleSpec defines the desired state of Autoscaler
// +k8s:openapi-gen=true
type ScaleSpec struct {
	Threshold int32 `json:"threshold"`
	Amount    int32 `json:"amount"`
}

// QueueAutoScalerStatus defines the observed state of QueueAutoScaler
type QueueAutoScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Nodes are the names of the memcached pods
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
