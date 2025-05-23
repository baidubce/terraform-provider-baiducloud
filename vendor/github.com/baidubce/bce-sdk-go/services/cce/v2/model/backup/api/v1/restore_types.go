/*
Copyright 2017, 2019 the Velero contributors.

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

package v1

// RestoreSpec defines the specification for a Velero restore.
type RestoreSpec struct {
	// BackupName is the unique name of the Velero backup to restore
	// from.
	BackupName string `json:"backupName"`

	// ScheduleName is the unique name of the Velero schedule to restore
	// from. If specified, and BackupName is empty, Velero will restore
	// from the most recent successful backup created from this schedule.
	// +optional
	ScheduleName string `json:"scheduleName,omitempty"`

	// IncludedNamespaces is a slice of namespace names to include objects
	// from. If empty, all namespaces are included.
	// +optional
	// +nullable
	IncludedNamespaces []string `json:"includedNamespaces,omitempty"`

	// ExcludedNamespaces contains a list of namespaces that are not
	// included in the restore.
	// +optional
	// +nullable
	ExcludedNamespaces []string `json:"excludedNamespaces,omitempty"`

	// IncludedResources is a slice of resource names to include
	// in the restore. If empty, all resources in the backup are included.
	// +optional
	// +nullable
	IncludedResources []string `json:"includedResources,omitempty"`

	// ExcludedResources is a slice of resource names that are not
	// included in the restore.
	// +optional
	// +nullable
	ExcludedResources []string `json:"excludedResources,omitempty"`

	// NamespaceMapping is a map of source namespace names
	// to target namespace names to restore into. Any source
	// namespaces not included in the map will be restored into
	// namespaces of the same name.
	// +optional
	NamespaceMapping map[string]string `json:"namespaceMapping,omitempty"`

	// LabelSelector is a LabelSelector to filter with
	// when restoring individual objects from the backup. If empty
	// or nil, all objects are included. Optional.
	// +optional
	// +nullable
	LabelSelector *LabelSelector `json:"labelSelector,omitempty"`

	// OrLabelSelectors is list of LabelSelector to filter with
	// when restoring individual objects from the backup. If multiple provided
	// they will be joined by the OR operator. LabelSelector as well as
	// OrLabelSelectors cannot co-exist in restore request, only one of them
	// can be used
	// +optional
	// +nullable
	OrLabelSelectors []*LabelSelector `json:"orLabelSelectors,omitempty"`

	// RestorePVs specifies whether to restore all included
	// PVs from snapshot
	// +optional
	// +nullable
	RestorePVs *bool `json:"restorePVs,omitempty"`

	// RestoreStatus specifies which resources we should restore the status
	// field. If nil, no objects are included. Optional.
	// +optional
	// +nullable
	RestoreStatus *RestoreStatusSpec `json:"restoreStatus,omitempty"`

	// PreserveNodePorts specifies whether to restore old nodePorts from backup.
	// +optional
	// +nullable
	PreserveNodePorts *bool `json:"preserveNodePorts,omitempty"`

	// IncludeClusterResources specifies whether cluster-scoped resources
	// should be included for consideration in the restore. If null, defaults
	// to true.
	// +optional
	// +nullable
	IncludeClusterResources *bool `json:"includeClusterResources,omitempty"`

	// Hooks represent custom behaviors that should be executed during or post restore.
	// +optional
	Hooks RestoreHooks `json:"hooks,omitempty"`

	// ExistingResourcePolicy specifies the restore behavior for the Kubernetes resource to be restored
	// +optional
	// +nullable
	ExistingResourcePolicy PolicyType `json:"existingResourcePolicy,omitempty"`

	// ItemOperationTimeout specifies the time used to wait for RestoreItemAction operations
	// The default value is 1 hour.
	// +optional
	ItemOperationTimeout string `json:"itemOperationTimeout,omitempty"`

	// ResourceModifier specifies the reference to JSON resource patches that should be applied to resources before restoration.
	// +optional
	// +nullable
	ResourceModifier *TypedLocalObjectReference `json:"resourceModifier,omitempty"`

	// UploaderConfig specifies the configuration for the restore.
	// +optional
	// +nullable
	UploaderConfig *UploaderConfigForRestore `json:"uploaderConfig,omitempty"`
}

// UploaderConfigForRestore defines the configuration for the restore.
type UploaderConfigForRestore struct {
	// WriteSparseFiles is a flag to indicate whether write files sparsely or not.
	// +optional
	// +nullable
	WriteSparseFiles *bool `json:"writeSparseFiles,omitempty"`
}

// RestoreHooks contains custom behaviors that should be executed during or post restore.
type RestoreHooks struct {
	Resources []RestoreResourceHookSpec `json:"resources,omitempty"`
}

type RestoreStatusSpec struct {
	// IncludedResources specifies the resources to which will restore the status.
	// If empty, it applies to all resources.
	// +optional
	// +nullable
	IncludedResources []string `json:"includedResources,omitempty"`

	// ExcludedResources specifies the resources to which will not restore the status.
	// +optional
	// +nullable
	ExcludedResources []string `json:"excludedResources,omitempty"`
}

// RestoreResourceHookSpec defines one or more RestoreResrouceHooks that should be executed based on
// the rules defined for namespaces, resources, and label selector.
type RestoreResourceHookSpec struct {
	// Name is the name of this hook.
	Name string `json:"name"`

	// IncludedNamespaces specifies the namespaces to which this hook spec applies. If empty, it applies
	// to all namespaces.
	// +optional
	// +nullable
	IncludedNamespaces []string `json:"includedNamespaces,omitempty"`

	// ExcludedNamespaces specifies the namespaces to which this hook spec does not apply.
	// +optional
	// +nullable
	ExcludedNamespaces []string `json:"excludedNamespaces,omitempty"`

	// IncludedResources specifies the resources to which this hook spec applies. If empty, it applies
	// to all resources.
	// +optional
	// +nullable
	IncludedResources []string `json:"includedResources,omitempty"`

	// ExcludedResources specifies the resources to which this hook spec does not apply.
	// +optional
	// +nullable
	ExcludedResources []string `json:"excludedResources,omitempty"`

	// LabelSelector, if specified, filters the resources to which this hook spec applies.
	// +optional
	// +nullable
	LabelSelector *LabelSelector `json:"labelSelector,omitempty"`

	//// PostHooks is a list of RestoreResourceHooks to execute during and after restoring a resource.
	//// +optional
	//PostHooks []RestoreResourceHook `json:"postHooks,omitempty"`
}

// ExecRestoreHook is a hook that uses pod exec API to execute a command inside a container in a pod
type ExecRestoreHook struct {
	// Container is the container in the pod where the command should be executed. If not specified,
	// the pod's first container is used.
	// +optional
	Container string `json:"container,omitempty"`

	// Command is the command and arguments to execute from within a container after a pod has been restored.
	// +kubebuilder:validation:MinItems=1
	Command []string `json:"command"`

	// OnError specifies how Velero should behave if it encounters an error executing this hook.
	// +optional
	OnError HookErrorMode `json:"onError,omitempty"`

	// ExecTimeout defines the maximum amount of time Velero should wait for the hook to complete before
	// considering the execution a failure.
	// +optional
	ExecTimeout string `json:"execTimeout,omitempty"`

	// WaitTimeout defines the maximum amount of time Velero should wait for the container to be Ready
	// before attempting to run the command.
	// +optional
	WaitTimeout string `json:"waitTimeout,omitempty"`

	// WaitForReady ensures command will be launched when container is Ready instead of Running.
	// +optional
	// +nullable
	WaitForReady *bool `json:"waitForReady,omitempty"`
}

// RestorePhase is a string representation of the lifecycle phase
type RestorePhase string

const (
	// RestorePhaseNew means the restore has been created but not
	// yet processed by the RestoreController
	RestorePhaseNew RestorePhase = "New"

	// RestorePhaseFailedValidation means the restore has failed
	// the controller's validations and therefore will not run.
	RestorePhaseFailedValidation RestorePhase = "FailedValidation"

	// RestorePhaseInProgress means the restore is currently executing.
	RestorePhaseInProgress RestorePhase = "InProgress"

	// RestorePhaseWaitingForPluginOperations means the restore of
	// Kubernetes resources and other async plugin operations was
	// successful and plugin operations are still ongoing.  The
	// restore is not complete yet.
	RestorePhaseWaitingForPluginOperations RestorePhase = "WaitingForPluginOperations"

	// RestorePhaseWaitingForPluginOperationsPartiallyFailed means
	// the restore of Kubernetes resources and other async plugin
	// operations partially failed (final phase will be
	// PartiallyFailed) and other plugin operations are still
	// ongoing.  The restore is not complete yet.
	RestorePhaseWaitingForPluginOperationsPartiallyFailed RestorePhase = "WaitingForPluginOperationsPartiallyFailed"

	// RestorePhaseCompleted means the restore has run successfully
	// without errors.
	RestorePhaseCompleted RestorePhase = "Completed"

	// RestorePhasePartiallyFailed means the restore has run to completion
	// but encountered 1+ errors restoring individual items.
	RestorePhasePartiallyFailed RestorePhase = "PartiallyFailed"

	// RestorePhaseFailed means the restore was unable to execute.
	// The failing error is recorded in status.FailureReason.
	RestorePhaseFailed RestorePhase = "Failed"

	// PolicyTypeNone means velero will not overwrite the resource
	// in cluster with the one in backup whether changed/unchanged.
	PolicyTypeNone PolicyType = "none"

	// PolicyTypeUpdate means velero will try to attempt a patch on
	// the changed resources.
	PolicyTypeUpdate PolicyType = "update"
)

// RestoreStatus captures the current status of a Velero restore
type RestoreStatus struct {
	// Phase is the current state of the Restore
	// +optional
	Phase RestorePhase `json:"phase,omitempty"`

	// ValidationErrors is a slice of all validation errors (if
	// applicable)
	// +optional
	// +nullable
	ValidationErrors []string `json:"validationErrors,omitempty"`

	// Warnings is a count of all warning messages that were generated during
	// execution of the restore. The actual warnings are stored in object storage.
	// +optional
	Warnings int `json:"warnings,omitempty"`

	// Errors is a count of all error messages that were generated during
	// execution of the restore. The actual errors are stored in object storage.
	// +optional
	Errors int `json:"errors,omitempty"`

	// FailureReason is an error that caused the entire restore to fail.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// StartTimestamp records the time the restore operation was started.
	// The server's time is used for StartTimestamps
	// +optional
	// +nullable
	StartTimestamp *Time `json:"startTimestamp,omitempty"`

	// CompletionTimestamp records the time the restore operation was completed.
	// Completion time is recorded even on failed restore.
	// The server's time is used for StartTimestamps
	// +optional
	// +nullable
	CompletionTimestamp *Time `json:"completionTimestamp,omitempty"`

	// Progress contains information about the restore's execution progress. Note
	// that this information is best-effort only -- if Velero fails to update it
	// during a restore for any reason, it may be inaccurate/stale.
	// +optional
	// +nullable
	Progress *RestoreProgress `json:"progress,omitempty"`

	// RestoreItemOperationsAttempted is the total number of attempted
	// async RestoreItemAction operations for this restore.
	// +optional
	RestoreItemOperationsAttempted int `json:"restoreItemOperationsAttempted,omitempty"`

	// RestoreItemOperationsCompleted is the total number of successfully completed
	// async RestoreItemAction operations for this restore.
	// +optional
	RestoreItemOperationsCompleted int `json:"restoreItemOperationsCompleted,omitempty"`

	// RestoreItemOperationsFailed is the total number of async
	// RestoreItemAction operations for this restore which ended with an error.
	// +optional
	RestoreItemOperationsFailed int `json:"restoreItemOperationsFailed,omitempty"`

	// HookStatus contains information about the status of the hooks.
	// +optional
	// +nullable
	HookStatus *HookStatus `json:"hookStatus,omitempty"`
}

// RestoreProgress stores information about the restore's execution progress
type RestoreProgress struct {
	// TotalItems is the total number of items to be restored. This number may change
	// throughout the execution of the restore due to plugins that return additional related
	// items to restore
	// +optional
	TotalItems int `json:"totalItems,omitempty"`
	// ItemsRestored is the number of items that have actually been restored so far
	// +optional
	ItemsRestored int `json:"itemsRestored,omitempty"`
}

// Restore is a Velero resource that represents the application of
// resources from a Velero backup to a target Kubernetes cluster.
type Restore struct {
	TypeMeta `json:",inline"`

	// +optional
	ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec RestoreSpec `json:"spec,omitempty"`

	// +optional
	Status RestoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RestoreList is a list of Restores.
type RestoreList struct {
	TypeMeta `json:",inline"`

	// +optional
	ListMeta `json:"metadata"`

	Items []Restore `json:"items"`
}

// PolicyType helps specify the ExistingResourcePolicy
type PolicyType string
