package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/firstlfq/rec/internal/model"
)

// FSStorage 基于文件系统的存储实现
// 数据保存在 ~/.rec/ 目录下
// 日志文件格式：YYYY-MM-DD.md，每行 "HH:MM 内容"
// 计划文件：plans.md
type FSStorage struct {
	dir string
}

// NewFSStorage 创建文件系统存储，dir 为数据目录
func NewFSStorage(dir string) *FSStorage {
	os.MkdirAll(dir, 0755)
	return &FSStorage{dir: dir}
}

// ==============================
// 日志
// ==============================

// AppendEntry 追加一条日志到指定日期的文件
func (s *FSStorage) AppendEntry(date time.Time, content string) error {
	filename := s.logFilename(date)
	line := fmt.Sprintf("%s %s\n", date.Format("15:04"), content)
	return appendToFile(filename, line)
}

// GetEntries 读取指定日期的全部日志
func (s *FSStorage) GetEntries(date time.Time) ([]model.Entry, error) {
	filename := s.logFilename(date)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Entry{}, nil
		}
		return nil, err
	}

	var entries []model.Entry
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 解析 "HH:MM 内容"
		if len(line) < 5 || line[2] != ':' {
			// 格式不对，整行当内容
			entries = append(entries, model.Entry{
				Time:    date,
				Content: line,
			})
			continue
		}
		t, err := time.Parse("15:04", line[:5])
		entryTime := date
		if err == nil {
			entryTime = time.Date(
				date.Year(), date.Month(), date.Day(),
				t.Hour(), t.Minute(), 0, 0, date.Location(),
			)
		}
		entries = append(entries, model.Entry{
			Time:    entryTime,
			Content: strings.TrimSpace(line[5:]),
		})
	}
	return entries, nil
}

// ==============================
// 搜索
// ==============================

// SearchEntries 在所有日志文件中搜索关键词
func (s *FSStorage) SearchEntries(keyword string) ([]model.SearchResult, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	keyword = strings.ToLower(keyword)
	var results []model.SearchResult

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if entry.Name() == "plans.md" || !strings.HasSuffix(name, ".md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.dir, name))
		if err != nil {
			continue
		}

		date := strings.TrimSuffix(name, ".md")
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.Contains(strings.ToLower(line), keyword) {
				timePart := ""
				content := line
				if len(line) >= 5 && line[2] == ':' {
					timePart = line[:5]
					content = strings.TrimSpace(line[5:])
				}
				results = append(results, model.SearchResult{
					Date:    date,
					Time:    timePart,
					Content: content,
				})
			}
		}
	}

	// 按日期排序（最新的在前）
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date > results[j].Date
	})

	return results, nil
}

// ==============================
// 计划
// ==============================

// AddPlan 添加一条计划
func (s *FSStorage) AddPlan(content string) (int, error) {
	filename := filepath.Join(s.dir, "plans.md")

	// 读取现有计划，计算下一个 ID
	nextID := 1
	if data, err := os.ReadFile(filename); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if id := extractPlanID(line); id >= nextID {
				nextID = id + 1
			}
		}
	}

	line := fmt.Sprintf("- [%d] 📌 %s\n", nextID, content)
	if err := appendToFile(filename, line); err != nil {
		return 0, err
	}
	return nextID, nil
}

// GetPlans 获取所有计划
func (s *FSStorage) GetPlans() ([]model.Plan, error) {
	filename := filepath.Join(s.dir, "plans.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Plan{}, nil
		}
		return nil, err
	}

	var plans []model.Plan
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		id := extractPlanID(line)
		if id == 0 {
			continue
		}
		content := extractPlanContent(line)
		done := strings.Contains(line, "✅")
		plans = append(plans, model.Plan{
			ID:      id,
			Content: content,
			Done:    done,
		})
	}
	return plans, nil
}

// CompletePlan 标记计划为已完成
func (s *FSStorage) CompletePlan(id int) error {
	filename := filepath.Join(s.dir, "plans.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	found := false
	prefix := fmt.Sprintf("- [%d]", id)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			// 替换 📌 为 ✅
			lines[i] = strings.Replace(trimmed, "📌", "✅", 1)
			// 如果还没加 ✅，在末尾加
			if !strings.Contains(lines[i], "✅") {
				lines[i] = lines[i] + " ✅"
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到计划 #%d", id)
	}

	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}

// ==============================
// 辅助方法
// ==============================

func (s *FSStorage) logFilename(date time.Time) string {
	return filepath.Join(s.dir, date.Format("2006-01-02")+".md")
}

// 从计划行中提取 ID，如 "- [3] 📌 买牛奶" → 3
func extractPlanID(line string) int {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "- [") {
		return 0
	}
	end := strings.Index(line, "]")
	if end < 3 {
		return 0
	}
	id, err := strconv.Atoi(line[3:end])
	if err != nil {
		return 0
	}
	return id
}

// 从计划行中提取内容，如 "- [3] 📌 买牛奶" → "买牛奶"
func extractPlanContent(line string) string {
	line = strings.TrimSpace(line)
	// 去掉 "- [N] " 前缀
	end := strings.Index(line, "]")
	if end < 0 {
		return line
	}
	content := strings.TrimSpace(line[end+1:])
	// 去掉 emoji 标记
	content = strings.TrimLeft(content, " 📌✅")
	return strings.TrimSpace(content)
}

func appendToFile(filename, line string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}
