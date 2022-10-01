package p2p

import "golang.org/x/sys/unix"

func CheckDirectoryFreeSpace(dir string) uint64 {
	var stat unix.Statfs_t
	unix.Statfs(dir, &stat)
	return stat.Bavail * uint64(stat.Bsize)
}
