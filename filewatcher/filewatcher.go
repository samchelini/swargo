package filewatcher

import (
	"syscall"
)

const (
	// inotify watch and read masks
	IN_ACCESS        = syscall.IN_ACCESS
	IN_ATTRIB        = syscall.IN_ATTRIB
	IN_CLOSE_WRITE   = syscall.IN_CLOSE_WRITE
	IN_CLOSE_NOWRITE = syscall.IN_CLOSE_NOWRITE
	IN_CREATE        = syscall.IN_CREATE
	IN_DELETE        = syscall.IN_DELETE
	IN_DELETE_SELF   = syscall.IN_DELETE_SELF
	IN_MODIFY        = syscall.IN_MODIFY
	IN_MOVE_SELF     = syscall.IN_MOVE_SELF
	IN_MOVED_FROM    = syscall.IN_MOVED_FROM
	IN_MOVED_TO      = syscall.IN_MOVED_TO
	IN_OPEN          = syscall.IN_OPEN

	// inotify watch macros
	IN_ALL_EVENTS = syscall.IN_ALL_EVENTS
	IN_MOVE       = syscall.IN_MOVE
	IN_CLOSE      = syscall.IN_CLOSE

	// additional inotify watch masks
	IN_DONT_FOLLOW = syscall.IN_DONT_FOLLOW
	IN_EXCL_UNLINK = syscall.IN_EXCL_UNLINK
	IN_MASK_ADD    = syscall.IN_MASK_ADD
	IN_ONESHOT     = syscall.IN_ONESHOT
	IN_ONLYDIR     = syscall.IN_ONLYDIR
	IN_MASK_CREATE = 0x10000000

	// inotify read masks
	IN_IGNORED    = syscall.IN_IGNORED
	IN_ISDIR      = syscall.IN_ISDIR
	IN_Q_OVERFLOW = syscall.IN_Q_OVERFLOW
	IN_UNMOUNT    = syscall.IN_UNMOUNT
)

type FileWatcher struct {
	fd  int // inotify file descriptor
	buf [syscall.SizeofInotifyEvent]byte
}

// initialize inotify and return new FileWatcher instance
func NewFileWatcher() (*FileWatcher, error) {
	fd, err := syscall.InotifyInit()
	return &FileWatcher{fd: fd}, err
}

// add file to the watch list for the specified event
func (fw *FileWatcher) AddWatch(file string, event uint32) error {
	_, err := syscall.InotifyAddWatch(fw.fd, file, event)
	return err
}

// watch for events
func (fw *FileWatcher) Watch() error {
	_, err := syscall.Read(fw.fd, fw.buf[:])
	return err
}

// close inotify file desciptor
func (fw *FileWatcher) Close() error {
	err := syscall.Close(fw.fd)
	return err
}
