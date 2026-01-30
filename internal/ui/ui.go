package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/mattn/go-runewidth"
)

var (
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
	dimmed  = color.New(color.Faint).SprintFunc()
)

// ClearScreen 清屏
func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// PressAnyKey 等待按任意键继续
func PressAnyKey() {
	fmt.Println()
	fmt.Printf("  %s", dimmed("按 Enter 键返回主菜单..."))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Logo ASCII 艺术 Logo
func Logo() {
	// 盾牌守护者 + 右侧信息居中
	fmt.Println()
	fmt.Printf("  %s\n", cyan("╭━━━━━━━╮"))
	fmt.Printf("  %s    %s %s\n", cyan("┃ ◉   ◉ ┃"), bold("certctl"), dimmed("v1.0.0"))
	fmt.Printf("  %s    %s\n", cyan("┃   ▽   ┃"), dimmed("SSL Certificate Manager"))
	fmt.Printf("  %s\n", cyan("╰┳━━━━━┳╯"))
	fmt.Printf("   %s\n", cyan("┃ ◢◣ ┃"))
	fmt.Printf("   %s\n", cyan("╰━━━━╯"))
	fmt.Println()
}

// Header 显示标题框
func Header(title string) {
	width := 50
	fmt.Println()
	fmt.Printf("  ┌%s┐\n", strings.Repeat("─", width))
	titleLen := displayWidth(title)
	padding := (width - titleLen) / 2
	fmt.Printf("  │%s%s%s│\n", strings.Repeat(" ", padding), cyan(title), strings.Repeat(" ", width-titleLen-padding))
	fmt.Printf("  └%s┘\n", strings.Repeat("─", width))
}

// StatusLine 显示状态行
func StatusLine(label, value string) {
	fmt.Printf("  %s: %s\n", dimmed(label), cyan(value))
}

// Select 交互式选择菜单
func Select(label string, items []string) (int, string, error) {
	fmt.Println()

	prompt := promptui.Select{
		Label: fmt.Sprintf("  %s %s", yellow("?"), bold(label)),
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("  %s {{ . | cyan }}", cyan(">")),
			Inactive: fmt.Sprintf("  %s {{ . }}", dimmed(">")),
			Selected: fmt.Sprintf("  %s {{ . | green }}", green("✔")),
		},
		Size: 10,
	}

	idx, result, err := prompt.Run()
	return idx, result, err
}

// SelectWithExit 带退出选项的选择菜单
func SelectWithExit(label string, items []string, exitText string) (int, string, error) {
	// 添加空行 + 退出选项
	allItems := append(items, "", exitText)

	fmt.Println()

	prompt := promptui.Select{
		Label: fmt.Sprintf("%s %s", yellow("?"), bold(label)),
		Items: allItems,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("{{ if eq . \"\" }}{{ else if eq . \"%s\" }}%s{{ else }}%s {{ . | cyan }}{{ end }}", exitText, yellow(exitText), cyan(">")),
			Inactive: fmt.Sprintf("{{ if eq . \"\" }}{{ else if eq . \"%s\" }}%s{{ else }}%s {{ . }}{{ end }}", exitText, dimmed(exitText), dimmed(">")),
			Selected: fmt.Sprintf("%s {{ . | green }}", green("✔")),
		},
		Size: 10,
	}

	idx, result, err := prompt.Run()
	if result == exitText || result == "" {
		return -1, "", fmt.Errorf("user exit")
	}
	return idx, result, err
}

// Input 交互式输入
func Input(label string, defaultVal string) (string, error) {
	prompt := promptui.Prompt{
		Label:   fmt.Sprintf("  %s %s", cyan("›"), label),
		Default: defaultVal,
	}

	result, err := prompt.Run()
	return result, err
}

// InputSecret 交互式密码输入
func InputSecret(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("  %s %s", cyan("›"), label),
		Mask:  '*',
	}

	result, err := prompt.Run()
	return result, err
}

// ConfirmPrompt 交互式确认
func ConfirmPrompt(label string) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("  %s %s", cyan("›"), label),
		IsConfirm: true,
	}

	result, _ := prompt.Run()
	return strings.ToLower(result) == "y" || result == ""
}

// Success 输出成功信息
func Success(msg string) {
	fmt.Printf("  %s %s\n", green("✔"), msg)
}

// Error 输出错误信息
func Error(msg string) {
	fmt.Printf("  %s %s\n", red("✖"), msg)
}

// Info 输出信息
func Info(msg string) {
	fmt.Printf("  %s %s\n", cyan("ℹ"), msg)
}

// Warning 输出警告信息
func Warning(msg string) {
	fmt.Printf("  %s %s\n", yellow("⚠"), msg)
}

