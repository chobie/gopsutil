// +build darwin

package gopsutil

import (
	"os/exec"
	"strings"
)

func NetIOCounters(pernic bool) ([]NetIOCountersStat, error) {
	out, err := exec.Command("/usr/sbin/netstat", "-ibdn").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	ret := make([]NetIOCountersStat, 0, len(lines)-1)

	for _, line := range lines {
		values := strings.Fields(line)
		if len(values) < 1 || values[0] == "Name" {
			continue
		}
		base := 1
		// sometimes Address is ommitted
		if len(values) < 12 {
			base = 0
		}

		n := NetIOCountersStat{
			Name:        values[0],
			PacketsRecv: mustParseUint64(values[base+3]),
			Errin:       mustParseUint64(values[base+4]),
			//Dropin:      mustParseUint64(values[base+5]),
			BytesRecv:   mustParseUint64(values[base+5]),
			PacketsSent: mustParseUint64(values[base+6]),
			Errout:      mustParseUint64(values[base+7]),
			BytesSent:   mustParseUint64(values[base+8]),
			Dropout:     mustParseUint64(values[base+9]),
		}
		ret = append(ret, n)
	}

	return ret, nil
}
