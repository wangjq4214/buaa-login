package encryption_test

import (
	"testing"

	"github.com/wangjq4214/buaa-login/encryption"
)

func TestGetBase64(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		expected string
	}{
		{"empty string", "", ""},
		{"simple case", "132456", "9F9x0JHI"},
		{"two = test case", "mldyyds", "Wa6YrcpYMv=="},
		{"one = test case", "mldyydss", "Wa6YrcpYM79="},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if ans := encryption.GetBase64(c.s); ans != c.expected {
				t.Fatalf("GetBase64(%v) should be %v, but got %v", c.s, c.expected, ans)
			}
		})
	}
}
