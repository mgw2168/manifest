package v1alpha1

// ClusterState defines cluster state.
type FrontState string

const (
	ClusterNameLabelKey = "kubesphere.io/cluster"
	NamespaceLabelKey   = "kubesphere.io/namespace"

	// custom resource status
	ManifestCreating = "ManifestCreating"
	ManifestCreated  = "ManifestCreated"
	Failed           = "Failed"
	Error            = "Error"

	// front state
	FrontCreating       string = "Creating"      // 创建中
	FrontUpdating       string = "InProgress"    // 更新中
	FrontCompleted      string = "Completed"     // 创建完成
	FrontRunning        string = "Running"       // 运行中
	FrontClosed         string = "Closed"        // 已关机
	FrontCreateFailed   string = "CreateFailed"  // 创建失败
	FrontUpdateFailed   string = "UpdateFailed"  // 更新失败
	FrontTerminating    string = "Terminating"   // 删除中
	StatusBootstrapping string = "Bootstrapping" // 从已有的数据源启动集群中
	StatusBootstrapped  string = "Bootstrapped"  // 从已有的数据源启动集群完成
	StatusRestoring     string = "Restoring"     // 恢复中

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
