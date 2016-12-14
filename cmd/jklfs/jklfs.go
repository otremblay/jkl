package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"otremblay.com/jkl"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type jklfs struct {
	pathfs.FileSystem
	issuePerDirs map[string]*jkl.JiraIssue
}

func (j *jklfs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch name {
	case "current_sprint":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	if _, ok := j.issuePerDirs[name]; ok {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	pathPieces := strings.Split(name, "/")
	path := strings.Join(pathPieces[0:2], "/")
	if i, ok := j.issuePerDirs[path]; ok {
		if path+"/description" == name {
			return &fuse.Attr{
				Mode: fuse.S_IFREG | 0644, Size: uint64(len(i.Fields.Description)),
			}, fuse.OK
		}
	}
	return nil, fuse.ENOENT
}

func (j *jklfs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c = []fuse.DirEntry{{Name: "current_sprint", Mode: fuse.S_IFDIR}}
		return c, fuse.OK
	}
	if name == "current_sprint" {
		issues, err := jkl.List("sprint in openSprints() and project = 'DO'")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, fuse.ENOENT
		}
		c = make([]fuse.DirEntry, len(issues))
		for i, issue := range issues {
			c[i] = fuse.DirEntry{Name: issue.Key, Mode: fuse.S_IFDIR}
			j.issuePerDirs["current_sprint/"+issue.Key] = issue
		}
		return c, fuse.OK
	}

	if _, ok := j.issuePerDirs[name]; ok {
		c = []fuse.DirEntry{fuse.DirEntry{Name: "description", Mode: fuse.S_IFREG}}
		return c, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (j *jklfs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	pathPieces := strings.Split(name, "/")
	path := strings.Join(pathPieces[0:2], "/")
	if i, ok := j.issuePerDirs[path]; ok {
		if path+"/description" == name {
			return nodefs.NewDataFile([]byte(i.Fields.Description)), fuse.OK
		}
	}
	return nil, fuse.ENOENT
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  jklfs MOUNTPOINT")
	}
	nfs := pathfs.NewPathNodeFs(&jklfs{pathfs.NewDefaultFileSystem(), map[string]*jkl.JiraIssue{}}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
