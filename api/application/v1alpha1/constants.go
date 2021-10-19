package v1alpha1

const (
	// custom resource status
	Creating      = "Creating"
	Failed        = "Failed"
	Error         = "Error"
	AlreadyExists = "AlreadyExists"

	// kind of operator cr
	DBTypeMysql      = "MySQL"

	// kind of cluster
	KindPostgreSQLCluster = "PostgreSQLCluster"
	KindPgCluster = "Pgcluster"
	KindPgClusterVersion = "Pgcluster"

	// cluster status
	ClusterStatusUnknown = "unknown"
	// suffix of secret name
	SuffixSecretName = "-userpassword-secret"
)
