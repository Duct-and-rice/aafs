package aahub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Duct-and-rice/aafs/node"
	"golang.org/x/text/encoding/japanese"
)

// Provider is a provider for aahub
type Provider struct {
	Folders []Folder
}

// NewProvider returns new provider
func NewProvider() *Provider {
	folders := []Folder{}
	return &Provider{
		Folders: folders,
	}
}

// Folder is folder of aahub
type Folder struct {
	Kind    int      `json:"kind"`
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Folders []Folder `json:"folders"`
}

type aahubFolderRoot struct {
	Folders []Folder `json:"folders"`
}

// Name returns aahub name
func (provider *Provider) Name() string {
	return "aahub"
}

// FetchFiles get directory from aahub
func (provider *Provider) FetchFiles() error {
	fmt.Printf("start fetch directories\n")
	url := "https://aa-storage.aahub.org/folders.json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	data := new(aahubFolderRoot)

	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}

	provider.Folders = data.Folders
	fmt.Printf("end fetch directories\n")
	return nil
}

// Files returns files as []node.Node
func (provider *Provider) Files() []node.Node {
	folders := provider.Folders
	result := make([]node.Node, len(folders))
	for i, folder := range folders {
		result[i] = folder.File()
	}
	return result
}

// File returns folder to files
func (f *Folder) File() node.Node {
	if len(f.Folders) == 0 {
		// if len of f.Folders, f is a a file
		return node.NewFileNode(f.Path, f.Name+".mlt", fetchFile)
	}

	children := make([]node.Node, len(f.Folders))
	for i, child := range f.Folders {
		children[i] = child.File()
	}

	return node.NewDirNode(children, f.Name)
}

type aahubMLT struct {
	Name    string     `json:"name,omitempty"`
	Path    string     `json:"path,omitempty"`
	Updated string     `json:"updated,omitempty"`
	AAs     []*aahubAA `json:"aa,omitempty"`
}
type aahubAA struct {
	Value string   `json:"value,omitempty"`
	Tag   []string `json:"tag,omitempty"`
}

var memo = map[string][]byte{}

func fetchFile(file *node.FileNode) error {
	if memo[file.Path] != nil {
		file.File.Content = memo[file.Path]
		return nil
	}
	url := "https://aa-storage.aahub.org/folders/" + file.Path + ".json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := new(aahubMLT)
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}

	b := strings.Builder{}

	encoder := japanese.ShiftJIS.NewEncoder()

	for i, aa := range data.AAs {
		b.WriteString(aa.Value)
		if i < len(data.AAs)-1 {
			b.WriteString("\n[SPLIT]\n")
		}
	}
	result, err := encoder.String(strings.Replace(b.String(), "\n", "\r\n", -1))
	if err != nil {
		return err
	}

	memo[file.Path] = []byte(result)

	file.File.Content = memo[file.Path]
	file.File.Size = uint64(len(memo[file.Path]))

	if len(memo) > 128 {
		for i := range memo {
			memo[i] = nil
		}
	}
	return nil
}
