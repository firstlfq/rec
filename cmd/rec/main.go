package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/firstlfq/rec/internal/app"
	"github.com/firstlfq/rec/internal/cmd"
	"github.com/firstlfq/rec/internal/storage"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取 HOME 目录失败: %v\n", err)
		os.Exit(1)
	}

	recDir := filepath.Join(home, ".rec")
	store := storage.NewFSStorage(recDir)

	application := app.New("rec", store)

	// 注册子命令
	application.Register(&cmd.HelpCmd{})
	application.Register(&cmd.LogCmd{})
	application.Register(&cmd.TodayCmd{})
	application.Register(&cmd.YesterdayCmd{})
	application.Register(&cmd.SearchCmd{})
	application.Register(&cmd.PlanCmd{})
	application.Register(&cmd.DoneCmd{})

	if err := application.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}
}
