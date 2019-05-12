package provider

import (
	"github.com/Duct-and-rice/aafs/aahub"
	"github.com/Duct-and-rice/aafs/node"
)

// Provider is for aahub.org
type Provider interface {
	Name() string
	FetchFiles() error
	Files() []node.Node
}

// NewProvider returns name's provider
func NewProvider(name string) Provider {
	switch name {
	case "aahub":
		p := aahub.NewProvider()
		return p

	default:
		return nil
	}
}
