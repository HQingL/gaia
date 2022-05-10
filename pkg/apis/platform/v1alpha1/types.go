package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster",categories=gaia
// +k8s:openapi-gen=true
type Target struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec TargetSpec `json:"spec,omitempty"`
	// +optional
	Status TargetStatus `json:"status,omitempty"`
}

type TargetSpec struct {
	// +optional
	Self bool `json:"self,omitempty"`
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
	// +required
	ReportFrequency *int64 `json:"reportFrequency,omitempty" protobuf:"varint,5,opt,name=reportFrequency"`
	// +required
	CollectFrequency *int64 `json:"collectFrequency,omitempty" protobuf:"varint,5,opt,name=collectFrequency"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	ParentURL string `json:"parenturl,omitempty" protobuf:"bytes,1,opt,name=parenturl"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	BootstrapToken string `json:"bootstraptoken,omitempty" protobuf:"bytes,1,opt,name=bootstraptoken"`
}

type TargetStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetList contains a list of Target
type TargetList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Target `json:"items"`
}

// ManagedClusterSpec defines the desired state of ManagedCluster
type ManagedClusterSpec struct {
	// ClusterID, a Random (Version 4) UUID, is a unique value in time and space value representing for child cluster.
	// It is typically generated by the controller agent on the successful creation of a "self-cluster" Lease
	// in the child cluster.
	// Also it is not allowed to change on PUT operations.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	ClusterID types.UID `json:"clusterId"`
	// Taints has the "effect" on any resource that does not tolerate the Taint.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
}

