package model

// ApplicationEnvironment is a type for accessing the server environment
type ApplicationEnvironment string

const (
	// ApplicationEnvironmentLocal environment for the loca environment
	ApplicationEnvironmentLocal ApplicationEnvironment = "local"
	// ApplicationEnvironmentDev environment for the dev environment
	ApplicationEnvironmentDev ApplicationEnvironment = "dev"
	// ApplicationEnvironmentStaging environment for the staging environment
	ApplicationEnvironmentStaging ApplicationEnvironment = "staging"
	// ApplicationEnvironmentProduction environment for the production environment
	ApplicationEnvironmentProduction ApplicationEnvironment = "production"
)
