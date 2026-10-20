package ui

import (
	"autoboost/discord"
	"autoboost/utils"
	"fmt"
	"strings"
)

func PrintBanner() {
	fmt.Print(utils.BrightMagenta + `
  ██╗     ██╗   ██╗███╗   ███╗██╗███╗   ██╗██╗ ██████╗ ██╗   ██╗███████╗
  ██║     ██║   ██║████╗ ████║██║████╗  ██║██║██╔═══██╗██║   ██║██╔════╝
  ██║     ██║   ██║██╔████╔██║██║██╔██╗ ██║██║██║   ██║██║   ██║███████╗
  ██║     ██║   ██║██║╚██╔╝██║██║██║╚██╗██║██║██║   ██║██║   ██║╚════██║
  ███████╗╚██████╔╝██║ ╚═╝ ██║██║██║ ╚████║██║╚██████╔╝╚██████╔╝███████║
  ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═╝ ╚═════╝  ╚═════╝ ╚══════╝
` + utils.Reset + utils.Cyan + `
  ▸ Tự Động Hóa Boost Server Discord
  ▸ Đa tài khoản │ Tự động chuyển boost (Transfer) │ Chống Rate-limit
` + utils.Reset + utils.Dim + `  v1.0.0  phát triển bởi Luminious` + utils.Reset + `
  ` + utils.Dim + `─────────────────────────────────────────────────────────────────────` + utils.Reset + `
`)
}

func PrintAccountTable(accounts []*discord.AccountStatus) {
	fmt.Printf("\n  %s📋 Trạng Thái Danh Sách Tài Khoản%s\n\n", utils.Bold+utils.BrightCyan, utils.Reset)

	hdr := fmt.Sprintf("  %-4s %-22s %-16s %-6s %-9s %-10s %-8s %s",
		"#", "Tài khoản", "Nitro", "Tổng", "Sẵn sàng", "Đang dùng", "CD", "Trạng thái")
	fmt.Println(utils.Bold + utils.BrightWhite + hdr + utils.Reset)
	fmt.Println("  " + strings.Repeat("─", 94))

	for i, acc := range accounts {
		num := fmt.Sprintf("%d", i+1)

		name := "—"
		if acc.User != nil {
			name = acc.User.DisplayName()
		}
		name = truncate(name, 20)

		nitro := utils.Red + "❌ Không" + utils.Reset
		if acc.HasNitro {
			nitro = utils.BrightGreen + "✅ " + acc.NitroType.String() + utils.Reset
		}

		total, ready, busy, cd := "—", "—", "—", "—"
		if acc.HasNitro && acc.Error == "" {
			total = fmt.Sprintf("%d", len(acc.Slots))
			ready = utils.BrightGreen + fmt.Sprintf("%d", len(acc.ReadySlots)) + utils.Reset
			busy = utils.BrightYellow + fmt.Sprintf("%d", len(acc.BusySlots)) + utils.Reset
			cd = utils.BrightRed + fmt.Sprintf("%d", len(acc.CooldownSlots)) + utils.Reset
		}

		status := slotStatusLabel(acc)
		if acc.Error != "" {
			status = utils.Yellow + "⚠ " + truncate(acc.Error, 28) + utils.Reset
		}

		fmt.Printf("  %-4s %-22s %-25s %-6s %-18s %-19s %-17s %s\n",
			num, name, nitro, total, ready, busy, cd, status)
	}
	fmt.Println("  " + strings.Repeat("─", 94))
}

func PrintSummary(accounts []*discord.AccountStatus) {
	totalReady, totalBusy, totalCooldown, noNitro, errored := 0, 0, 0, 0, 0
	for _, acc := range accounts {
		if acc.Error != "" {
			errored++
			continue
		}
		if !acc.HasNitro {
			noNitro++
			continue
		}
		totalReady += len(acc.ReadySlots)
		totalBusy += len(acc.BusySlots)
		totalCooldown += len(acc.CooldownSlots)
	}

	fmt.Printf("\n  %sThống Kê Tổng Quan:%s\n", utils.Bold+utils.BrightWhite, utils.Reset)
	fmt.Printf("  ├─ %s✔ Sẵn sàng boost:    %s%d slot%s\n", utils.BrightGreen, utils.Bold, totalReady, utils.Reset)
	fmt.Printf("  ├─ %s↩ Cần chuyển server: %s%d slot%s (đang boost server khác)\n", utils.BrightYellow, utils.Bold, totalBusy, utils.Reset)
	fmt.Printf("  ├─ %s⏳ Đang chờ cooldown: %s%d slot%s\n", utils.BrightRed, utils.Bold, totalCooldown, utils.Reset)
	fmt.Printf("  ├─ %s— Không có Nitro:    %s%d tài khoản%s\n", utils.Dim, utils.Bold, noNitro, utils.Reset)
	fmt.Printf("  └─ %s⚠ Lỗi / Không vào đc: %s%d tài khoản%s\n\n", utils.Yellow, utils.Bold, errored, utils.Reset)
}

