package encryption_test

import (
	"testing"

	"github.com/wangjq4214/buaa-login/encryption"
)

func TestGetSHA1(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		expected string
	}{
		{"empty string", "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"simple case 1", "123456", "7c4a8d09ca3762af61e59520943dc26494f8941b"},
		{"simple case 2", "mldyyds", "7e27bffdd616bbd5fe6bf52c14bbd60ad2c1bd0b"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if ans := encryption.GetSHA1(c.s); ans != c.expected {
				t.Fatalf("GetSHA1(%v) should be %v, but got %v", c.s, c.expected, ans)
			}
		})
	}
}
