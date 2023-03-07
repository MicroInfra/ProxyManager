package service

import (
	"fmt"
	"os/exec"
)

func run_command() {
	fmt.Println("Hello, world.")
	cmd := exec.Command("echo", "Hello, world.")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Print the output
	fmt.Println(string(stdout))
}
