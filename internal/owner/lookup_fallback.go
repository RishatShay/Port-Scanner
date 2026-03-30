//go:build !linux && !windows

package owner

// lookup isn't implemented for this OS. Learning which process holds a
// socket here would require either cgo (e.g. macOS's proc_pidinfo) or
// shelling out to a tool like lsof/netstat, which is exactly the
// OS-dependent, shell-based approach this package is meant to avoid.
// Callers fall back to reporting the port as open without an owner.
func lookup(protocol string, port int) (Owner, bool) {
	return Owner{}, false
}
