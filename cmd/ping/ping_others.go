//go:build !windows

package ping

import (
	"log"
	"os/exec"
)

func Ping(ip string) error {
	args := []string{ip, "-c", "4"}

	path, err := exec.LookPath("ping")
	if err != nil {
		log.Fatalf("We got an error while finding the ping command, err: %v\n", err.Error())
	}

	cmd := exec.Command(path, args...)
	return cmd.Run()
}
