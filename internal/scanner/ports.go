package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePorts turns a string like "22,80,443,8000-8100" into a list of port
// numbers. Ranges are inclusive on both ends.
func ParsePorts(input string) ([]int, error) {
	var ports []int

	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if start, end, ok := strings.Cut(part, "-"); ok {
			from, err := strconv.Atoi(strings.TrimSpace(start))
			if err != nil {
				return nil, fmt.Errorf("invalid range %q: %w", part, err)
			}
			to, err := strconv.Atoi(strings.TrimSpace(end))
			if err != nil {
				return nil, fmt.Errorf("invalid range %q: %w", part, err)
			}
			if from > to {
				return nil, fmt.Errorf("invalid range %q: start is greater than end", part)
			}
			for p := from; p <= to; p++ {
				ports = append(ports, p)
			}
			continue
		}

		port, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", part, err)
		}
		ports = append(ports, port)
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports found in %q", input)
	}

	return ports, nil
}
