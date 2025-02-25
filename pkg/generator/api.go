package generator

import (
	"github.com/cidverse/go-vcsapp/pkg/platform/api"
)

type Config struct {
	Directory  string
	Platform   api.Platform
	Repository api.Repository
}

// Generator provides a common interface for all generators
type Generator interface {
	Name() string
	Generate() error
	SetOutputDirectory(string)
	GetOutputDirectory() string
}
