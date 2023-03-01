package utils

import (
	"zkctl/cmd/pkg/ui"

	"github.com/skratchdot/open-golang/open"
)

func TryOpenBrowser(writer *ui.Writer, message string, url string) {
	writer.PrintUrl(message, url)
	open.Run(url)
}
