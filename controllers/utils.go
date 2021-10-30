package controllers

import "github.com/manifest/api/application/v1alpha1"

func convertObjState(state string) (frontState string) {
	switch state {
	case v1alpha1.ClusterInitState, v1alpha1.StatusCreating:
		frontState = v1alpha1.FrontCreating
	case v1alpha1.ClusterUpdateState, v1alpha1.StatusInProgress, v1alpha1.PgclusterStateProcessed:
		frontState = v1alpha1.FrontUpdating
	case v1alpha1.StatusCompleted, v1alpha1.PgclusterStateCreated:
		frontState = v1alpha1.FrontCompleted
	case v1alpha1.ClusterReadyState, v1alpha1.StatusRunning, v1alpha1.PgclusterStateInitialized:
		frontState = v1alpha1.FrontRunning
	case v1alpha1.ClusterCloseState, v1alpha1.PgclusterStateShutdown:
		frontState = v1alpha1.FrontClosed
	case v1alpha1.StatusCreateFailed:
		frontState = v1alpha1.FrontCreateFailed
	case v1alpha1.StatusUpdateFailed:
		frontState = v1alpha1.FrontUpdateFailed
	case v1alpha1.StatusTerminating:
		frontState = v1alpha1.FrontTerminating
	case v1alpha1.PgclusterStateBootstrapping:
		frontState = v1alpha1.StatusBootstrapping
	case v1alpha1.PgclusterStateBootstrapped:
		frontState = v1alpha1.StatusBootstrapped
	case v1alpha1.PgclusterStateRestore:
		frontState = v1alpha1.StatusRestoring
	default:
		frontState = v1alpha1.ClusterStatusUnknown
	}
	return
}
