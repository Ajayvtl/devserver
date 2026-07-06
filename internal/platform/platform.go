package platform

import (
	"context"
	"runtime"
)

type Platform struct {
	ID      string   `yaml:"id" json:"id"`
	Name    string   `yaml:"name" json:"name"`
	Version string   `yaml:"version" json:"version"`
	Like    []string `yaml:"like,omitempty" json:"like,omitempty"`
}

type Detector interface {
	Detect(context.Context) (Platform, error)
}

type SystemDetector struct{}

func NewDetector() *SystemDetector {
	return &SystemDetector{}
}

func (d *SystemDetector) Detect(context.Context) (Platform, error) {
	return Platform{
		ID:   runtime.GOOS,
		Name: runtime.GOOS,
	}, nil
}
