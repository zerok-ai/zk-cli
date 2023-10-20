package shell

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"time"
	"zkctl/cmd/internal"
	"zkctl/cmd/pkg/utils"

	"zkctl/cmd/pkg/ui"
)

const ShellToUse = "bash"

func ShelloutWithSpinner(command, spinnerText, successText, failureText string) (string, error) {

	// fmt.Println("command: " + command)
	// send the output to console as well as to a buffer
	var stdBuffer bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)

	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdBuffer

	spinner := ui.GlobalWriter.NewSpinner(spinnerText)
	if successText != "" {
		spinner.SetStopMessage(successText)
	}

	if failureText != "" {
		spinner.SetStopFailMessage(failureText)
	}

	spinner.Start()
	defer spinner.WriteStop()

	err := cmd.Run()

	if err == nil {
		return stdBuffer.String(), err
	}

	ui.GlobalWriter.PrintflnWithPrefixlnAndArrow("Got error %v", err)

	spinner.WriteStopFail()
	return stdBuffer.String(), err
}

type ExecWithSpinner func() error

func RunWithSpinner(task ExecWithSpinner, spinnerText, successText, failureText string) (string, error) {

	// fmt.Println("command: " + command)
	// send the output to console as well as to a buffer
	var stdBuffer bytes.Buffer

	spinner := ui.GlobalWriter.NewSpinner(spinnerText)
	if successText != "" {
		spinner.SetStopMessage(successText)
	}

	if failureText != "" {
		spinner.SetStopFailMessage(failureText)
	}

	spinner.Start()
	defer spinner.WriteStop()

	err := task()

	if err == nil {
		return stdBuffer.String(), err
	}

	ui.GlobalWriter.PrintflnWithPrefixlnAndArrow("Got error %v", err)

	spinner.WriteStopFail()
	return stdBuffer.String(), err
}

func Shellout(command string) (string, error) {
	printLogsOnConsole := viper.Get(internal.VerboseKeyFlag) == true
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
	out, err := Shellout(command)
	diff := time.Since(startTime)
	if err == nil {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("%s [time taken: %v]", successMsg, diff))
	}
	return out, err
}

// execute and print error only
func ExecWithLogsDurationAndSuccessM(command string, successMsg string) (string, error) {
	startTime := time.Now()
	out, err := Shellout(command)
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

func ExecuteShellFile(shellFile, inputParameters string) error {

	// make the file executable
	out, err := Shellout("chmod +x " + shellFile)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}

	out, err = Shellout(shellFile + inputParameters)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}
	return nil
}

func ExecuteShellFileWithSpinner(shellFile, inputParameters, spinnerText, successMessage, errorMessage string) error {

	// make the file executable
	out, err := ShelloutWithSpinner("chmod +x "+shellFile, "checking for appropriate permissions", "", "failed to get permission to run "+shellFile)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}

	out, err = ShelloutWithSpinner(shellFile+inputParameters, spinnerText, successMessage, errorMessage)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}
	return nil
}

func ExecuteEmbeddedFileWithSpinner(content embed.FS, filePath, inputParameters, spinnerText, successMessage, errorMessage string) error {
	installDbCode := utils.GetEmbeddedFileContents(filePath, content)

	tmpScriptName := GetPWD() + "/tmp_install_zk_file.sh"

	defer DeleteFile(tmpScriptName)

	err := utils.WriteTextToFile(installDbCode, tmpScriptName)
	if err != nil {
		return err
	}
	return ExecuteShellFileWithSpinner(tmpScriptName, inputParameters, spinnerText, successMessage, errorMessage)
}

func DeleteFile(filePath string) {
	err := utils.DeleteFile(filePath)
	if err != nil {
		//TODO: How to log the error here?
	}
}
