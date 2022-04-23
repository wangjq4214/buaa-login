package encryption_test

import (
	"testing"

	"github.com/wangjq4214/buaa-login/internal/encryption"
)

func TestGetXencode(t *testing.T) {
	cases := []struct {
		name     string
		msg      string
		key      string
		expected string
	}{
		{"empty string", "{\"username\":\"201626203044@cmcc\",\"password\":\"15879684798qq\",\"ip\":\"10.128.96.249\",\"acid\":\"1\",\"enc_ver\":\"srun_bx1\"}", "e6843f26b8544327a3a25978dd3c5f89e6b745df1732993b88fe082c13a34cb9", "13GwOQhjyto7UD3YETXHKszNW6cgCyaeZxbRoFRKgNRJTnTSqC/awYNrdZP1cgJfTPvesb2/jkTwKUtyOvK1yZkmA25ShYWGhKahj/1p0QO3aW/8Ue8NRUy0QcMqtvyS3XdsFypSV9EO10kTcp1PXHvhH64="},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if ans := encryption.GetBase64(encryption.GetXencode(c.msg, c.key)); ans != c.expected {
				t.Fatalf("GetXencode(%v, %v) should be %v, but got %v", c.msg, c.key, c.expected, ans)
			}
		})
	}
}
