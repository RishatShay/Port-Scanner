//go:build windows

package owner

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	modIPHlpAPI             = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTCPTable = modIPHlpAPI.NewProc("GetExtendedTcpTable")
	procGetExtendedUDPTable = modIPHlpAPI.NewProc("GetExtendedUdpTable")

	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess               = modKernel32.NewProc("OpenProcess")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle               = modKernel32.NewProc("CloseHandle")
)

const (
	afInet              = 2 // AF_INET
	tcpTableOwnerPidAll = 5 // TCP_TABLE_OWNER_PID_ALL
	udpTableOwnerPid    = 1 // UDP_TABLE_OWNER_PID

	errorInsufficientBuffer        = 122
	processQueryLimitedInformation = 0x1000
	maxPath                        = 260
)

// tcpRowOwnerPid mirrors the Win32 MIB_TCPROW_OWNER_PID struct. Ports are
// stored in network byte order inside the low 16 bits of the DWORD.
type tcpRowOwnerPid struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

// udpRowOwnerPid mirrors the Win32 MIB_UDPROW_OWNER_PID struct.
type udpRowOwnerPid struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

// lookup asks the Windows IP helper API which PID owns port, then resolves
// that PID to an executable name. Both are native Win32 calls
// (GetExtendedTcpTable/GetExtendedUdpTable, QueryFullProcessImageName) - no
// shelling out to netstat, no cgo.
func lookup(protocol string, port int) (Owner, bool) {
	var pid uint32
	var ok bool

	switch protocol {
	case "tcp":
		pid, ok = findTCPOwnerPid(port)
	case "udp":
		pid, ok = findUDPOwnerPid(port)
	default:
		return Owner{}, false
	}
	if !ok {
		return Owner{}, false
	}

	name, ok := processNameByPid(pid)
	if !ok {
		return Owner{}, false
	}

	return Owner{PID: int(pid), Name: name}, true
}

func findTCPOwnerPid(port int) (uint32, bool) {
	buf, err := fetchOwnerTable(procGetExtendedTCPTable, tcpTableOwnerPidAll)
	if err != nil {
		return 0, false
	}

	rowSize := int(unsafe.Sizeof(tcpRowOwnerPid{}))
	return scanOwnerTable(buf, rowSize, func(row unsafe.Pointer) (uint16, uint32) {
		r := (*tcpRowOwnerPid)(row)
		return uint16(r.LocalPort), r.OwningPid
	}, port)
}

func findUDPOwnerPid(port int) (uint32, bool) {
	buf, err := fetchOwnerTable(procGetExtendedUDPTable, udpTableOwnerPid)
	if err != nil {
		return 0, false
	}

	rowSize := int(unsafe.Sizeof(udpRowOwnerPid{}))
	return scanOwnerTable(buf, rowSize, func(row unsafe.Pointer) (uint16, uint32) {
		r := (*udpRowOwnerPid)(row)
		return uint16(r.LocalPort), r.OwningPid
	}, port)
}

// scanOwnerTable walks a MIB_*TABLE_OWNER_PID buffer (a leading
// dwNumEntries DWORD followed by that many fixed-size rows) looking for a
// row whose local port matches.
func scanOwnerTable(buf []byte, rowSize int, read func(unsafe.Pointer) (localPort uint16, pid uint32), port int) (uint32, bool) {
	if len(buf) < 4 {
		return 0, false
	}

	numEntries := *(*uint32)(unsafe.Pointer(&buf[0]))
	const headerSize = 4 // sizeof(dwNumEntries), rows are DWORD-aligned already

	for i := uint32(0); i < numEntries; i++ {
		offset := headerSize + int(i)*rowSize
		if offset+rowSize > len(buf) {
			break
		}

		localPort, pid := read(unsafe.Pointer(&buf[offset]))
		if int(ntohs(localPort)) == port {
			return pid, true
		}
	}

	return 0, false
}

// fetchOwnerTable calls proc twice: once with no buffer to learn the
// required size (it returns ERROR_INSUFFICIENT_BUFFER), then again with a
// buffer of that size to actually fill it in.
func fetchOwnerTable(proc *syscall.LazyProc, tableClass uintptr) ([]byte, error) {
	var size uint32
	var buf []byte

	for retry := 0; retry < 5; retry++ {
		var bufPtr uintptr
		if len(buf) > 0 {
			bufPtr = uintptr(unsafe.Pointer(&buf[0]))
		}

		ret, _, _ := proc.Call(
			bufPtr,
			uintptr(unsafe.Pointer(&size)),
			0, // bOrder: don't need results sorted
			uintptr(afInet),
			tableClass,
			0,
		)

		switch ret {
		case 0:
			return buf, nil
		case errorInsufficientBuffer:
			buf = make([]byte, size)
		default:
			return nil, fmt.Errorf("GetExtended*Table failed with code %d", ret)
		}
	}

	return nil, fmt.Errorf("GetExtended*Table: buffer never large enough")
}

func ntohs(v uint16) uint16 {
	return (v >> 8) | (v << 8)
}

func processNameByPid(pid uint32) (string, bool) {
	handle, _, _ := procOpenProcess.Call(
		uintptr(processQueryLimitedInformation),
		0,
		uintptr(pid),
	)
	if handle == 0 {
		return "", false
	}
	defer procCloseHandle.Call(handle)

	var buf [maxPath]uint16
	size := uint32(len(buf))
	ret, _, _ := procQueryFullProcessImageName.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", false
	}

	path := syscall.UTF16ToString(buf[:size])
	if i := strings.LastIndexAny(path, `\/`); i >= 0 {
		path = path[i+1:]
	}
	return path, true
}
