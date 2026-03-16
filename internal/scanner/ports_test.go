package scanner

import (
	"reflect"
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:  "single port",
			input: "80",
			want:  []int{80},
		},
		{
			name:  "list of ports",
			input: "22,80,443",
			want:  []int{22, 80, 443},
		},
		{
			name:  "range",
			input: "1-5",
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "mixed list and range",
			input: "22,100-102",
			want:  []int{22, 100, 101, 102},
		},
		{
			name:    "invalid port",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid range",
			input:   "10-5",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePorts(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
