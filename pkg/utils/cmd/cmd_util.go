package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SetLogger 设置命令日志，并返回 logr.Logger
func SetLogger(cmd *cobra.Command, verbosity uint32) logr.Logger {
	// 设置日志级别
	logrusLogger := logrus.New()
	switch verbosity {
	case 1:
		logrusLogger.SetLevel(logrus.DebugLevel)
	case 2:
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}
	// 将 logger 注入上下文
	logger := logrusr.New(logrusLogger)
	cmd.SetContext(logr.NewContext(cmd.Context(), logger))

	return logger
}

// ChangeWorkingDirectory 切换命令工作目录
func ChangeWorkingDirectory(cmd *cobra.Command, path string) error {
	defer func() {
		pwd, _ := os.Getwd()
		logger := logr.FromContextOrDiscard(cmd.Context())
		logger.V(1).Info(fmt.Sprintf("working directory: %q", pwd))
	}()

	if path == "" {
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get absolute path of %q error: %w", path, err)
	}
	if err := os.Chdir(absPath); err != nil {
		return fmt.Errorf("change working directory to %q error: %w", absPath, err)
	}
	// chdir 之后需要更新一下 PWD 变量，否则 os.Getwd 会判断错误
	if err := os.Setenv("PWD", absPath); err != nil {
		return fmt.Errorf("set env PWD to %q error: %w", absPath, err)
	}

	return nil
}
