package owner

import "testing"

func TestIsLocalHost(t *testing.T) {
	cases := map[string]bool{
		"localhost":       true,
		"127.0.0.1":       true,
		"::1":             true,
		"scanme.nmap.org": false,
		"192.0.2.1":       false, // TEST-NET-1, never loopback
	}

	for host, want := range cases {
		if got := IsLocalHost(host); got != want {
			t.Errorf("IsLocalHost(%q) = %v, want %v", host, got, want)
		}
	}
}
