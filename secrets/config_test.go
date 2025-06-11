package secrets

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetConfig(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	t.Run("throws an error if no version found", func(t *testing.T) {
		_, err := getConfig(filepath.Join(dir, "../tests/secrets/version-missing.yml"))
		if err.Error() != "`version` must be passed" {
			t.Errorf("Expected error for missing version, got: %v", err)
		}
	})

	t.Run("throws an error if mismatched version found", func(t *testing.T) {
		_, err := getConfig(filepath.Join(dir, "../tests/secrets/version-wrong.yml"))
		if err.Error() != "Only major version `1` is supported" {
			t.Errorf("Expected error for missing version, got: %v", err)
		}
	})

}

func TestGetYamlConfig(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	t.Run("returns data with environment variables expanded", func(t *testing.T) {

		name := "TEST_VAL"
		value := "rAFZDFAFEWRAGDSWCF"
		os.Setenv(name, value)
		data, err := getYamlConfig(filepath.Join(dir, "../tests/secrets/expand.yml"))
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if data == nil || *data == "" {
			t.Error("Expected non-empty data, got empty string")
		}

		if !strings.Contains(*data, value) {
			t.Errorf("Expected data to contain %s, got: %s", value, *data)
		}
	})
}
