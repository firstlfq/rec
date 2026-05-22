package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	// 确保 ~/.rec 目录存在
	home, _ := os.UserHomeDir()
	recDir := filepath.Join(home, ".rec")
	os.MkdirAll(recDir, 0755)

	cmd := os.Args[1]

	switch cmd {
	case "help", "--help", "-h":
		printHelp()

	case "today":
		showDay(recDir, time.Now())

	case "yesterday":
		showDay(recDir, time.Now().AddDate(0, 0, -1))

	case "搜", "search", "s":
		if len(os.Args) < 3 {
			fmt.Println("用法: rec 搜 <关键词>")
			return
		}
		keyword := strings.Join(os.Args[2:], " ")
		searchLogs(recDir, keyword)

	case "plan", "p":
		if len(os.Args) < 3 {
			showPlans(recDir)
			return
		}
		task := strings.Join(os.Args[2:], " ")
		addPlan(recDir, task)

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("用法: rec done <任务编号>")
			return
		}
		completePlan(recDir, os.Args[2])

	default:
		// 默认：记录一条日志
		msg := strings.Join(os.Args[1:], " ")
		addLog(recDir, msg)
	}
}

func printHelp() {
	fmt.Println(`rec — 个人日志记忆工具

用法:
  rec <消息>            记录一条日志
  rec today             查看今天的日志
  rec yesterday         查看昨天的日志
  rec 搜 <关键词>        搜索所有日志
  rec plan <任务>        添加工作计划
  rec plan              查看工作计划
  rec done <编号>        标记任务完成
  rec help              显示帮助

数据保存在 ~/.rec/ 目录下，纯文本格式。`)
}

func addLog(recDir string, msg string) {
	now := time.Now()
	filename := filepath.Join(recDir, now.Format("2006-01-02")+".md")
	line := fmt.Sprintf("%s %s\n", now.Format("15:04"), msg)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "写入失败: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString(line)

	fmt.Printf("✅ 已记录 (%s)\n", now.Format("15:04"))
}

func showDay(recDir string, date time.Time) {
	filename := filepath.Join(recDir, date.Format("2006-01-02")+".md")
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("📋 %s\n   (暂无记录)\n", date.Format("2006-01-02"))
		return
	}

	fmt.Printf("📋 %s\n", date.Format("2006-01-02"))
	fmt.Println(strings.Repeat("─", 40))
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			fmt.Println("  " + line)
		}
	}
}

func searchLogs(recDir string, keyword string) {
	entries, _ := os.ReadDir(recDir)
	found := false

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if entry.Name() == "plans.md" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(recDir, entry.Name()))
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		date := strings.TrimSuffix(entry.Name(), ".md")

		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				if !found {
				fmt.Print("🔍 搜索结果:\n\n")
					found = true
				}
				fmt.Printf("  %s  %s\n", date, line)
			}
		}
	}

	if !found {
		fmt.Printf("🔍 未找到包含「%s」的记录", keyword)
	}
}

func addPlan(recDir string, task string) {
	filename := filepath.Join(recDir, "plans.md")

	entries, _ := os.ReadDir(recDir)
	maxID := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "plan-") && strings.HasSuffix(entry.Name(), ".md") {
			var id int
			fmt.Sscanf(entry.Name(), "plan-%d.md", &id)
			if id > maxID {
				maxID = id
			}
		}
	}

	// 读取现有计划，算下一个编号
	existing := ""
	if data, err := os.ReadFile(filename); err == nil {
		existing = string(data)
	}

	nextID := 1
	lines := strings.Split(existing, "\n")
	for _, l := range lines {
		var id int
		if strings.HasPrefix(l, "- [") {
			if _, err := fmt.Sscanf(l, "- [%d]", &id); err == nil {
				if id >= nextID {
					nextID = id + 1
				}
			}
		}
	}

	line := fmt.Sprintf("- [%d] 📌 %s\n", nextID, task)
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(line)

	fmt.Printf("✅ 已添加任务 #%d: %s\n", nextID, task)
}

func showPlans(recDir string) {
	filename := filepath.Join(recDir, "plans.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("📋 工作计划\n  (暂无计划)")
		return
	}

	fmt.Println("📋 工作计划")
	fmt.Println(strings.Repeat("─", 40))
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			fmt.Println("  " + line)
		}
	}
}

func completePlan(recDir string, idStr string) {
	filename := filepath.Join(recDir, "plans.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取计划失败: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		if strings.Contains(line, "[ "+idStr+"]") || strings.Contains(line, "["+idStr+"]") {
			// 替换标记
			lines[i] = strings.Replace(line, "📌", "✅", 1)
			found = true
			fmt.Printf("✅ 任务 #%s 已完成\n", idStr)
			break
		}
	}

	if !found {
		// 尝试在行尾追加完成标记
		for i, line := range lines {
			var id int
			if _, err := fmt.Sscanf(line, "- [%d]", &id); err == nil && fmt.Sprintf("%d", id) == idStr {
				lines[i] = line + " ✅"
				found = true
				fmt.Printf("✅ 任务 #%s 已完成\n", idStr)
				break
			}
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "未找到任务 #%s\n", idStr)
		return
	}

	os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