// Title 输出标题
func Title(msg string) {
	fmt.Printf("\n  %s\n", bold(msg))
}

// NewSpinner 创建一个新的 spinner
func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "  "
	s.Suffix = " " + msg
	s.Color("cyan")
	return s
}

// Box 输出一个框
func Box(title string, lines []string) {
	width := 60

	titleWidth := displayWidth(title)
	topPadding := width - titleWidth - 5
	if topPadding < 0 {
		topPadding = 0
	}
	fmt.Printf("\n  ┌─ %s %s┐\n", cyan(title), dimmed(strings.Repeat("─", topPadding)))

	fmt.Println("  │" + strings.Repeat(" ", width) + "│")
	for _, line := range lines {
		padding := width - displayWidth(stripANSI(line)) - 2
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("  │  %s%s│\n", line, strings.Repeat(" ", padding))
	}
	fmt.Println("  │" + strings.Repeat(" ", width) + "│")
	fmt.Printf("  └%s┘\n", dimmed(strings.Repeat("─", width)))
}

// AIBox 美化 AI 诊断输出
func AIBox(content string) {
	boxWidth := 68 // 内容区域宽度
	lines := strings.Split(content, "\n")

	// 打印一行内容的辅助函数
	printLine := func(text string, colorFunc func(...interface{}) string) {
		textWidth := runewidth.StringWidth(text)
		padding := boxWidth - textWidth
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("  %s %s%s %s\n", cyan("│"), colorFunc(text), strings.Repeat(" ", padding), cyan("│"))
	}

	// 无颜色函数
	noColor := func(a ...interface{}) string {
		if len(a) > 0 {
			return fmt.Sprint(a...)
		}
		return ""
	}

	// 顶部边框
	fmt.Println()
	fmt.Printf("  %s\n", cyan("┌"+strings.Repeat("─", boxWidth+2)+"┐"))
	
	// 标题行
	printLine("🤖 AI 诊断", bold)
	
	// 分隔线
	fmt.Printf("  %s\n", cyan("├"+strings.Repeat("─", boxWidth+2)+"┤"))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Printf("  %s%s%s\n", cyan("│"), strings.Repeat(" ", boxWidth+2), cyan("│"))
			continue
		}

		// 确定颜色函数和缩进
		colorFunc := noColor
		indent := "   " // 默认续行缩进
		if strings.HasPrefix(line, "🔍") {
			colorFunc = yellow
		} else if strings.HasPrefix(line, "✅") {
			colorFunc = green
		} else if strings.HasPrefix(line, "💡") {
			colorFunc = magenta
		} else if len(line) >= 2 && line[0] >= '1' && line[0] <= '9' && line[1] == '.' {
			// 检测数字列表项 (1. 到 9.)
			colorFunc = dimmed
			line = "  " + line // 缩进
			indent = "    " // 列表项续行缩进更多
		}

		// 对长行进行换行处理，考虑续行缩进
		firstLineWidth := boxWidth
		continueLineWidth := boxWidth - runewidth.StringWidth(indent)
		
		wrappedLines := wrapTextWithIndent(line, firstLineWidth, continueLineWidth)
		for i, wline := range wrappedLines {
			if i > 0 {
				wline = indent + wline
			}
			printLine(wline, colorFunc)
		}
	}

	// 底部边框
	fmt.Printf("  %s\n", cyan("└"+strings.Repeat("─", boxWidth+2)+"┘"))
}

// wrapText 将文本按指定宽度换行
func wrapText(text string, maxWidth int) []string {
	return wrapTextWithIndent(text, maxWidth, maxWidth)
}

// wrapTextWithIndent 将文本按指定宽度换行，支持首行和续行不同宽度
func wrapTextWithIndent(text string, firstWidth, continueWidth int) []string {
	if runewidth.StringWidth(text) <= firstWidth {
		return []string{text}
	}

	var result []string
	var current strings.Builder
	currentWidth := 0
	lineNum := 0
	maxWidth := firstWidth

	for _, r := range text {
		rw := runewidth.RuneWidth(r)
		if currentWidth+rw > maxWidth && currentWidth > 0 {
			result = append(result, current.String())
			current.Reset()
			currentWidth = 0
			lineNum++
			maxWidth = continueWidth // 后续行使用续行宽度
		}
		current.WriteRune(r)
		currentWidth += rw
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// stripANSI 去除 ANSI 转义序列
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// displayWidth 计算显示宽度
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r > 127 {
			width += 2
		} else {
			width++
		}
	}
	return width
}

