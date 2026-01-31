package update

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		// Basic cases
		{"same version", "1.0.0", "1.0.0", false},
		{"newer major", "2.0.0", "1.0.0", true},
		{"newer minor", "1.1.0", "1.0.0", true},
		{"newer patch", "1.0.1", "1.0.0", true},
		{"older major", "1.0.0", "2.0.0", false},
		{"older minor", "1.0.0", "1.1.0", false},
		{"older patch", "1.0.0", "1.0.1", false},

		// With v prefix
		{"v prefix latest", "v1.1.0", "1.0.0", true},
		{"v prefix current", "1.1.0", "v1.0.0", true},
		{"v prefix both", "v1.1.0", "v1.0.0", true},

		// Different lengths
		{"longer latest", "1.0.1", "1.0", true},
		{"longer current", "1.0", "1.0.1", false},
		{"much newer", "2.5.3", "1.2.1", true},

		// Edge cases
		{"double digits", "1.10.0", "1.9.0", true},
		{"triple digits", "1.0.100", "1.0.99", true},

		// Non-numeric parts (fall back to string comparison)
		{"non-numeric part", "1.x.0", "1.w.0", true},  // "x" > "w" lexically
		{"mixed numeric non-numeric", "2.0.0", "1.beta.0", true}, // 2 > 1
		{"all non-numeric", "b.0.0", "a.0.0", true}, // "b" > "a" lexically
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNewerVersion(tt.latest, tt.current)
			if result != tt.expected {
				t.Errorf("IsNewerVersion(%q, %q) = %v, want %v",
					tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}
