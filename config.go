package storagenomad

const (
	// DefaultPrefix defines the default prefix in secure variable store
	DefaultPrefix = "caddytls"

	// DefaultValuePrefix sets a prefix to secure variables to check validation
	DefaultValuePrefix = "caddy-storage-nomad"

	// DefaultTimeout is the default timeout for Nomad connections
	DefaultTimeout = 10

	// EnvNamePrefix defines the env variable name to override Secure Var key prefix
	EnvNamePrefix = "CADDY_CLUSTERING_NOMAD_PREFIX"

	// EnvValuePrefix defines the env variable name to override Secure Var value prefix
	EnvValuePrefix = "CADDY_CLUSTERING_NOMAD_VALUEPREFIX"
)
