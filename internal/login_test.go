package internal_test

import (
	"testing"

	"github.com/wangjq4214/buaa-login/internal"
)

func TestLogin(t *testing.T) {
	lm := internal.NewLoginManager(internal.NewLoginManagerParams{
		Username:           "by2106105",
		Password:           "",
		N:                  "200",
		AcID:               "1",
		Enc:                "srun_bx1",
		VType:              "1",
		LoginPageURL:       "https://gw.buaa.edu.cn/srun_portal_pc?ac_id=1&theme=buaa",
		GetChallengeApiURL: "http://gw.buaa.edu.cn/cgi-bin/get_challenge",
		LoginApiURL:        "http://gw.buaa.edu.cn/cgi-bin/srun_portal",
	})

	err := lm.Login()
	if err != nil {
		t.Fatalf("We got an error that is %v", err.Error())
	}
}
