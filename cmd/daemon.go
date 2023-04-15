package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	login "github.com/wangjq4214/buaa-login"
)

const envName = "BUAA_LOGIN_DAEMON"
const envValue = "CHILD"

var (
	daemonUsername string
	daemonPassword string
	daemonIP       string
	daemonLog      string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start a daemon to login the campus network.",
	Run: func(cmd *cobra.Command, args []string) {
		val := os.Getenv(envName)
		if val != envValue {
			fatherHandler()
			return
		}

		childHandler()
	},
}

func fatherHandler() {
	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,
		Env:  os.Environ(),
	}

	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))
	cmd.Start()
}

func childHandler() {
	timer := time.NewTimer(0)
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Can not get current dir, err: %v\n", err.Error())
	}
	stdout, err := os.OpenFile(filepath.Join(pwd, daemonLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Can not open file, err: %v\n", err.Error())
	}

	os.Stdout = stdout
	os.Stderr = stdout

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(stdout)

	lm := login.NewLoginManager(login.NewLoginManagerParams{
		Username:           daemonUsername,
		Password:           daemonPassword,
		N:                  "200",
		AcID:               "1",
		Enc:                "srun_bx1",
		VType:              "1",
		LoginPageURL:       "https://gw.buaa.edu.cn/srun_portal_pc?ac_id=1&theme=buaa",
		GetChallengeApiURL: "https://gw.buaa.edu.cn/cgi-bin/get_challenge",
		LoginApiURL:        "https://gw.buaa.edu.cn/cgi-bin/srun_portal",
	})

	for {
		<-timer.C

		sysType := runtime.GOOS
		args := []string{daemonIP}
		if sysType != "windows" {
			args = append(args, "-c", "4")
		}

		path, err := exec.LookPath("ping")
		if err != nil {
			log.Fatalf("We got an error while finding the ping command, err: %v\n", err.Error())
		}

		cmd := exec.Command(path, args...)
		err = cmd.Run()
		if err == nil {
			log.Println("The PC is online.")
			timer.Reset(5 * time.Second)
			continue
		} else {
			log.Printf("We got an error while ping, err: %v\n", err.Error())
		}

		for i := 0; i < 5; i++ {
			err := lm.Login()
			if err != nil {
				log.Printf("We got an error while login, err: %v\n", err.Error())
			} else {
				break
			}
		}

		timer.Reset(5 * time.Second)
	}
}

func init() {
	daemonCmd.Flags().StringVarP(&daemonUsername, "username", "u", "", "Your buaa gw username.")
	daemonCmd.MarkFlagRequired("username")

	daemonCmd.Flags().StringVarP(&daemonPassword, "password", "p", "", "Your buaa gw password.")
	daemonCmd.MarkFlagRequired("password")

	daemonCmd.Flags().StringVar(&daemonIP, "ip", "baidu.com", "Test ip with ping.")
	daemonCmd.Flags().StringVar(&daemonLog, "log", "login.log", "Specify the log file path.")

	rootCmd.AddCommand(daemonCmd)
}
