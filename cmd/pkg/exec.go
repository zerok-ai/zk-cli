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

func Shellout(command string) (string, error) {

	// send the output to console as well as to a buffer
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Run()
	return stdBuffer.String(), err
}

// execute and print error only
func execOnShell(command string) (string, time.Duration, error) {
	startTime := time.Now()
	out, err := Shellout(command)
	diff := time.Since(startTime)
	return out, diff, err
}

func ExecOnShellM(command string, successMsg string) (string, error) {
	out, diff, err := execOnShell(command)
	if err == nil {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("%s [time taken: %v]\n", successMsg, diff))
	}
	return out, err
}

func GetPWD() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return mydir
}
