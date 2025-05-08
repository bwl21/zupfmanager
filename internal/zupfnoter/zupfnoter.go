package zupfnoter

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	_ "embed"
)

var (
	//go:embed zupfnoter-cli.min.js
	zupfnoter []byte

	ZupfnoterPath string
)

func init() {
	ZupfnoterPath = os.Getenv("ZUPFNOTER_PATH")
	if ZupfnoterPath == "" {
		tempFile, err := os.CreateTemp("", "zupfnoter-*.js")
		if err != nil {
			panic(fmt.Errorf("failed to create temp file: %w", err))
		}
		tempFile.Write(zupfnoter)
		tempFile.Close()
		ZupfnoterPath = tempFile.Name()
	}
}

func Run(ctx context.Context, args ...string) (string, string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	if os.Getenv("ZUPFNOTER_DEBUG") != "" {
		var files []string
		for _, arg := range args {
			fc, _ := os.ReadFile(arg)
			files = append(files, string(fc))
		}
		slog.Info("running zupfnoter", "args", args, "files", files)
		return "", "", nil
	}

	cmd := exec.CommandContext(ctx, "node", append([]string{ZupfnoterPath}, args...)...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}
