//go:build !windows

package utils

import (
	"golang.org/x/sys/unix"
)

func CheckDirectoryFreeSpace(dir string) uint64 {
	var stat unix.Statfs_t
	unix.Statfs(dir, &stat)
	return stat.Bavail * uint64(stat.Bsize)
}

// Input `size` in bytes
func GetBytesFromMBytes(size uint64) uint64 {
	return size * 1024 * 1024
}

func GetMBytesFromBytes(size int64) int64 {
	return size / 1024 / 1024
}
func GetKBytesFromBytes(size int64) int64 {
	return size / 1024
}

func IsSpaceEnough(dir string, size uint64) bool {
	return CheckDirectoryFreeSpace(dir)-size > 0
}
