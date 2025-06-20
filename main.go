package main

import (
	"fmt"
	"os"

	"kiberon-labs/vault-env/v2/models"
	"kiberon-labs/vault-env/v2/secrets"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var semverVersion string = "0.0.0"

func main() {
	var input string
	var collection string
	var output string
	var vaultAddress string
	var defaultToken string
	var format string
	verbose := false

	log.SetLevel(log.FatalLevel)

	app := &cli.App{
		Name:                 "vault-env",
		EnableBashCompletion: true,
		Usage:                "used to retrieve Hashicorp vault secrets",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Value:       ".vault-env.yml",
				Usage:       "the input file to parse.",
				Destination: &input,
			},
			&cli.StringFlag{
				Name:        "token",
				EnvVars:     []string{"VAULT_TOKEN"},
				Usage:       "the default vault token to use when no login requested",
				Destination: &defaultToken,
			},
			&cli.StringFlag{
				Name:        "collection",
				Value:       "default",
				Usage:       "the collection of secrets to pull",
				Destination: &collection,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "",
				Usage:       "the default file output",
				Destination: &output,
			},
			&cli.StringFlag{
				Name:        "vaultAddr",
				EnvVars:     []string{"VAULT_ADDR"},
				Usage:       "the address of the vault instance",
				Destination: &vaultAddress,
			},
			&cli.StringFlag{
				Name:        "format",
				Value:       string(models.NIX),
				Usage:       "the output format to use. Use `WINDOWS` if you want to output in .bat compatible format",
				Destination: &format,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Value:       false,
				Usage:       "sets the output to be verbose",
				Destination: &verbose,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "retrieves version information",
				Flags: []cli.Flag{},
				Action: func(c *cli.Context) error {
					fmt.Println(semverVersion)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {

			if verbose {
				log.SetLevel(log.TraceLevel)
			}

			if vaultAddress == "" {
				return fmt.Errorf("vault address is required, please set the VAULT_ADDR environment variable or use the --vaultAddr flag")
			}

			err := secrets.RetrieveSecrets(&models.SecretCtx{
				VaultAddress: vaultAddress,
				InputFile:    input,
				Collection:   collection,
				DefaultToken: defaultToken,
				Output:       output,
				Format:       format,
			})
			if err != nil {
				return err
			}
			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Println("A fatal error has occurred")
		log.Fatalln(err)
	} else {
		os.Exit(0)
	}
}
