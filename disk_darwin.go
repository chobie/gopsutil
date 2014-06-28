// +build darwin

package gopsutil

import "C"

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	MNT_WAIT   = 1
	MFSNAMELEN = 16 /* length of type name including null */
	MNAMELEN   = 88 /* size of on/from name bufs */
)

// sys/mount.h
const (
	MNT_RDONLY      = 0x00000001 /* read only filesystem */
	MNT_SYNCHRONOUS = 0x00000002 /* filesystem written synchronously */
	MNT_NOEXEC      = 0x00000004 /* can't exec from filesystem */
	MNT_NOSUID      = 0x00000008 /* don't honor setuid bits on fs */
	MNT_UNION       = 0x00000020 /* union with underlying filesystem */
	MNT_ASYNC       = 0x00000040 /* filesystem written asynchronously */
	MNT_SUIDDIR     = 0x00100000 /* special handling of SUID on dirs */
	MNT_SOFTDEP     = 0x00200000 /* soft updates being done */
	MNT_NOSYMFOLLOW = 0x00400000 /* do not follow symlinks */
	MNT_GJOURNAL    = 0x02000000 /* GEOM journal support enabled */
	MNT_MULTILABEL  = 0x04000000 /* MAC support for individual objects */
	MNT_ACLS        = 0x08000000 /* ACL support enabled */
	MNT_NOATIME     = 0x10000000 /* disable update of file access time */
	MNT_NOCLUSTERR  = 0x40000000 /* disable cluster read */
	MNT_NOCLUSTERW  = 0x80000000 /* disable cluster write */
	MNT_NFS4ACLS    = 0x00000010
)

type Statfs struct {
	FVersion     uint32           /* structure version number */
	FType        uint32           /* type of filesystem */
	FFlags       uint64           /* copy of mount exported flags */
	FBsize       uint64           /* filesystem fragment size */
	FIosize      uint64           /* optimal transfer block size */
	FBlocks      uint64           /* total data blocks in filesystem */
	FBfree       uint64           /* free blocks in filesystem */
	FBavail      int64            /* free blocks avail to non-superuser */
	FFiles       uint64           /* total file nodes in filesystem */
	FFfree       int64            /* free nodes avail to non-superuser */
	FSyncwrites  uint64           /* count of sync writes since mount */
	FAsyncwrites uint64           /* count of async writes since mount */
	FSyncreads   uint64           /* count of sync reads since mount */
	FAsyncreads  uint64           /* count of async reads since mount */
	FSpare       [10]uint64       /* unused spare */
	FNamemax     uint32           /* maximum filename length */
	FOwner       uint32           /* user that mounted the filesystem */
	FFsid        int32            /* filesystem id */
	FCharspare   [80]byte         /* spare string space */
	FFstypename  [MFSNAMELEN]byte /* filesystem type name */
	FMntfromname [MNAMELEN]byte   /* mounted filesystem */
	FMntonname   [MNAMELEN]byte   /* directory on which mounted */
}

func DiskPartitions(all bool) ([]DiskPartitionStat, error) {
	var ret []DiskPartitionStat

	// get length
	count, err := Getfsstat(nil, MNT_WAIT)
	if err != nil {
		return ret, err
	}

	fs := make([]Statfs, count)
	_, err = Getfsstat(fs, MNT_WAIT)

	for _, stat := range fs {
		opts := "rw"
		if stat.FFlags&MNT_RDONLY != 0 {
			opts = "ro"
		}
		if stat.FFlags&MNT_SYNCHRONOUS != 0 {
			opts += ",sync"
		}
		if stat.FFlags&MNT_NOEXEC != 0 {
			opts += ",noexec"
		}
		if stat.FFlags&MNT_NOSUID != 0 {
			opts += ",nosuid"
		}
		if stat.FFlags&MNT_UNION != 0 {
			opts += ",union"
		}
		if stat.FFlags&MNT_ASYNC != 0 {
			opts += ",async"
		}
		if stat.FFlags&MNT_SUIDDIR != 0 {
			opts += ",suiddir"
		}
		if stat.FFlags&MNT_SOFTDEP != 0 {
			opts += ",softdep"
		}
		if stat.FFlags&MNT_NOSYMFOLLOW != 0 {
			opts += ",nosymfollow"
		}
		if stat.FFlags&MNT_GJOURNAL != 0 {
			opts += ",gjounalc"
		}
		if stat.FFlags&MNT_MULTILABEL != 0 {
			opts += ",multilabel"
		}
		if stat.FFlags&MNT_ACLS != 0 {
			opts += ",acls"
		}
		if stat.FFlags&MNT_NOATIME != 0 {
			opts += ",noattime"
		}
		if stat.FFlags&MNT_NOCLUSTERR != 0 {
			opts += ",nocluster"
		}
		if stat.FFlags&MNT_NOCLUSTERW != 0 {
			opts += ",noclusterw"
		}
		if stat.FFlags&MNT_NFS4ACLS != 0 {
			opts += ",nfs4acls"
		}

		d := DiskPartitionStat{
			Mountpoint: byteToString(stat.FMntonname[:]),
			Fstype:     byteToString(stat.FFstypename[:]),
			Opts:       opts,
		}
		ret = append(ret, d)
	}

	return ret, nil
}

func DiskIOCounters() (map[string]DiskIOCountersStat, error) {
	return nil, errors.New("not implemented yet")

	// statinfo->devinfo->devstat
	// /usr/include/devinfo.h

	// get length
	count, err := Getfsstat(nil, MNT_WAIT)
	if err != nil {
		return nil, err
	}

	fs := make([]Statfs, count)
	_, err = Getfsstat(fs, MNT_WAIT)

	ret := make(map[string]DiskIOCountersStat, 0)
	for _, stat := range fs {
		name := byteToString(stat.FMntonname[:])
		d := DiskIOCountersStat{
			Name:       name,
			ReadCount:  stat.FSyncwrites + stat.FAsyncwrites,
			WriteCount: stat.FSyncreads + stat.FAsyncreads,
		}

		ret[name] = d
	}

	return ret, nil
}

// Getfsstat is borrowed from pkg/syscall/syscall_freebsd.go
// change Statfs_t to Statfs in order to get more information
func Getfsstat(buf []Statfs, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	var bufsize uintptr
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
		bufsize = unsafe.Sizeof(Statfs{}) * uintptr(len(buf))
	}
	r0, _, e1 := syscall.Syscall(syscall.SYS_GETFSSTAT, uintptr(_p0), bufsize, uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}
