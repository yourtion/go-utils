package ipc

import (
	"os"
	"strconv"
	"syscall"
)

// SocketPair create a socket pair
func SocketPair() ([]*os.File, error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	return []*os.File{
		os.NewFile(uintptr(fds[0]), "@/nodejs/left"+strconv.Itoa(fds[0])),
		os.NewFile(uintptr(fds[1]), "@/nodejs/right"+strconv.Itoa(fds[1])),
	}, err
}
