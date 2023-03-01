package shell

import (
	"fmt"
	"io"
	"os"

	"bytes"
	"os/exec"

	"time"

	"zkctl/cmd/pkg/ui"
)

const ShellToUse = "bash"

func Shellout(command string, printLogsOnConsole bool) (string, error) {

	fmt.Println("command: "+command)
	// send the output to console as well as to a buffer
	var stdBuffer bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)

	var mw io.Writer
	if printLogsOnConsole {
		mw = io.MultiWriter(os.Stdout, &stdBuffer)
	} else {
		mw = io.MultiWriter(&stdBuffer)
	}
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Run()
	return stdBuffer.String(), err
}

// execute and print error only
func ExecWithDurationAndSuccessM(command string, successMsg string) (string, error) {
	startTime := time.Now()
	out, err := Shellout(command, false)
	diff := time.Since(startTime)
	if err == nil {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("%s [time taken: %v]\n", successMsg, diff))
	}
	return out, err
}

// execute and print error only
func ExecWithLogsDurationAndSuccessM(command string, successMsg string) (string, error) {
	startTime := time.Now()
	out, err := Shellout(command, true)
	diff := time.Since(startTime)
	if err == nil {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("%s [time taken: %v]\n", successMsg, diff))
	} else {
		ui.GlobalWriter.Printf("\n ExecWithLogsDurationAndSuccessM err --- %v\n", err)
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
