package v1alpha1

// ClusterState defines cluster state.
type FrontState string

const (
	// custom resource status
	Creating      = "Creating"
	Failed        = "Failed"
	Error         = "Error"
	AlreadyExists = "AlreadyExists"

	// front state
	FrontCreating       string = "Creating"
	FrontUpdating       string = "InProgress"
	FrontCompleted      string = "Completed"
	FrontRunning        string = "Running"
	FrontClosed         string = "Closed"
	FrontCreateFailed   string = "Create Failed"
	FrontUpdateFailed   string = "Update Failed"
	FrontTerminating    string = "Terminating"
	StatusBootstrapping string = "Bootstrapping"
	StatusBootstrapped  string = "Bootstrapped"
	StatusRestoring     string = "Restoring"

	// MySQL state
	// ClusterInitState  indicates whether the cluster is initializing.
	ClusterInitState string = "Initializing"
	// ClusterUpdateState indicates whether the cluster is being updated.
	ClusterUpdateState string = "Updating"
	// ClusterReadyState indicates whether all containers in the pod are ready.
	ClusterReadyState string = "Ready"
	// ClusterCloseState indicates whether the cluster is closed.
	ClusterCloseState string = "Closed"

	// ClickHouse state
	StatusCreating     = "Creating"
	StatusInProgress   = "InProgress"
	StatusCompleted    = "Completed"
	StatusRunning      = "Running"
	StatusCreateFailed = "CreateFailed"
	StatusUpdateFailed = "UpdateFailed"
	StatusTerminating  = "Terminating"

	// PostgreSQL state
	PgclusterStateCreated       string = "pgcluster Created"
	PgclusterStateProcessed     string = "pgcluster Processed"
	PgclusterStateInitialized   string = "pgcluster Initialized"
	PgclusterStateBootstrapping string = "pgcluster Bootstrapping"
	PgclusterStateBootstrapped  string = "pgcluster Bootstrapped"
	PgclusterStateRestore       string = "pgcluster Restoring"
	PgclusterStateShutdown      string = "pgcluster Shutdown"

	// kind of operator cr
	DBTypeMysql = "MySQL"

	// kind of cluster
	KindPostgreSQLCluster = "PostgreSQLCluster"
	KindMysqlCluster      = "MysqlCluster"
	KindClickHouseCluster = "ClickHouseInstallation"
	KindPgCluster         = "Pgcluster"
	KindPgClusterVersion  = "radondb.com/v1"

	// cluster status
	ClusterStatusUnknown = "unknown"
	// suffix of secret name
	SuffixSecretName = "-userpassword-secret"
)
