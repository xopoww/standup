package config

import (
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func Load(r io.Reader, dst any) error {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(dst)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func LoadFile(path string, dst any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	return Load(f, dst)
}
