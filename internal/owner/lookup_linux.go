//go:build linux

package owner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// lookup finds the process bound to port by walking /proc, the same data
// source tools like ss/lsof read on Linux. No shelling out is needed:
// /proc/net/<proto> maps a listening port to a socket inode, and
// /proc/<pid>/fd links inodes to the process that owns them.
func lookup(protocol string, port int) (Owner, bool) {
	inode, ok := findInode(protocol, port)
	if !ok {
		return Owner{}, false
	}

	pid, ok := findPidByInode(inode)
	if !ok {
		return Owner{}, false
	}

	name, err := processName(pid)
	if err != nil {
		return Owner{}, false
	}

	return Owner{PID: pid, Name: name}, true
}

// findInode scans /proc/net/<protocol> and its IPv6 counterpart for a
// socket bound to port, and returns its inode as a string.
func findInode(protocol string, port int) (string, bool) {
	wantPort := fmt.Sprintf("%04X", port)

	for _, path := range []string{"/proc/net/" + protocol, "/proc/net/" + protocol + "6"} {
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		inode, found := scanProcNet(f, protocol, wantPort)
		f.Close()
		if found {
			return inode, true
		}
	}

	return "", false
}

func scanProcNet(f *os.File, protocol, wantPort string) (string, bool) {
	sc := bufio.NewScanner(f)
	sc.Scan() // header line

	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 10 {
			continue
		}

		_, hexPort, ok := strings.Cut(fields[1], ":") // local_address = "IP:PORT" in hex
		if !ok || !strings.EqualFold(hexPort, wantPort) {
			continue
		}

		// TCP sockets are only "open" for scanning purposes while
		// listening (state 0A); UDP has no such state, any bound
		// socket counts.
		if protocol == "tcp" && fields[3] != "0A" {
			continue
		}

		return fields[9], true // inode
	}

	return "", false
}

// findPidByInode walks /proc/<pid>/fd looking for a symlink to
// socket:[inode]. Processes owned by other users are silently skipped
// (permission denied), same as they would be for an unprivileged lsof.
func findPidByInode(inode string) (int, bool) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, false
	}

	target := "socket:[" + inode + "]"

	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := "/proc/" + entry.Name() + "/fd"
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(fdDir + "/" + fd.Name())
			if err == nil && link == target {
				return pid, true
			}
		}
	}

	return 0, false
}

// processName reads the short process name from /proc/<pid>/comm.
func processName(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
