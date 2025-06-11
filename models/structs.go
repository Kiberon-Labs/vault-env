package models

//SecretCtx : The inputs we receive from the CLI that need to be  accessible
type SecretCtx struct {
	// VaultAddress is the address of the Vault server
	VaultAddress string
	// InputFile is the path to the input file containing the secrets configuration
	InputFile string
	// Collection is the name of the collection of secrets to retrieve
	Collection string
	// Output is the path to the output file where secrets will be written
	Output string
	// DefaultToken is the token to use for authentication with Vault
	DefaultToken string
	// Format is the output format for the secrets, e.g., NIX or WINDOWS
	Format string
}

type SecretsRoot struct {
	Secrets Secrets `yaml:"secrets,omitempty"`
}

type Secrets struct {
	Version    *string      `yaml:"version"`
	Collection []Collection `yaml:"collections,omitempty"`
	/**
	 * Output is the default output file for all secrets
	 * If a secret has a file specified, it will override this value
	 */
	Output *string `yaml:"output,omitempty"`
}

// A Collections is a named  group of secrets.
// It can be used to logically group secrets together, such as by application or environment.
type Collection struct {
	Name   *string  `yaml:"name"`
	Values []Secret `yaml:"values"`
}

// An Alias is an optional field on a Secret that allows you to map a secret field to a different name.
type Alias struct {
	Field *string `yaml:"field"`
	Name  *string `yaml:"name"`
}

// A secret defines an individual secret from Vault along with metadata about how to retrieve it and process it.
type Secret struct {
	Engine *string `yaml:"engine,omitempty"`
	Root   *string `yaml:"root,omitempty"`
	Path   *string `yaml:"path,omitempty"`
	Field  *string `yaml:"field,omitempty"`
	// ENV or FILE
	Type SecretOutputType `yaml:"type,omitempty"`
	File *string          `yaml:"file,omitempty"`
	//Only usable by kv-v2, allows specifying a specific version of a secret. Omitting this will retrieve the latest version.
	Version *int `yaml:"version,omitempty"`
	/**
	 * Aliases are used to map a secret to a different name
	 * IE if you have a secret named `password` but want to use it as `PASSWORD`
	 */
	Aliases []Alias `yaml:"aliases,omitempty"`
	//  A specific namespace to use for the secret.
	Namespace *string `yaml:"namespace,omitempty"`
}
