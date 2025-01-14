package utils

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"zkctl/cmd/pkg/ui"

	"github.com/spf13/viper"
)

const (
	ZEROK_DIR_PATH_FLAG       = "zkdir"
	BACKEND_CLI_NAME          = "daemon"
	BACKEND_CLI_AUTH_FILENAME = "auth.json"
	ERROR_DUMP_FILENAME       = "dump"

	px_dir_sym_name = ".pixie"
	px_dir_name     = "px"
	pxrepo_dir      = "pxrepo"
)

func WriteTextToFile(text, filePath string) error {
	// Open the file for appending (or create if it doesn't exist)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(text)
	return err
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string, filedDir string, filename string) error {

	// create the directory
	err := os.MkdirAll(filedDir, os.ModePerm)
	if err != nil {
		ui.GlobalWriter.PrintNoticeMessage(fmt.Sprint(err))
	}

	// Get the data
	filepath := fmt.Sprintf("%s%c%s", filedDir, os.PathSeparator, filename)
	resp, err := http.Get(url)
	if err != nil {
		return ui.RetryableError(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadFileAndShowSpinner(url string, filedDir string, filename string) error {
	var err error

	spinner := ui.GlobalWriter.NewSpinner("Downloading necessary files")
	spinner.SetStopMessage("download successful")
	spinner.SetStopFailMessage("download failed")

	spinner.Start()
	defer spinner.WriteStop()

	err = DownloadFile(url, filedDir, filename)
	if err == nil {
		return nil
	}
	spinner.WriteStopFail()

	if errors.Is(err, ui.ErrSpinnerTimeout) {
		return errors.New("timeout waiting for download")
	}

	return err
}

// This internally calls DownloadFile and then gives it executable permission
func DownloadExecutableFile(url string, filedDir string, filename string, showSpinner bool) error {

	var err error

	// download the file
	filepath := fmt.Sprintf("%s%c%s", filedDir, os.PathSeparator, filename)
	if showSpinner {
		err = DownloadFileAndShowSpinner(url, filedDir, filename)
	} else {
		err = DownloadFile(url, filedDir, filename)
	}
	if err != nil {
		return err
	}

	// add execute permission to the file for current user
	stats, err := os.Stat(filepath)
	if err != nil {
		return ui.RetryableError(err)
	}
	err = os.Chmod(filepath, stats.Mode()|0100)
	if err != nil {
		return ui.RetryableError(err)
	}

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateDirAndSymLinkIfNotExists(oldname, newname string) error {

	exist := Exists(oldname)
	if !exist {
		err := os.MkdirAll(oldname, os.ModePerm)
		if err != nil {
			ui.GlobalWriter.PrintNoticeMessage(fmt.Sprint(err))
			return fmt.Errorf("unable to create dir %s", oldname)
		}
	}

	exist = Exists(newname)
	if exist {
		// // check if it is a symlink
		// fileInfo, _ := os.Lstat(newname)
		// fmt.Printf("%s mode = %v\n", fileInfo.Name(), fileInfo.Mode())

		// if fileInfo.Mode()&os.ModeSymlink != 0 {
		// 	fmt.Printf("%s is a symlink\n", fileInfo.Name())
		// }
		return nil
	}

	err := os.Symlink(oldname, newname)
	return err
}

func InitializeFolders() error {
	// create simlink
	err := CreateDirAndSymLinkIfNotExists(getPxDirPath(), getPxDirSymbolicPath())
	if err != nil {
		panic(err)
	}
	return nil
}

func ResetErrorDumpfile() {
	prettyTime := time.Now().Format("Monday, January 2, 2006 03:04:05 PM MST")
	initLogString := "\n---------------------------------------------------------------\n"
	initLogString = fmt.Sprintf("%sCli dump file appended at %s\n", initLogString, prettyTime)
	dumpPath := GetErrorDumpPath()
	err := WriteTextToFile(initLogString, dumpPath)
	if err != nil {
		return
	}
}

func DumpErrorAndPrintLocation(errorMessage string, printOnConsole bool) {
	dumpPath := GetErrorDumpPath()
	err := WriteTextToFile(errorMessage, dumpPath)
	if err == nil && printOnConsole {
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("Error:Check %s for details\n", dumpPath))
	}
	if printOnConsole {
		ui.LogAndPrintError(fmt.Errorf(errorMessage+": %v\n", errorMessage))
	}
}

func GetPxRepoDir(zkdir string) string {
	return fmt.Sprintf("%s%c%s", zkdir, os.PathSeparator, pxrepo_dir)
}

func getPxDirPath() string {
	zkPath := viper.Get(ZEROK_DIR_PATH_FLAG)
	return fmt.Sprintf("%s%c%s", zkPath, os.PathSeparator, px_dir_name)
}

func getPxDirSymbolicPath() string {
	var err error

	var baseDir string
	if baseDir, err = os.UserHomeDir(); err != nil {
		baseDir = os.TempDir()
	}

	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, px_dir_sym_name)
}

func getBackendCLIDir() string {
	return viper.GetString(ZEROK_DIR_PATH_FLAG)
}

func GetBackendCLIPath() string {
	return fmt.Sprintf("%s%c%s", getBackendCLIDir(), os.PathSeparator, BACKEND_CLI_NAME)
}

func GetBackendAuthPath() string {
	return fmt.Sprintf("%s%c%s", getPxDirSymbolicPath(), os.PathSeparator, BACKEND_CLI_AUTH_FILENAME)
}

func GetErrorDumpPath() string {
	return fmt.Sprintf("%s%c%s", getBackendCLIDir(), os.PathSeparator, ERROR_DUMP_FILENAME)
}

func BackendCLIExists() error {
	_, err := os.Stat(GetBackendCLIPath())
	return err
}

func DownloadBackendCLI(url string) error {
	return DownloadExecutableFile(url, getBackendCLIDir(), BACKEND_CLI_NAME, true)
}

func GetEmbeddedFileContents(filename string, content embed.FS) string {
	data, err := content.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading embedded file:", err)
		return ""
	}
	contentString := string(data)
	return contentString
}
