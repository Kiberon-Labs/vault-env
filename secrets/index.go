package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kiberon-labs/vault-env/v2/models"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type secretVal struct {
	Name   string
	Value  string
	config models.Secret
}

func getSecret(secret models.Secret, client *api.Client) ([]secretVal, error) {

	if *secret.Engine == "kv-v2" {
		return getKV2(secret, client)
	}

	return getNormalSecret(secret, client)
}

func processSingleValue(output string, secret models.Secret) secretVal {

	secretName := strings.ToUpper(*secret.Field)

	//Check if there is an alias
	if len(secret.Aliases) > 0 {
		secretName = *(secret.Aliases[0].Name)
	}

	return secretVal{
		Name:   secretName,
		Value:  output,
		config: secret,
	}
}

func processMultiValue(keyVals map[string]interface{}, secret models.Secret) []secretVal {

	retArr := make([]secretVal, 0)
	aliasLookup := make(map[string]models.Alias)

	for _, a := range secret.Aliases {
		aliasLookup[*a.Field] = a
	}

	for key, rawValue := range keyVals {

		value := fmt.Sprintf("%v", rawValue)

		newName := strings.ToUpper(key)

		if *secret.Engine == "aws" {

			switch newName {
			//Rename to expected so we can use on CLI immediately
			case "ACCESS_KEY":
				newName = "AWS_ACCESS_KEY_ID"
			case "SECRET_KEY":
				newName = "AWS_SECRET_ACCESS_KEY"
			case "SECURITY_TOKEN":
				newName = "AWS_SESSION_TOKEN"
			//Suppress useless info
			case "LEASE_ID":
				fallthrough
			case "LEASE_DURATION":
				fallthrough
			case "LEASE_RENEWABLE":
				continue
			}

			if newName == "AWS_SESSION_TOKEN" && value == "<nil>" {
				continue
			}
		}

		candidate, exists := aliasLookup[newName]

		if exists {
			newName = *candidate.Name
		}

		//Lookup aliasing if necessary

		retArr = append(retArr, secretVal{
			Name:   newName,
			Value:  value,
			config: secret,
		})
	}
	return retArr
}

// This is specifically written to not return the meta data
func getKV2(secret models.Secret, client *api.Client) ([]secretVal, error) {

	secretPath := formatPath(*secret.Engine, *secret.Root, *secret.Path)
	additionalData := make(map[string][]string)

	if secret.Version != nil {
		additionalData["version"] = []string{fmt.Sprintf("%d", *secret.Version)}
	}

	output, err := client.Logical().ReadWithData(secretPath, additionalData)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, fmt.Errorf("No secret detected at %s", secretPath)
	}

	if output.Data == nil || output.Data["data"] == nil {
		return nil, fmt.Errorf("Expected data not found at %s", secretPath)
	}

	data := output.Data["data"].(map[string]interface{})

	//Single value request
	if secret.Field != nil {
		singleStringSecret := fmt.Sprintf("%v", data[*secret.Field])
		return []secretVal{processSingleValue(singleStringSecret, secret)}, nil
	}

	return processMultiValue(data, secret), nil
}

func getNormalSecret(secret models.Secret, client *api.Client) ([]secretVal, error) {

	var root = ""

	secretPath := formatPath(*secret.Engine, root, *secret.Path)
	additionalData := make(map[string][]string)
	if secret.Field != nil {
		additionalData["field"] = []string{*secret.Field}
	}

	output, err := client.Logical().ReadWithData(secretPath, additionalData)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if output == nil {
		return nil, fmt.Errorf("No secret detected at %s", secretPath)
	}

	//Single value request
	if secret.Field != nil {
		singleStringSecret := fmt.Sprintf("%v", output.Data[*secret.Field])
		return []secretVal{processSingleValue(singleStringSecret, secret)}, nil
	}

	return processMultiValue(output.Data, secret), nil

}

func setupVaultClient(ctx *models.SecretCtx) (*api.Client, error) {

	clientConfig := api.DefaultConfig()
	clientConfig.Address = ctx.VaultAddress

	if ctx.VaultAddress == "" {
		return nil, fmt.Errorf("Vault address is not set")
	}

	client, err := api.NewClient(clientConfig)

	if err != nil {
		return nil, fmt.Errorf("Error creating vault client: %v", err)
	}
	return client, nil
}

func RetrieveSecrets(ctx *models.SecretCtx) error {

	client, err := setupVaultClient(ctx)
	if err != nil {
		return err
	}

	secrets, err := getConfig(ctx.InputFile)

	if err != nil {
		return err
	}

	returnedSecrets := make([]secretVal, 0)

	var foundCollection *models.Collection

	for _, v := range secrets.Collection {
		if *v.Name == ctx.Collection {
			foundCollection = &v
			break
		}
	}

	if foundCollection == nil {
		log.Trace("No applicable secrets detected")
		return nil
	}

	//ROLE / TOKEN lookup
	client.SetToken(ctx.DefaultToken)

	log.Debugln(fmt.Sprintf("Found %d secret(s) to retrieve", len(foundCollection.Values)))
	for _, v := range foundCollection.Values {

		if v.Namespace != nil {
			client.SetNamespace(*v.Namespace)
		}

		vals, err := getSecret(v, client)
		if err != nil {
			log.Println("Could not retrieve secret")
			return err
		}

		returnedSecrets = append(returnedSecrets, vals...)
	}

	fsLookup := make(map[string]*os.File)

	for _, v := range returnedSecrets {

		if v.config.Type == models.EnvironmentSecret {

			if ctx.Format == string(models.NIX) {
				fmt.Printf("%s=\"%s\"\n", v.Name, v.Value)
			} else {
				fmt.Printf("set %s=\"%s\"\n", v.Name, v.Value)
			}
		} else {

			var pathVal string

			if len(ctx.Output) > 0 {
				pathVal = ctx.Output
			} else if (v.config.File != nil) && (len(*v.config.File) > 0) {
				pathVal = *v.config.File
			} else {
				pathVal = models.DefaultDotEnvOutput
			}

			if !filepath.IsAbs(pathVal) {
				pathVal = filepath.Join(filepath.Dir(ctx.InputFile), pathVal)
			}

			f, exists := fsLookup[pathVal]
			if !exists {
				f, err = os.OpenFile(pathVal, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0440)
				if err != nil {
					log.Println("Failed to open file")
					return err
				}
				defer f.Close()
				fsLookup[pathVal] = f
			}

			if len(v.Name) != 0 {

				var outputString string
				if ctx.Format == string(models.NIX) {
					outputString = fmt.Sprintf("%s=\"%s\"\n", v.Name, v.Value)
				} else {
					outputString = fmt.Sprintf("set %s=\"%s\"\n", v.Name, v.Value)
				}

				if _, err = f.WriteString(outputString); err != nil {
					log.Println("Failed to write to file")
					return err
				}
			}
		}
	}

	return nil
}
