package main

import (
	"fmt"
	"os"
	"time"
)
import "github.com/hanwen/go-fuse/fuse"
import "github.com/hanwen/go-fuse/fuse/nodefs"
import "io/ioutil"

func NewJklfsFile() (nodefs.File, error) {
	f, err := ioutil.TempFile("", "jklfile")
	if err != nil {
		return nil, err
	}
	return &jklfile{f}, nil
}

type jklfile struct {
	*os.File
}

func (f *jklfile) InnerFile() nodefs.File {
	return nil
}

func (f *jklfile) String() string {
	return fmt.Sprintf("jklfile(%s)", f.Name())
}

func (f *jklfile) Write(data []byte, off int64) (uint32, fuse.Status) {
	n, err := f.File.WriteAt(data, off)
	if err != nil {
		return fuse.EACCES
	}
	return uint32(n), fuse.OK
}

func (f *jklfile) Fsync(flag int) (code fuse.Status) {
	return fuse.OK
}

func (f *jklfile) Truncate(size uint64) fuse.Status {
	return fuse.EPERM
}

func (f *jklfile) Chmod(mode uint32) fuse.Status {
	return fuse.EPERM
}

func (f *jklfile) Chown(uid uint32, gid uint32) fuse.Status {
	return fuse.EPERM
}

func (f *jklfile) Allocate(off uint64, sz uint64, mode uint32) fuse.Status {
	return fuse.EPERM
}

func (f *jklfile) Flush() fuse.Status {
	return fuse.OK
}

func (f *jklfile) GetAttr(out *fuse.Attr) fuse.Status {
	return fuse.OK
}

func (f *jklfile) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	return nil, fuse.OK
}
func (f *jklfile) Release() {

}
func (f *jklfile) SetInode(i *nodefs.Inode) {}

func (f *jklfile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	return fuse.EPERM
}
