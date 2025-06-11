package models

// Constants for the project
const (
	DefaultDotEnvOutput = "./.env.vault"
)

type SecretOutputType string

const (
	EnvironmentSecret SecretOutputType = "ENV"
	FILE              SecretOutputType = "FILE"
)

type OSFormat string

const (
	WINDOWS OSFormat = "WINDOWS"
	NIX     OSFormat = "FILE"
)
