package main

import (
	"log"
	"os"

	"github.com/Duct-and-rice/aafs/fs"
	nodefs "github.com/hanwen/go-fuse/fuse/nodefs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("aafs", "HukuTemp as FS")
	provider   = app.Flag("provider", "Provider").Short('p').Default("aahub").Enum("aahub", "yaruyomi")
	mountpoint = app.Flag("mountpoint", "Mountpoint").Short('m').Default("/mnt/aafs").String()
	debug      = app.Flag("debug", "is debug").Short('d').Bool()
)

func main() {
	app.Parse(os.Args[1:])

	root := fs.NewRoot(*provider)

	opts := nodefs.Options{}
	opts.Debug = *debug
	s, _, err := nodefs.MountRoot(*mountpoint, root, &opts)
	if err != nil {
		log.Fatal("Mount Error:", err)
	}
	s.Serve()
}
