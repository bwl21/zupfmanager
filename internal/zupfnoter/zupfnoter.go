package zupfnoter

import (
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

func Run(ctx context.Context, args ...string) error {
	if os.Getenv("ZUPFNOTER_DEBUG") != "" {
		var files []string
		for _, arg := range args {
			fc, _ := os.ReadFile(arg)
			files = append(files, string(fc))
		}
		slog.Info("running zupfnoter", "args", args, "files", files)
		return nil
	}

	cmd := exec.CommandContext(ctx, "node", append([]string{ZupfnoterPath}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
