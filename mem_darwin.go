// +build darwin
package gopsutil

/*
#include <stdio.h>
#include <stdlib.h>
#include <mach/mach_init.h>
#include <mach/mach_host.h>
 */
import "C"

import (
	"fmt"
	"unsafe"
)

func vm_info(vmstat *C.vm_statistics_data_t) error {
	var count C.mach_msg_type_number_t = C.HOST_VM_INFO_COUNT
	status := C.host_statistics(
		C.host_t(C.mach_host_self()),
		C.HOST_VM_INFO,
		C.host_info_t(unsafe.Pointer(vmstat)),
		&count)

	if status != C.KERN_SUCCESS {
		return fmt.Errorf("host_statistics=%d", status)
	}

	return nil
}

func VirtualMemory() (*VirtualMemoryStat, error) {
	var vm_stat C.vm_statistics_data_t
	var total uint64

	if err := sysctlbyname("hw.memsize", &total); err != nil {
		return nil, err
	}

	err := vm_info(&vm_stat)
	if err != nil {
		return nil, err
	}

	kern := uint64(vm_stat.inactive_count) * uint64(C.vm_page_size)
	free := uint64(vm_stat.free_count) * uint64(C.vm_page_size)
	used := total - free
	ret := &VirtualMemoryStat{
		Total: total,
		Available: free + kern,
		Used: used,
		UsedPercent: float64(used)/ float64(total),
		Free: free,
		Active: uint64(vm_stat.active_count) * uint64(C.vm_page_size),
		Inactive: uint64(vm_stat.inactive_count) * uint64(C.vm_page_size),
		//Buffers
		//Cached
		Wired: uint64(vm_stat.wire_count) * uint64(C.vm_page_size),
	}

	return ret, nil
}

type xsw_usage struct {
	Total, Avail, Used uint64
}

func SwapMemory() (*SwapMemoryStat, error) {
	sw_usage := xsw_usage{}
	if err := sysctlbyname("vm.swapusage", &sw_usage); err != nil {
		return nil, err
	}

	ret := &SwapMemoryStat{
		Total: sw_usage.Total,
		Used: sw_usage.Used,
		Free: sw_usage.Avail,
		UsedPercent: float64(sw_usage.Used) / float64(sw_usage.Total),
		//SIN
		//SOUT
	}
	return ret, nil
}
