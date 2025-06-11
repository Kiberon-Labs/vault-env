//go:build unit_test
// +build unit_test

package secrets_test

import (
	"os"
	"path/filepath"
	"strings"

	"kiberon-labs/vault-env/v2/secrets"

	"kiberon-labs/vault-env/v2/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secrets", func() {

	Describe("Environment support", func() {

		It("Injects the correct environment values", func() {

			name := "TEST_VAL"
			value := "rAFZDFAFEWRAGDSWCF"
			os.Setenv(name, value)
			absPath, _ := filepath.Abs("../tests/.secrets.yml")
			val, err := secrets.GetYamlConfig(absPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Contains(string(val), value)).To(Equal(true))
		})
	})
})