// Prompt 获取用户输入（简单版本）
func Prompt(msg string) string {
	fmt.Printf("  %s %s: ", cyan("›"), msg)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptSecret 获取用户输入（简单版本）
func PromptSecret(msg string) string {
	fmt.Printf("  %s %s: ", cyan("›"), msg)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Confirm 确认提示（简单版本）
func Confirm(msg string) bool {
	fmt.Printf("  %s %s %s ", cyan("›"), msg, dimmed("(Y/n)"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "" || input == "y" || input == "yes"
}

// DomainList 输出域名列表
func DomainList(domains []string) {
	for _, d := range domains {
		fmt.Printf("    %s %s\n", dimmed("•"), d)
	}
}

// DNSRecord 输出 DNS 记录信息
func DNSRecord(recordName, recordType, recordValue, fullRecord string) {
	lines := []string{
		"请在 DNS 服务商添加以下 TXT 记录:",
		"",
		fmt.Sprintf("%s  %s", dimmed("主机记录"), bold(recordName)),
		fmt.Sprintf("%s  %s", dimmed("记录类型"), recordType),
		fmt.Sprintf("%s  %s", dimmed("记录值  "), cyan(recordValue)),
		"",
		strings.Repeat("─", 40),
		fmt.Sprintf("%s 完整记录: %s", yellow("💡"), fullRecord),
	}
	Box("DNS 验证", lines)
}

// CertResult 输出证书结果
func CertResult(certPath, keyPath, expiry string) {
	// 获取文件名
	certFile := filepath.Base(certPath)
	keyFile := filepath.Base(keyPath)
	dir := filepath.Dir(certPath)

	fmt.Println()
	fmt.Printf("  %s %s\n", green("✔"), bold("证书申请成功"))
	fmt.Println()
	fmt.Printf("    %s %s\n", dimmed("目录:"), dir)
	fmt.Printf("    %s %s\n", dimmed("证书:"), green(certFile))
	fmt.Printf("    %s %s\n", dimmed("私钥:"), green(keyFile))
	fmt.Printf("    %s %s\n", dimmed("有效期:"), bold(expiry))
	fmt.Println()
	fmt.Printf("    %s\n", dimmed("Nginx 配置:"))
	fmt.Printf("    %s %s;\n", cyan("ssl_certificate"), certPath)
	fmt.Printf("    %s %s;\n", cyan("ssl_certificate_key"), keyPath)
	fmt.Println()
}

// Step 显示步骤进度
func Step(current, total int, msg string) {
	circles := ""
	for i := 1; i <= total; i++ {
		if i < current {
			circles += green("●") + " "
		} else if i == current {
			circles += cyan("●") + " "
		} else {
			circles += dimmed("○") + " "
		}
	}
	fmt.Printf("\n  %s %s %s\n", circles, dimmed(fmt.Sprintf("[%d/%d]", current, total)), bold(msg))
}

// Detail 显示详细信息
func Detail(msg string) {
	fmt.Printf("    %s %s\n", dimmed("·"), msg)
}

// ErrorWithHint 显示错误及建议
func ErrorWithHint(msg string, hints []string) {
	fmt.Printf("\n  %s %s\n", red("✖"), msg)
	if len(hints) > 0 {
		fmt.Printf("\n  %s 可能的原因：\n", yellow("💡"))
		for _, hint := range hints {
			fmt.Printf("    %s %s\n", dimmed("•"), hint)
		}
	}
	fmt.Println()
}

// ProgressDone 步骤完成
func ProgressDone(msg string) {
	fmt.Printf("    %s %s\n", green("✔"), msg)
}

// ProgressFail 步骤失败
func ProgressFail(msg string) {
	fmt.Printf("    %s %s\n", red("✖"), msg)
}

// StepProgress 带动态进度的步骤显示
type StepProgress struct {
	total   int
	current int
	spin    *spinner.Spinner
}

// NewStepProgress 创建步骤进度
func NewStepProgress(total int) *StepProgress {
	return &StepProgress{total: total}
}

// Next 进入下一步
func (s *StepProgress) Next(msg string) {
	if s.spin != nil {
		s.spin.Stop()
	}

	s.current++
	circles := ""
	for i := 1; i <= s.total; i++ {
		if i < s.current {
			circles += green("●") + " "
		} else if i == s.current {
			circles += cyan("●") + " "
		} else {
			circles += dimmed("○") + " "
		}
	}
	fmt.Printf("\n  %s %s %s\n", circles, dimmed(fmt.Sprintf("[%d/%d]", s.current, s.total)), bold(msg))
}

// Done 完成所有步骤
func (s *StepProgress) Done(msg string) {
	if s.spin != nil {
		s.spin.Stop()
	}
	circles := ""
	for i := 1; i <= s.total; i++ {
		circles += green("●") + " "
	}
	fmt.Printf("\n  %s %s\n", circles, green("✔ "+msg))
}
