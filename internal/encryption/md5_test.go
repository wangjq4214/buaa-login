package encryption_test

import (
	"testing"

	"github.com/wangjq4214/buaa-login/internal/encryption"
)

func TestGetMD5(t *testing.T) {
	cases := []struct {
		name     string
		token    string
		password string
		expected string
	}{
		{"empty string", "", "", "74e6f7298a9c2d168935f58c001bad88"},
		{"simple case", "711ab370231392679fe06523b119a8fe096f5ed9bd206b4de8d7b5b994bbc3e5", "15879684798qq", "b7cc5da95734d0161fadc8ad87855e75"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if ans := encryption.GetMD5(c.token, c.password); ans != c.expected {
				t.Fatalf("GetMD5(%v, %v) should be %v, but got %v", c.token, c.password, c.expected, ans)
			}
		})
	}
}
