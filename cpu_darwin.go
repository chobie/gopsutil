// +build darwin

package gopsutil

/*
#include <mach/mach.h>
#include <mach/mach_error.h>
#include <mach/mach_host.h>
#include <stdlib.h>
 */
import "C"

import (
	"encoding/binary"
	"unsafe"
	"bytes"
	"fmt"
	"syscall"
)

// TODO: get per cpus
func CPUTimes(percpu bool) ([]CPUTimesStat, error) {
	var count C.mach_msg_type_number_t
	var cpuload *C.processor_cpu_load_info_data_t
	var ncpu C.natural_t

	status := C.host_processor_info(C.host_t(C.mach_host_self()),
		C.PROCESSOR_CPU_LOAD_INFO,
		&ncpu,
		(*C.processor_info_array_t)(unsafe.Pointer(&cpuload)),
		&count)

	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_processor_info error=%d", status)
	}

	// jump through some cgo casting hoops and ensure we properly free
	// the memory that cpuload points to
	target := C.vm_map_t(C.mach_task_self_)
	address := C.vm_address_t(uintptr(unsafe.Pointer(cpuload)))
	defer C.vm_deallocate(target, address, C.vm_size_t(ncpu))

	// the body of struct processor_cpu_load_info
	// aka processor_cpu_load_info_data_t
	var cpu_ticks [C.CPU_STATE_MAX]uint32

	// copy the cpuload array to a []byte buffer
	// where we can binary.Read the data
	size := int(ncpu) * binary.Size(cpu_ticks)
	buf := C.GoBytes(unsafe.Pointer(cpuload), C.int(size))
	bbuf := bytes.NewBuffer(buf)

	ret := []CPUTimesStat{}
	err := binary.Read(bbuf, binary.LittleEndian, &cpu_ticks)
	if err != nil {
		return nil, err
	}

	user := float32(cpu_ticks[C.CPU_STATE_USER])
	system := float32(cpu_ticks[C.CPU_STATE_SYSTEM])
	idle := float32(cpu_ticks[C.CPU_STATE_IDLE])
	nice := float32(cpu_ticks[C.CPU_STATE_NICE])
	//total := user + system + idle + nice
	cpu := CPUTimesStat{
		User: user,
		System: system ,
		Idle: idle,
		Nice: nice,
	}
	ret = append(ret, cpu)

	return ret, nil
}

// Returns only one CPUInfoStat on FreeBSD
func CPUInfo() ([]CPUInfoStat, error) {
	var host_info C.host_basic_info_data_t
	var count C.mach_msg_type_number_t = C.HOST_VM_INFO_COUNT
	var cpu_type_name, cpu_subtype_name *C.char;
	var freq int64
	var cache_size int32

	status := C.host_info(C.host_t(C.mach_host_self()),
		C.HOST_BASIC_INFO,
		C.host_info_t(unsafe.Pointer(&host_info)),
		&count)

	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_processor_info error=%d", status)
	}
	fmt.Printf("cpuinfo: %+v\n", host_info)
	C.slot_name(host_info.cpu_type, host_info.cpu_subtype, &cpu_type_name,
		&cpu_subtype_name);
	fmt.Printf("name: %s, name: %s\n", C.GoString(cpu_type_name), C.GoString(cpu_subtype_name))

	ret := []CPUInfoStat{}
	brand, _ := syscall.Sysctl("machdep.cpu.brand_string")
	sysctlbyname("hw.cpufrequency", &freq)
	sysctlbyname("machdep.cpu.cache.size", &cache_size)

	for i := 0; i < int(host_info.max_cpus); i++ {
		info := CPUInfoStat{
			CPU: int32(i),
			ModelName: brand,
			Mhz: float64(freq),
			CacheSize: cache_size,
		}
		ret = append(ret, info)
	}
	return ret, nil
}
