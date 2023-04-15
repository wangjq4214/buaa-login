package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	login "github.com/wangjq4214/buaa-login"
	"github.com/wangjq4214/buaa-login/cmd/ping"
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

func getLogPath(pwd, daemonLog string, idx uint32) string {
	return filepath.Join(pwd, fmt.Sprintf("%v-%v.log", daemonLog, idx))
}

func setLog(pwd string, idx *uint32) *os.File {
	stdout, err := os.OpenFile(getLogPath(pwd, daemonLog, *idx), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	*idx++

	if err != nil {
		log.Printf("Can not open file, err: %v\n", err.Error())
	}

	old_out := os.Stdout

	os.Stdout = stdout
	os.Stderr = stdout

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(stdout)

	return old_out
}

func logFile() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Can not get current dir, err: %v\n", err.Error())
	}

	var idx uint32
	setLog(pwd, &idx)

	go func() {
		stat, _ := os.Stat(getLogPath(pwd, daemonLog, idx-1))
		if stat.Size() > 5*1024*1024 {
			log.Println("remove")
			old_file := setLog(pwd, &idx)
			old_file.Close()

			os.Remove(old_file.Name())
		}

		time.Sleep(30 * 60 * time.Second)
	}()
}

func childHandler() {
	timer := time.NewTimer(0)

	logFile()

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

		func() {
			defer timer.Reset(5 * time.Second)

			err := ping.Ping(daemonIP)
			if err == nil {
				log.Println("The PC is online.")
				return
			} else {
				log.Printf("We got an error while ping, err: %v\n", err.Error())
			}

			for i := 0; i < 5; i++ {
				err := lm.Login()
				if err != nil {
					log.Printf("We got an error while login, err: %v\n", err.Error())
					continue
				}

				return
			}
		}()
	}
}

func init() {
	daemonCmd.Flags().StringVarP(&daemonUsername, "username", "u", "", "Your buaa gw username.")
	daemonCmd.MarkFlagRequired("username")

	daemonCmd.Flags().StringVarP(&daemonPassword, "password", "p", "", "Your buaa gw password.")
	daemonCmd.MarkFlagRequired("password")

	daemonCmd.Flags().StringVar(&daemonIP, "ip", "baidu.com", "Test ip with ping.")
	daemonCmd.Flags().StringVar(&daemonLog, "log", "login", "Specify the log file path.")

	rootCmd.AddCommand(daemonCmd)
}
