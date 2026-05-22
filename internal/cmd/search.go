package cmd

import (
	"fmt"
	"strings"

	"github.com/firstlfq/rec/internal/storage"
)

// SearchCmd 搜索：rec 搜 <关键词>
type SearchCmd struct{}

func (c *SearchCmd) Name() string      { return "搜" }
func (c *SearchCmd) Aliases() []string { return []string{"search", "s"} }
func (c *SearchCmd) Execute(args []string, store storage.Storage) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: rec 搜 <关键词>")
	}

	keyword := strings.Join(args, " ")
	results, err := store.SearchEntries(keyword)
	if err != nil {
		return fmt.Errorf("搜索失败: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("🔍 未找到包含「%s」的记录\n", keyword)
		return nil
	}

	fmt.Printf("🔍 搜索「%s」共 %d 条结果:\n\n", keyword, len(results))
	for _, r := range results {
		if r.Time != "" {
			fmt.Printf("  %s  %s  %s\n", r.Date, r.Time, r.Content)
		} else {
			fmt.Printf("  %s  %s\n", r.Date, r.Content)
		}
	}
	return nil
}
