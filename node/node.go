package node

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
)

// Node is the expression of nodes
type Node interface {
	nodefs.Node
	GetName() string
	IsDir() bool
}

// FetchFileFunc is the type for provide how file is dowloaded
type FetchFileFunc func(f *FileNode) error

// FileNode is the expression of files
type FileNode struct {
	nodefs.Node
	Path      string
	Name      string
	FetchFile FetchFileFunc
	File      File
}

// GetName returns the name of the file
func (file *FileNode) GetName() string {
	return file.Name
}

// IsDir returns false
func (file *FileNode) IsDir() bool {
	return false
}

// File is file
type File struct {
	nodefs.File

	Content []byte
	Size    uint64
}

// Read reads from file
func (f *File) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	return fuse.ReadResultData(f.Content), fuse.OK
}

// Open returns file
func (file *FileNode) Open(flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	err := file.FetchFile(file)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	return nodefs.NewDataFile(file.File.Content), fuse.OK
}

// GetAttr get attr
func (file *FileNode) GetAttr(out *fuse.Attr, f nodefs.File, ctx *fuse.Context) fuse.Status {
	out.Size = file.File.Size
	out.Mode = fuse.S_IFREG | 0444
	out.Ctime = uint64(time.Now().Unix())
	out.Mtime = uint64(time.Now().Unix())

	return fuse.OK
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
func NewFileNode(path string, name string, f FetchFileFunc) *FileNode {
	return &FileNode{
		Node:      nodefs.NewDefaultNode(),
		Path:      path,
		Name:      name,
		FetchFile: f,
		File:      File{File: nodefs.NewDefaultFile(), Size: 0},
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
