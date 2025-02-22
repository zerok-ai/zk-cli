package ui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/viper"

	"cloud.google.com/go/errorreporting"
)

const (
	ASSUME_YES_FLAG = "yes"
)

type Writer struct {
	writen []string
}

var (
	greenStatusOk    = color.GreenString(statusOK)
	redStatusErr     = color.RedString(statusErr)
	writenStatusOk   = "V"
	writenStatusErr  = "X"
	writenStatusWarn = "!"

	ErrorReporter *errorreporting.Client
)

func NewWriter() *Writer {
	return &Writer{
		writen: []string{},
	}
}

var GlobalWriter = NewWriter()
var QuietWriter = NewWriter()

func (w *Writer) MarshalJSON() ([]byte, error) {
	return json.Marshal((w.Dump()))
}

func (w *Writer) Writeln(message string) {
	w.addMessage(fmt.Sprintln(message))
}

func (w *Writer) Println(message string) {
	w.addMessage(message)
	fmt.Println(message)
}

func (w *Writer) PrintlnWithPrefixln(message string) {
	w.addMessage(message)
	fmt.Printf("\n%s\n", message)
}

func (w *Writer) PrintflnWithPrefixln(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	w.addMessage(message)
	fmt.Printf("\n%s\n", message)
}

func (w *Writer) PrintflnWithPrefixArrow(format string, args ...interface{}) {
	message := fmt.Sprintf("→ "+format, args...)
	w.addMessage(message)
	fmt.Printf("%s\n", message)
}

func (w *Writer) PrintflnWithPrefixlnAndArrow(format string, args ...interface{}) {
	message := fmt.Sprintf("→ "+format, args...)
	w.addMessage(message)
	fmt.Printf("\n%s\n", message)
}

func (w *Writer) Printf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	w.addMessage(formatted)
	fmt.Print(formatted)
}

func (w *Writer) PrintUrl(message string, url string) {
	w.addMessage(fmt.Sprintf("%s%s", message, url))
	fmt.Printf("%s%s\n", message, w.UrlLink(url))
}

func (w *Writer) Errorf(format string, args ...interface{}) error {
	formatted := fmt.Sprintf(format, args...)
	w.addMessage(formatted)
	return errors.New(formatted)
}

func (w *Writer) PrintSuccessMessage(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusOk, message))
	fmt.Printf("%s %s", greenStatusOk, message)
}

func (w *Writer) PrintSuccessMessageln(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusOk, message))
	fmt.Printf("%s %s\n", greenStatusOk, message)
}

func (w *Writer) PrintlnSuccessMessageln(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusOk, message))
	fmt.Printf("\n%s %s\n", greenStatusOk, message)
}

func (w *Writer) PrintErrorMessage(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusErr, message))
	fmt.Printf("%s %s", redStatusErr, message)
}

func (w *Writer) PrintErrorMessageln(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusErr, message))
	fmt.Printf("%s %s\n", redStatusErr, message)
}

func (w *Writer) PrintWarningMessage(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusWarn, message))
	fmt.Printf("%s %s", statusWarning, message)
}

func (w *Writer) PrintWarningMessageln(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusWarn, message))
	fmt.Printf("%s %s\n", statusWarning, message)
}

func (w *Writer) PrintlnWarningMessageln(message string) {
	w.addMessage(fmt.Sprintf("%s %s", writenStatusWarn, message))
	fmt.Printf("\n%s %s\n", statusWarning, message)
}

func (w *Writer) PrintNoticeMessage(message string) {
	w.addMessage(message)
	fmt.Printf("🚨 %s", message)
}

func (w *Writer) UrlLink(url string) string {
	return color.New(color.FgBlue).Add(color.Underline).Sprint(url)
}

func (w *Writer) NewSpinner(message string) *Spinner {
	w.addMessage(message)
	return newSpinner(w, message)
}

func (w *Writer) YesNoPrompt(message string, defaultValue bool) bool {
	if viper.GetBool(ASSUME_YES_FLAG) {
		return true
	}

	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}

	var answer bool
	survey.AskOne(prompt, &answer)
	w.addMessage(fmt.Sprintf("%s %t", message, answer))
	return answer
}

func (w *Writer) QuestionPrompt(message string) string {

	var qs = []*survey.Question{
		{
			Name:   "response",
			Prompt: &survey.Input{Message: message},
		},
	}
	answer := struct {
		Response string // the name of the valiable is mapped with the name of the question
	}{}

	// ask the questions
	err := survey.Ask(qs, &answer)
	if err != nil {
		return ""
	}

	w.addMessage(fmt.Sprintf("%s %s", message, answer.Response))
	return answer.Response
}

func (w *Writer) MultiSelectPrompt(message string, options, defaults []string) []string {
	if viper.GetBool(ASSUME_YES_FLAG) {
		return defaults
	}

	prompt := &survey.MultiSelect{
		Options: options,
		Default: defaults,
		Message: message,
	}

	var response []string
	survey.AskOne(prompt, &response)
	w.addMessage(fmt.Sprintf("%s %v", message, response))
	return response
}

func (w *Writer) timeFormat(message string) string {
	timeFormatted := time.Now().Format(time.RFC3339)

	return fmt.Sprintf("%s - %s", timeFormatted, message)
}

func (w *Writer) addMessage(message string) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		w.writen = append(w.writen, w.timeFormat(line))
	}
}

func (w *Writer) Dump() string {
	return strings.Join(w.writen, "\n")
}

func LogAndPrintError(err error) {
	ErrorReporter.Report(errorreporting.Entry{
		Error: err,
	})
	log.Print(err)
}

// initialize the error reporter and return a `cleanup` function.
// cleanup function should be ideally called at the end of the program
func InitializeErrorReportor(ctx context.Context) func() {

	// Sets your Google Cloud Platform project ID.
	projectID := "zerok-dev"

	var err error
	ErrorReporter, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName:    "zkctl",
		ServiceVersion: "v1.0",
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	var cleanup func()
	if err != nil {
		log.Fatal(err)
	} else {
		cleanup = func() {
			if err := ErrorReporter.Close(); err != nil {
				log.Printf("failed to report errors to Cloud Error Reporting: %v", err)
			}
		}
	}

	return cleanup
}
