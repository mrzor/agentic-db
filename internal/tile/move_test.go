package tile

import "testing"

func TestParseMove(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		output   string
		src, dst int
		wantErr  bool
	}{
		{"plain", []string{"2", "1"}, defaultOutput, 2, 1, false},
		{"pinned output", []string{"--output", "eDP-1", "5", "3"}, "eDP-1", 5, 3, false},
		{"missing dst", []string{"2"}, "", 0, 0, true},
		{"non-numeric", []string{"a", "1"}, "", 0, 0, true},
		{"zero index", []string{"0", "1"}, "", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, src, dst, err := parseMove(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if output != tc.output || src != tc.src || dst != tc.dst {
				t.Fatalf("got (%q, %d, %d), want (%q, %d, %d)", output, src, dst, tc.output, tc.src, tc.dst)
			}
		})
	}
}
