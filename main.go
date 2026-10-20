package main

import (
	"autoboost/config"
	"autoboost/core"
	"autoboost/discord"
	"autoboost/storage"
	"autoboost/ui"
	"autoboost/utils"
	"fmt"
	"os"
)

func main() {

	enableWindowsANSI()

	ui.PrintBanner()

	cfg, err := config.FromFile("config.json")
	if err != nil {
		fatal("Đọc file config.json", err)
	}

	tokensFile := cfg.TokensFile
	if _, statErr := os.Stat(tokensFile); os.IsNotExist(statErr) {
		fmt.Printf("  %s⚠ Không tìm thấy file %s%s\n", utils.Yellow, tokensFile, utils.Reset)
		tokensFile = ui.AskTokensFile()
	}

	spin := ui.NewSpinner("Đang đọc token từ " + tokensFile)
	spin.Start()
	tokens, err := storage.LoadTokens(tokensFile)
	if err != nil {
		spin.StopFail(err.Error())
		os.Exit(1)
	}
	spin.StopSuccess(fmt.Sprintf("Đã nạp thành công %d token", len(tokens)))

	var proxyFn func() string
	if cfg.ProxyFile != "" {
		pool, pErr := utils.LoadProxies(cfg.ProxyFile)
		if pErr == nil {
			fmt.Printf("  %s✔ Đã nạp %d proxy vào hệ thống%s\n\n", utils.BrightGreen, pool.Len(), utils.Reset)
			proxyFn = pool.Next
		} else {
			fmt.Printf("  %s⚠ Lỗi file proxy (sẽ chạy trực tiếp không qua proxy): %v%s\n\n", utils.Yellow, pErr, utils.Reset)
		}
	}

	fmt.Printf("  %sĐang kiểm tra %d tài khoản với tối đa %d luồng song song...%s\n\n",
		utils.Dim, len(tokens), cfg.MaxConcurrent, utils.Reset)

	checkSpin := ui.NewSpinner(fmt.Sprintf("Đang kiểm tra 0 / %d", len(tokens)))
	checkSpin.Start()
	accounts := core.CheckAllAccounts(tokens, cfg, proxyFn)
	checkSpin.StopSuccess(fmt.Sprintf("Đã kiểm tra xong %d tài khoản", len(accounts)))

	ui.PrintAccountTable(accounts)
	ui.PrintSummary(accounts)
	ui.PrintCooldownInfo(accounts)

	guildID := ui.AskGuildID()
	if guildID == "" {
		fmt.Println(utils.Colorize(utils.Red, "  Chưa nhập ID máy chủ. Đang thoát chương trình."))
		os.Exit(0)
	}

	validateSpin := ui.NewSpinner("Đang kiểm tra máy chủ " + guildID)
	validateSpin.Start()
	guildOK := validateGuild(guildID, accounts, cfg, proxyFn)
	if !guildOK {
		validateSpin.StopFail("Không tìm thấy máy chủ hoặc tài khoản chưa tham gia server")
		if !ui.Confirm("Vẫn tiếp tục thực hiện chứ?") {
			os.Exit(0)
		}
	} else {
		validateSpin.StopSuccess("Xác thực máy chủ thành công ✓")
	}

	totalBusy := 0
	for _, acc := range accounts {
		totalBusy += len(acc.BusySlots)
	}
	includeTransfer := ui.AskHandleTransfer(totalBusy)

	maxAvailable := core.AvailableCount(accounts, includeTransfer)
	if maxAvailable == 0 {
		fmt.Println(utils.Colorize(utils.Red, "\n  Không còn slot nào khả dụng để boost. Đang thoát chương trình."))
		os.Exit(0)
	}

	amount := ui.AskBoostAmount(maxAvailable)
	if amount == 0 {
		fmt.Println(utils.Colorize(utils.Red, "\n  Đã hủy thao tác."))
		os.Exit(0)
	}
	if amount == -1 {
		amount = maxAvailable
	}

	plan := core.BuildPlan(accounts, guildID, amount, includeTransfer)
	fmt.Printf("\n  %sKế hoạch: Lên lịch %d lượt boost → máy chủ %s%s\n",
		utils.BrightCyan, plan.TotalScheduled, guildID, utils.Reset)

	if !ui.Confirm("Bắt đầu thực thi kế hoạch ngay bây giờ?") {
		fmt.Println(utils.Colorize(utils.Yellow, "  Đã hủy thao tác."))
		os.Exit(0)
	}

	results := core.ExecutePlan(plan, cfg, proxyFn)

	entries := core.ToLogEntries(results)
	exports, exportErr := storage.WriteAllExports(entries, guildID, cfg.LogDir)
	fmt.Println()
	fmt.Printf("  %s📂 Tệp báo cáo xuất ra:%s\n", utils.Bold+utils.BrightCyan, utils.Reset)
	if exports != nil {
		if exports.JSONPath != "" {
			fmt.Printf("  %s├─ JSON (dữ liệu thô):     %s%s%s\n", utils.Dim, utils.Reset, exports.JSONPath, utils.Reset)
		}
		if exports.TXTPath != "" {
			fmt.Printf("  %s├─ TXT  (báo cáo chi tiết): %s%s%s\n", utils.Dim, utils.Reset, exports.TXTPath, utils.Reset)
		}
		if exports.CSVPath != "" {
			fmt.Printf("  %s└─ CSV  (bảng tính Excel):  %s%s%s\n", utils.Dim, utils.Reset, exports.CSVPath, utils.Reset)
		}
	}
	if exportErr != nil {
		fmt.Printf("  %s⚠ Có lỗi khi xuất một số tệp: %v%s\n", utils.Yellow, exportErr, utils.Reset)
	}
	fmt.Println()
}

func validateGuild(guildID string, accounts []*discord.AccountStatus, cfg *config.Config, proxyFn func() string) bool {
	for _, acc := range accounts {
		if acc.Error != "" || !acc.HasNitro {
			continue
		}
		opts := []discord.ClientOption{
			discord.WithUserAgent(cfg.UserAgent),
			discord.WithMaxRetries(1),
			discord.WithBaseURL(cfg.BaseURL),
		}
		if proxyFn != nil {
			if pURL := proxyFn(); pURL != "" {
				opts = append(opts, discord.WithProxy(pURL))
			}
		}
		client := discord.New(acc.Token, opts...)
		if _, err := client.GetGuildInfo(guildID); err == nil {
			return true
		}
		break
	}
	return false
}

func fatal(ctx string, err error) {
	fmt.Fprintf(os.Stderr, "\n  %s✘ Lỗi nghiêm trọng — %s: %v%s\n\n", utils.BrightRed, ctx, err, utils.Reset)
	os.Exit(1)
}

func enableWindowsANSI() {

}
