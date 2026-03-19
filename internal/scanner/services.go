package scanner

// wellKnownPorts maps common port numbers to the service that usually runs
// on them. It only covers the ports people run into most often.
var wellKnownPorts = map[int]string{
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	80:    "http",
	110:   "pop3",
	143:   "imap",
	443:   "https",
	445:   "smb",
	3306:  "mysql",
	3389:  "rdp",
	5432:  "postgresql",
	6379:  "redis",
	8080:  "http-alt",
	8443:  "https-alt",
	27017: "mongodb",
}

// ServiceName returns the commonly known service name for a port,
// or "unknown" if the port isn't in the list.
func ServiceName(port int) string {
	if name, ok := wellKnownPorts[port]; ok {
		return name
	}
	return "unknown"
}
