//go:build windows

package utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func CheckDirectoryFreeSpace(dir string) uint64 {
	// var stat unix.Statfs_t
	// unix.Statfs(dir, &stat)
	// return stat.Bavail * uint64(stat.Bsize)

	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	var fbytes int64
	_, _, err := c.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(dir))),
		uintptr(unsafe.Pointer(&fbytes)), nil, nil)

	return uint64(fbytes)
}

// Input `size` in bytes
func GetMBytesFromBytes(size int64) int64 {
	return size / 1024 / 1024
}

func IsSpaceEnough(dir string, size uint64) bool {
	return CheckDirectoryFreeSpace(dir)-size > 0
}
