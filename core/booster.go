package core

import (
	"autoboost/config"
	"autoboost/discord"
	"autoboost/storage"
	"autoboost/ui"
	"autoboost/utils"
	"fmt"
	"time"
)

type BoostResult struct {
	Timestamp time.Time
	Token     string
	Username  string
	UserID    string
	GuildID   string
	SlotID    string
	Action    string
	Success   bool
	Error     string
}

func ExecutePlan(plan *BoostPlan, cfg *config.Config, proxyFn func() string) []BoostResult {
	var results []BoostResult

	fmt.Println()
	fmt.Printf("  %s⚡ Bắt đầu tiến trình boost → Máy chủ %s%s\n\n",
		utils.Bold+utils.BrightMagenta, plan.TargetGuildID, utils.Reset)

	for idx, work := range plan.Work {
		result := BoostResult{
			Timestamp: time.Now(),
			Token:     work.Account.Token,
			GuildID:   plan.TargetGuildID,
			SlotID:    work.Slot.ID,
		}

		username := "không rõ"
		userID := ""
		if work.Account.User != nil {
			username = work.Account.User.DisplayName()
			userID = work.Account.User.ID
		}
		result.Username = username
		result.UserID = userID

		opts := buildClientOpts(cfg, proxyFn)
		client := discord.New(work.Account.Token, opts...)

		var err error
		if work.NeedsXfer {
			result.Action = "transfer"
			fmt.Printf("  %s[%d/%d]%s %s%-20s%s  đang chuyển boost...\n",
				utils.Dim, idx+1, plan.TotalScheduled, utils.Reset,
				utils.BrightYellow, username, utils.Reset)
			err = client.TransferBoost(work.Slot, plan.TargetGuildID)
		} else {
			result.Action = "boost"
			fmt.Printf("  %s[%d/%d]%s %s%-20s%s  đang boost...\n",
				utils.Dim, idx+1, plan.TotalScheduled, utils.Reset,
				utils.BrightCyan, username, utils.Reset)
			err = client.ApplyBoostSingle(plan.TargetGuildID, work.Slot.ID)
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Success = true
		}

		ui.PrintBoostResult(username, plan.TargetGuildID, work.Slot.ID, result.Action, result.Success, result.Error)
		results = append(results, result)

		if idx < len(plan.Work)-1 {
			utils.RandomDelay(cfg.MinBoostDelay, cfg.MaxBoostDelay)
		}
	}

	success, failed := 0, 0
	for _, r := range results {
		if r.Success {
			success++
		} else {
			failed++
		}
	}

	ui.PrintFinalSummary(success, failed, plan.SkippedCooldown)
	return results
}

func ToLogEntries(results []BoostResult) []storage.BoostLogEntry {
	entries := make([]storage.BoostLogEntry, len(results))
	for i, r := range results {
		entries[i] = storage.BoostLogEntry{
			Timestamp: r.Timestamp,
			Token:     r.Token,
			Username:  r.Username,
			UserID:    r.UserID,
			GuildID:   r.GuildID,
			SlotID:    r.SlotID,
			Action:    r.Action,
			Success:   r.Success,
			Error:     r.Error,
		}
	}
	return entries
}

func buildClientOpts(cfg *config.Config, proxyFn func() string) []discord.ClientOption {
	opts := []discord.ClientOption{
		discord.WithUserAgent(cfg.UserAgent),
		discord.WithMaxRetries(cfg.MaxRetries),
		discord.WithBaseURL(cfg.BaseURL),
	}
	if proxyFn != nil {
		if pURL := proxyFn(); pURL != "" {
			opts = append(opts, discord.WithProxy(pURL))
		}
	}
	return opts
}
