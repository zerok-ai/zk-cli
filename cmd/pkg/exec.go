package pkg

import (
	"fmt"
	"io"
	"os"

	"bytes"
	"os/exec"

	"time"

	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const ShellToUse = "bash"

func shellout(command string) error {

	// send the output to console as well as to a buffer
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = mw
	cmd.Stderr = mw

	return cmd.Run()
}

// execute and print error only
func execOnShell(command string) (time.Duration, error) {
	startTime := time.Now()
	err := shellout(command)
	diff := time.Since(startTime)
	return diff, err
}

func ExecOnShellM(command string, successMsg string) error {
	diff, err := execOnShell(command)
	if err == nil {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("%s [time taken: %v]\n", successMsg, diff))
	}
	return err
}

func GetPWD() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return mydir
}
