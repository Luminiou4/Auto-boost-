package ui

import (
	"autoboost/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func AskGuildID() string {
	prompt := promptui.Prompt{
		Label: utils.BrightCyan + "🎯 Nhập ID Máy Chủ (Guild ID)" + utils.Reset,
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if len(input) < 17 || len(input) > 20 {
				return fmt.Errorf("ID máy chủ phải có độ dài từ 17–20 chữ số")
			}
			for _, ch := range input {
				if ch < '0' || ch > '9' {
					return fmt.Errorf("ID máy chủ chỉ được chứa chữ số")
				}
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(utils.Colorize(utils.Red, "\n  Đã hủy thao tác."))
		return ""
	}
	return strings.TrimSpace(result)
}

func AskBoostAmount(maxAvailable int) int {
	fmt.Printf("  %sSố slot khả dụng: %s%s%d%s\n",
		utils.Dim, utils.Reset, utils.Bold, maxAvailable, utils.Reset)

	prompt := promptui.Prompt{
		Label: utils.BrightCyan + "⚡ Số lượng boost muốn dùng (1-" + strconv.Itoa(maxAvailable) + " hoặc all)" + utils.Reset,
		Validate: func(input string) error {
			input = strings.TrimSpace(strings.ToLower(input))
			if input == "all" {
				return nil
			}
			n, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("vui lòng nhập một số hoặc chữ \"all\"")
			}
			if n < 1 {
				return fmt.Errorf("số lượng tối thiểu là 1")
			}
			if n > maxAvailable {
				return fmt.Errorf("chỉ có tối đa %d slot khả dụng", maxAvailable)
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0
	}

	result = strings.TrimSpace(strings.ToLower(result))
	if result == "all" {
		return -1
	}
	n, _ := strconv.Atoi(result)
	return n
}

func AskHandleTransfer(busySlots int) bool {
	if busySlots == 0 {
		return false
	}
	fmt.Printf("\n  %s%d slot%s hiện đang boost ở các máy chủ khác.\n",
		utils.BrightYellow, busySlots, utils.Reset)

	prompt := promptui.Select{
		Label: utils.BrightYellow + "↩ Bạn có muốn chuyển (Transfer) các boost này sang máy chủ mục tiêu không?" + utils.Reset,
		Items: []string{
			"Có  — Rút boost khỏi server cũ để chuyển sang server mới",
			"Không — Bỏ qua slot đang bận, chỉ dùng slot ĐANG SẴN SÀNG",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return false
	}
	return i == 0
}

func Confirm(message string) bool {
	prompt := promptui.Select{
		Label: message,
		Items: []string{"Có (Yes)", "Không (No)"},
	}
	i, _, err := prompt.Run()
	if err != nil {
		return false
	}
	return i == 0
}

func AskTokensFile() string {
	prompt := promptui.Prompt{
		Label:   utils.BrightCyan + "📂 Nhập đường dẫn file chứa token" + utils.Reset,
		Default: "tokens.txt",
	}
	result, err := prompt.Run()
	if err != nil {
		return "tokens.txt"
	}
	return strings.TrimSpace(result)
}
