package pkg

import (
	"fmt"
	"log"
	"os"

	"bytes"
	"os/exec"

	"time"
)

const ShellToUse = "bash"

func shellout(command string, printSuccessLogs bool, printErrorLogs bool) (string, error) {
	/**/
	startTime := time.Now()

	var shellLog bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	if printSuccessLogs {
		cmd.Stdout = &shellLog
	}
	if printErrorLogs {
		cmd.Stderr = &shellLog
	}
	err := cmd.Run()

	diff := time.Since(startTime)

	fmt.Printf("[time taken: %v]\n", diff)

	return shellLog.String(), err
	/*/
	var err error
	return command, err
	/**/
}

// execute and print error only
func ExecOnShellsE(command string) {
	ExecOnShell(command, true, true)
}

// execute and print
func ExecOnShell(command string, printSuccessLogs bool, printErrorLogs bool) {
	// log.Printf("command: %v\n", command)
	out, err := shellout(command, printSuccessLogs, printErrorLogs)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println(out)
}

func GetPWD() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return mydir
}
