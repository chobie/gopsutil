// +build darwin

package gopsutil

/*
#include <utmpx.h>
 */
import "C"

import (
	"runtime"
	"syscall"
	"time"
	"os"
	"os/exec"
	"strings"
)

func HostInfo() (*HostInfoStat, error) {
	tv := syscall.Timeval{}
	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		return nil, err
	}

	uptime := time.Since(time.Unix(tv.Unix())).Seconds()
	info := &HostInfoStat{
		OS:             runtime.GOOS,
		Platform: "darwin",
		PlatformFamily: "darwin",
		Uptime: uint64(uptime),
	}

	hostname, err := os.Hostname()
	if err != nil {
		return info, err
	}
	info.Hostname = hostname

	version, err := syscall.Sysctl("kern.osrelease")
	if err != nil {
		return info, err
	}
	info.PlatformVersion = version

	return info, nil
}

func BootTime() (int64, error) {
	var boottime int64
	if err := sysctlbyname("kern.boottime", &boottime); err != nil {
		return 0, err
	}
	return boottime, nil
}

func Users() ([]UserStat, error) {
	var user *C.struct_utmpx = nil

	stat := []UserStat{}
	C.setutxent()
	defer C.endutxent()
	// TODO(chobie): Is this correct?
	for user = C.getutxent(); user != nil; user = C.getutxent() {
		stat = append(stat, UserStat{
			User: C.GoString(&user.ut_user[0]),
			Terminal: C.GoString(&user.ut_line[0]),
			Host: C.GoString(&user.ut_host[0]),
			Started: int(user.ut_tv.tv_sec),
		})
	}

	return stat, nil
}

func GetPlatformInformation() (string, string, string, error) {
	platform := ""
	family := ""
	version := ""

	out, err := exec.Command("uname", "-s").Output()
	if err == nil {
		platform = strings.ToLower(strings.TrimSpace(string(out)))
	}

	out, err = exec.Command("uname", "-r").Output()
	if err == nil {
		version = strings.ToLower(strings.TrimSpace(string(out)))
	}

	return platform, family, version, nil
}

func GetVirtualization() (string, string, error) {
	system := ""
	role := ""

	return system, role, nil
}
