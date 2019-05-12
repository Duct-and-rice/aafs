package node

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

// Node is the expression of nodes
type Node interface {
	nodefs.Node
	GetName() string
	IsDir() bool
}

// FileNode is the expression of files
type FileNode struct {
	nodefs.Node
	Path string
	Name string
}

// GetName returns the name of the file
func (file *FileNode) GetName() string {
	return file.Name
}

// IsDir returns false
func (file *FileNode) IsDir() bool {
	return false
}

// DirNode is the expression of directories
type DirNode struct {
	nodefs.Node
	Children []Node
	Name     string
}

// GetName returns the name of the directory
func (dir *DirNode) GetName() string {
	return dir.Name
}

// IsDir returns true
func (dir *DirNode) IsDir() bool {
	return true
}

// OpenDir returns dir entry
func (dir *DirNode) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	p := dir.Inode()
	children := dir.Children
	entries := make([]fuse.DirEntry, len(children))
	for i, child := range children {
		entries[i] = fuse.DirEntry{
			Mode: 444,
			Name: child.GetName(),
		}
		p.NewChild(child.GetName(), child.IsDir(), child)
	}
	return entries, fuse.OK
}

// NewFileNode returns new file node
func NewFileNode(path string, name string) *FileNode {
	return &FileNode{
		Node: nodefs.NewDefaultNode(),
		Path: path,
		Name: name,
	}
}

// NewDirNode returns new dir node
func NewDirNode(children []Node, name string) *DirNode {
	return &DirNode{
		Node:     nodefs.NewDefaultNode(),
		Children: children,
		Name:     name,
	}
}
