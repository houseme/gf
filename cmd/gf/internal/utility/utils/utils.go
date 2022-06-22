package utils

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	gofmtPath     = gproc.SearchBinaryPath("gofmt")     // gofmtPath is the binary path of command `gofmt`.
	goimportsPath = gproc.SearchBinaryPath("goimports") // gofmtPath is the binary path of command `goimports`.
)

func init() {
	// Wraps the command binary path with char '"' if there's space char in the path.
	if gstr.Contains(gofmtPath, " ") {
		gofmtPath = fmt.Sprintf(`"%s"`, gofmtPath)
	}
	if gstr.Contains(goimportsPath, " ") {
		goimportsPath = fmt.Sprintf(`"%s"`, goimportsPath)
	}
}

// GoFmt formats the source file using command `gofmt -w -s PATH`.
func GoFmt(path string) {
	if gofmtPath == "" {
		mlog.Fatal(`command "gofmt" not found`)
	}
	var command = fmt.Sprintf(`%s -w %s`, gofmtPath, path)
	result, err := gproc.ShellExec(context.Background(), command)
	if err != nil {
		mlog.Fatalf(`error executing command "%s": %s`, command, result)
	}
}

// GoImports formats the source file using command `goimports -w PATH`.
func GoImports(path string) {
	if goimportsPath == "" {
		mlog.Fatal(`command "goimports" not found`)
	}
	var command = fmt.Sprintf(`%s -w %s`, goimportsPath, path)
	result, err := gproc.ShellExec(context.Background(), command)
	if err != nil {
		mlog.Fatalf(`error executing command "%s": %s`, command, result)
	}
}

// IsFileDoNotEdit checks and returns whether file contains `do not edit` key.
func IsFileDoNotEdit(filePath string) bool {
	if !gfile.Exists(filePath) {
		return true
	}
	return gstr.Contains(gfile.GetContents(filePath), consts.DoNotEditKey)
}
