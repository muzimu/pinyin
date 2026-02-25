// py - 汉字转拼音首字母命令行工具
// 用法: echo "张三" | py
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
	"github.com/spf13/cobra"
)

var (
	upper      bool
	lower      bool
	keepNonHan bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "py",
		Short: "将汉字转换为拼音首字母",
		Long: `py 读取标准输入中的汉字，输出其拼音首字母。

示例:
  echo "张三" | py
  echo "张三" | py -u
  jq -r '.name' test.json | py`,
		RunE: run,
	}

	rootCmd.Flags().BoolVarP(&upper, "upper", "u", false, "首字母大写")
	rootCmd.Flags().BoolVarP(&lower, "lower", "l", true, "首字母小写（默认开启）")
	rootCmd.Flags().BoolVarP(&keepNonHan, "keep-non-han", "k", true, "保留非汉字字符（默认开启）")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// upper 显式指定时覆盖默认 lower
	if cmd.Flags().Changed("upper") && upper {
		lower = false
	}

	a := pinyin.NewArgs()
	a.Style = pinyin.FirstLetter

	scanner := bufio.NewScanner(os.Stdin)
	// 支持大输入场景，设置 1MB 缓冲区
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, len(buf))

	for scanner.Scan() {
		line := scanner.Text()
		result := convertLine(line, a)
		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取输入失败: %v\n", err)
		return err
	}

	return nil
}

// convertLine 将一行文本中的汉字转换为拼音首字母
func convertLine(line string, a pinyin.Args) string {
	var sb strings.Builder

	for _, r := range line {
		if unicode.Is(unicode.Han, r) {
			// 汉字：转换为拼音首字母
			p := pinyin.SinglePinyin(r, a)
			if len(p) > 0 && len(p[0]) > 0 {
				letter := string(p[0][0])
				if upper {
					letter = strings.ToUpper(letter)
				} else {
					letter = strings.ToLower(letter)
				}
				sb.WriteString(letter)
			}
		} else if keepNonHan {
			// 非汉字：根据配置决定是否保留
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
