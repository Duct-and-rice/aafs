package fs

import (
	"github.com/Duct-and-rice/aafs/provider"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

// Root is the expression of the root directory
type Root struct {
	nodefs.Node

	Provider provider.Provider
}

// OnMount downloads informations from provider and set them
func (root *Root) OnMount(conn *nodefs.FileSystemConnector) {
	err := root.Provider.FetchFiles()
	if err != nil {
		return
	}
	root.Node.OnMount(conn)
}

// OpenDir returns directory files
func (root *Root) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	p := root.Inode()
	children := root.Provider.Files()
	entries := make([]fuse.DirEntry, len(children))
	for i, child := range children {
		name := child.GetName()
		entries[i] = fuse.DirEntry{
			Mode: 444,
			Name: name,
		}
		if p.GetChild(name) == nil {
			p.NewChild(name, child.IsDir(), child)
		}
	}
	return entries, fuse.OK
}

// NewRoot returns a new root instance
func NewRoot(providerName string) *Root {

	return &Root{
		Node:     nodefs.NewDefaultNode(),
		Provider: provider.NewProvider(providerName),
	}
}
