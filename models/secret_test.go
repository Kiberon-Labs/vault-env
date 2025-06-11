package models

import (
	"testing"
)

type dummySecretOutputType string

func TestSecretInitDefaultsDefaultTypeAndFile(t *testing.T) {
	output := "default.txt"
	secret := &Secret{
		Type: "",
		File: nil,
	}
	secrets := &Secrets{
		Output: &output,
	}
	got, err := secret.InitDefaults(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != FILE {
		t.Errorf("expected type %v, got %v", FILE, got.Type)
	}
	if got.File != &output {
		t.Errorf("expected file output to be default, got %v", got.File)
	}
}

func TestSecretInitDefaultsExplicitTypeAndFile(t *testing.T) {
	file := "myfile.txt"
	secret := &Secret{
		Type: FILE,
		File: &file,
	}
	secrets := &Secrets{}
	got, err := secret.InitDefaults(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != FILE {
		t.Errorf("expected type %v, got %v", FILE, got.Type)
	}
	if got.File != &file {
		t.Errorf("expected file output to be %v, got %v", file, got.File)
	}
}