func PrintCooldownInfo(accounts []*discord.AccountStatus) {
	hasCooldown := false
	for _, acc := range accounts {
		if len(acc.CooldownSlots) > 0 {
			hasCooldown = true
			break
		}
	}
	if !hasCooldown {
		return
	}

	fmt.Printf("  %s⏳ Chi Tiết Thời Gian Chờ Cooldown%s\n\n", utils.Bold+utils.BrightYellow, utils.Reset)
	for _, acc := range accounts {
		for _, slot := range acc.CooldownSlots {
			name := "?"
			if acc.User != nil {
				name = acc.User.DisplayName()
			}
			remaining := slot.CooldownRemaining()
			fmt.Printf("  %s%-22s%s  slot: %s...  còn lại: %s%s%s\n",
				utils.Dim, name, utils.Reset,
				slot.ID[:8],
				utils.BrightYellow, remaining, utils.Reset,
			)
		}
	}
	fmt.Println()
}

func PrintBoostResult(username, guildID, slotID, action string, success bool, errMsg string) {
	icon := utils.BrightGreen + "✔" + utils.Reset
	detail := ""
	if !success {
		icon = utils.BrightRed + "✘" + utils.Reset
		detail = "  " + utils.Red + errMsg + utils.Reset
	}
	shortSlot := slotID
	if len(slotID) > 8 {
		shortSlot = slotID[:8] + "..."
	}
	actionVN := "boost mới"
	if action == "transfer" {
		actionVN = "chuyển server"
	}
	fmt.Printf("  %s  %-20s  server: %s  slot: %s  [%s]%s\n",
		icon, username, guildID, shortSlot, actionVN, detail)
}

func PrintFinalSummary(success, failed, skipped int) {
	fmt.Println()
	fmt.Println("  " + utils.Bold + strings.Repeat("─", 35) + utils.Reset)
	fmt.Printf("  %s✔ Đã boost:  %s%s%d slot%s\n", utils.BrightGreen, utils.Reset, utils.Bold, success, utils.Reset)
	fmt.Printf("  %s✘ Thất bại:  %s%s%d slot%s\n", utils.BrightRed, utils.Reset, utils.Bold, failed, utils.Reset)
	fmt.Printf("  %s⏭ Bỏ qua:    %s%s%d slot%s (đang cooldown)\n", utils.Yellow, utils.Reset, utils.Bold, skipped, utils.Reset)
	fmt.Println("  " + utils.Bold + strings.Repeat("─", 35) + utils.Reset)
	fmt.Println()
}

func slotStatusLabel(acc *discord.AccountStatus) string {
	if acc.Error != "" {
		return utils.Yellow + "⚠ LỖI" + utils.Reset
	}
	if !acc.HasNitro {
		return utils.Dim + "— Không Nitro" + utils.Reset
	}
	parts := []string{}
	if n := len(acc.ReadySlots); n > 0 {
		parts = append(parts, utils.BrightGreen+fmt.Sprintf("%d SẴN SÀNG", n)+utils.Reset)
	}
	if n := len(acc.BusySlots); n > 0 {
		parts = append(parts, utils.BrightYellow+fmt.Sprintf("%d ĐANG DÙNG", n)+utils.Reset)
	}
	if n := len(acc.CooldownSlots); n > 0 {
		parts = append(parts, utils.BrightRed+fmt.Sprintf("%d CD", n)+utils.Reset)
	}
	if len(parts) == 0 {
		return utils.Dim + "— Không có slot" + utils.Reset
	}
	return strings.Join(parts, " / ")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-2] + ".."
}