// ManagedClusterStatus defines the observed state of ManagedCluster
type ManagedClusterStatus struct {
	// lastObservedTime is the time when last status from the series was seen before last heartbeat.
	// RFC 3339 date and time at which the object was acknowledged by the controller Agent.
	// +optional
	LastObservedTime metav1.Time `json:"lastObservedTime,omitempty"`

	// k8sVersion is the Kubernetes version of the cluster
	// +optional
	KubernetesVersion string `json:"k8sVersion,omitempty"`

	// platform indicates the running platform of the cluster
	// +optional
	Platform string `json:"platform,omitempty"`

	// APIServerURL indicates the advertising url/address of managed Kubernetes cluster
	// +optional
	APIServerURL string `json:"apiserverURL,omitempty"`

	// Healthz indicates the healthz status of the cluster
	// which is deprecated since Kubernetes v1.16. Please use Livez and Readyz instead.
	// Leave it here only for compatibility.
	// +optional
	Healthz bool `json:"healthz"`

	// Livez indicates the livez status of the cluster
	// +optional
	Livez bool `json:"livez"`

	// Readyz indicates the readyz status of the cluster
	// +optional
	Readyz bool `json:"readyz"`

	// Allocatable is the sum of allocatable resources for nodes in the cluster
	// +optional
	Allocatable corev1.ResourceList `json:"allocatable,omitempty"`

	// Capacity is the sum of capacity resources for nodes in the cluster
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// Available is the sum of Available resources for nodes in the cluster
	// +optional
	Available corev1.ResourceList `json:"available,omitempty"`

	// ClusterCIDR is the CIDR range of the cluster
	// +optional
	ClusterCIDR string `json:"clusterCIDR,omitempty"`

	// ServcieCIDR is the CIDR range of the services
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// NodeStatistics is the info summary of nodes in the cluster
	// +optional
	NodeStatistics NodeStatistics `json:"nodeStatistics,omitempty"`

	// Conditions is an array of current cluster conditions.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// heartbeatFrequencySeconds is the frequency at which the agent reports current cluster status
	// +optional
	HeartbeatFrequencySeconds *int64 `json:"heartbeatFrequencySeconds,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=mcls,categories=gaia
// +kubebuilder:printcolumn:name="CLUSTER ID",type=string,JSONPath=`.spec.clusterId`,description="The unique id for the cluster"
// +kubebuilder:printcolumn:name="KUBERNETES",type=string,JSONPath=".status.k8sVersion"
// +kubebuilder:printcolumn:name="READYZ",type=string,JSONPath=".status.readyz"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ManagedCluster is the Schema for the managedclusters API
type ManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterSpec   `json:"spec,omitempty"`
	Status ManagedClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterList contains a list of ManagedCluster
type ManagedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedCluster `json:"items"`
}

type NodeStatistics struct {
	// ReadyNodes is the number of ready nodes in the cluster
	// +optional
	ReadyNodes int32 `json:"readyNodes,omitempty"`

	// NotReadyNodes is the number of not ready nodes in the cluster
	// +optional
	NotReadyNodes int32 `json:"notReadyNodes,omitempty"`

	// UnknownNodes is the number of unknown nodes in the cluster
	// +optional
	UnknownNodes int32 `json:"unknownNodes,omitempty"`

	// LostNodes is the number of states lost nodes in the cluster
	// +optional
	LostNodes int32 `json:"lostNodes,omitempty"`
}

const (
	// ClusterReady means cluster is ready.
	ClusterReady = "Ready"
)

// ClusterRegistrationRequestSpec defines the desired state of ClusterRegistrationRequest
type ClusterRegistrationRequestSpec struct {
	// ClusterID, a Random (Version 4) UUID, is a unique value in time and space value representing for child cluster.
	// It is typically generated by the gaia agent on the successful creation of a "self-cluster" Lease
	// in the child cluster.
	// Also it is not allowed to change on PUT operations.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	ClusterID types.UID `json:"clusterId"`

	// ClusterNamePrefix is the prefix of cluster name.
	// a lower case alphanumeric characters or '-', and must start and end with an alphanumeric character
	//
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:Pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9]))*"
	ClusterNamePrefix string `json:"clusterNamePrefix,omitempty"`

	// ClusterName is the cluster name.
	// a lower case alphanumeric characters or '-', and must start and end with an alphanumeric character
	//
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:Pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9]))*"
	ClusterName string `json:"clusterName,omitempty"`

	// ClusterLabels is the labels of the child cluster.
	//
	// +optional
	// +kubebuilder:validation:Type=object
	ClusterLabels map[string]string `json:"clusterLabels,omitempty"`
}

// ClusterRegistrationRequestStatus defines the observed state of ClusterRegistrationRequest
type ClusterRegistrationRequestStatus struct {
	// DedicatedNamespace is a dedicated namespace for the child cluster, which is created in the parent cluster.
	//
	// +optional
	DedicatedNamespace string `json:"dedicatedNamespace,omitempty"`

	// DedicatedToken is populated by parent cluster when Result is RequestApproved.
	// With this token, the client could have full access on the resources created in DedicatedNamespace.
	//
	// +optional
	DedicatedToken []byte `json:"token,omitempty"`

	// CACertificate is the public certificate that is the root of trust for parent cluster
	// The certificate is encoded in PEM format.
	//
	// +optional
	CACertificate []byte `json:"caCertificate,omitempty"`

	// Result indicates whether this request has been approved.
	// When all necessary objects have been created and ready for child cluster registration,
	// this field will be set to "Approved". If any illegal updates on this object, "Illegal" will be set to this filed.
	//
	// +optional
	Result *ApprovedResult `json:"result,omitempty"`

	// ErrorMessage tells the reason why the request is not approved successfully.
	//
	// +optional
	ErrorMessage string `json:"errorMessage,omitempty"`

	// ManagedClusterName is the name of ManagedCluster object in the parent cluster corresponding to the child cluster
	//
	// +optional
	ManagedClusterName string `json:"managedClusterName,omitempty"`
}

type ApprovedResult string

// These are the possible results for a cluster registration request.
const (
	RequestDenied   ApprovedResult = "Denied"
	RequestApproved ApprovedResult = "Approved"
	RequestFailed   ApprovedResult = "Failed"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster",shortName=clsrr,categories=gaia
// +kubebuilder:printcolumn:name="CLUSTER ID",type=string,JSONPath=`.spec.clusterId`,description="The unique id for the cluster"
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=`.status.result`,description="The status of current cluster registration request"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterRegistrationRequest is the Schema for the clusterregistrationrequests API
type ClusterRegistrationRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterRegistrationRequestSpec   `json:"spec,omitempty"`
	Status ClusterRegistrationRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRegistrationRequestList contains a list of ClusterRegistrationRequest
type ClusterRegistrationRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRegistrationRequest `json:"items"`
}
