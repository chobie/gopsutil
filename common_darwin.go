// +build darwin
package gopsutil

import (
	"syscall"
	"unsafe"
	"encoding/binary"
	"bytes"
)

func sysctlbyname(name string, data interface{}) (err error) {
	val, err := syscall.Sysctl(name)
	if err != nil {
		return err
	}

	buf := []byte(val)

	switch v := data.(type) {
	case *int32:
		*v = *(*int32)(unsafe.Pointer(&buf[0]))
		return
	case *int64:
		*v = *(*int64)(unsafe.Pointer(&buf[0]))
		return
	case *uint64:
		*v = *(*uint64)(unsafe.Pointer(&buf[0]))
		return
	case *float64:
		*v = *(*float64)(unsafe.Pointer(&buf[0]))
		return
	}

	bbuf := bytes.NewBuffer([]byte(val))
	return binary.Read(bbuf, binary.LittleEndian, data)
}
