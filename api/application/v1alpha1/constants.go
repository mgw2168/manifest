package v1alpha1

const (
	// custom resource status
	Creating = "Creating"
	Deleting = "Deleting"
	Failed   = "Failed"
	Error    = "Error"
	Created  = "Created"

	NamespaceLabelKey = "kubesphere.io/namespace"

	// kind of operator cr
	DBTypeClickHouse = "ClickHouseInstallation"
	DBTypePostgreSQL = "PostgreSQL"
	DBTypeMysql      = "MySQL"

	// type of the cluster application
	ClusterAppTypeClickHouse = "ClickHouse"
	ClusterAPPTypePostgreSQL = "PostgreSQL"
	ClusterAPPTypeMySQL      = "MySQL"

	// cluster status
	ClusterStatusUnknown = "unknown"
	// suffix of secret name
	SuffixSecretName = "-userpassword-secret"
)
