package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BoostLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	GuildID   string    `json:"guild_id"`
	SlotID    string    `json:"slot_id"`
	Action    string    `json:"action"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Duration  string    `json:"duration_ms,omitempty"`
}

type ExportResult struct {
	JSONPath string
	TXTPath  string
	CSVPath  string
}

func WriteAllExports(entries []BoostLogEntry, guildID string, logDir string) (*ExportResult, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("tạo thư mục log: %w", err)
	}

	ts := time.Now().Format("20060102_150405")
	result := &ExportResult{}
	var errs []string

	if path, err := writeJSON(entries, logDir, ts); err != nil {
		errs = append(errs, "JSON: "+err.Error())
	} else {
		result.JSONPath = path
	}

	if path, err := writeTXT(entries, guildID, logDir, ts); err != nil {
		errs = append(errs, "TXT: "+err.Error())
	} else {
		result.TXTPath = path
	}

	if path, err := writeCSV(entries, logDir, ts); err != nil {
		errs = append(errs, "CSV: "+err.Error())
	} else {
		result.CSVPath = path
	}

	if len(errs) > 0 {
		return result, fmt.Errorf("lỗi xuất file: %s", strings.Join(errs, "; "))
	}
	return result, nil
}

func writeJSON(entries []BoostLogEntry, logDir, ts string) (string, error) {
	filename := filepath.Join(logDir, fmt.Sprintf("luminious_%s.json", ts))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "")
	for _, entry := range entries {
		e := entry
		if len(e.Token) > 20 {
			e.Token = e.Token[:20] + "..."
		}
		if err := enc.Encode(e); err != nil {
			return filename, err
		}
	}
	return filename, nil
}

func writeTXT(entries []BoostLogEntry, guildID, logDir, ts string) (string, error) {
	filename := filepath.Join(logDir, fmt.Sprintf("luminious_%s.txt", ts))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sep := strings.Repeat("─", 74)

	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "  ██╗     ██╗   ██╗███╗   ███╗██╗███╗   ██╗██╗ ██████╗ ██╗   ██╗███████╗\n")
	fmt.Fprintf(f, "  ╚══════╝ ╚═════╝ ╚═╝    ╚═╝╚═╝╚═╝ ╚═══╝╚═╝ ╚═════╝  ╚═════╝ ╚══════╝\n")
	fmt.Fprintf(f, "\n  BÁO CÁO PHIÊN BOOST MÁY CHỦ DISCORD\n")
	fmt.Fprintf(f, "  %s\n", sep)
	fmt.Fprintf(f, "  Thời gian thực hiện : %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "  Máy chủ mục tiêu   : Guild %s\n", guildID)
	fmt.Fprintf(f, "  Tổng số thao tác   : %d\n", len(entries))

	success, failed, xfer := 0, 0, 0
	for _, e := range entries {
		if e.Success {
			success++
		} else {
			failed++
		}
		if e.Action == "transfer" {
			xfer++
		}
	}
	fmt.Fprintf(f, "  Thành công         : %d\n", success)
	fmt.Fprintf(f, "  Thất bại           : %d\n", failed)
	fmt.Fprintf(f, "  Chuyển server      : %d\n", xfer)
	fmt.Fprintf(f, "  %s\n\n", sep)

	fmt.Fprintf(f, "  %-8s %-22s %-14s %-12s %-36s %s\n",
		"Giờ", "Tài khoản", "Hành động", "Kết quả", "Slot ID", "Lỗi / Ghi chú")
	fmt.Fprintf(f, "  %s\n", sep)

	for _, e := range entries {
		icon := "✅"
		errMsg := ""
		if !e.Success {
			icon = "❌"
			errMsg = e.Error
			if len(errMsg) > 35 {
				errMsg = errMsg[:33] + ".."
			}
		}

		actionVN := "boost mới"
		if e.Action == "transfer" {
			actionVN = "chuyển server"
		}

		fmt.Fprintf(f, "  %-8s %s %-20s %-14s %-12s %-36s %s\n",
			e.Timestamp.Format("15:04:05"),
			icon,
			truncStr(e.Username, 20),
			actionVN,
			boolLabelVN(e.Success),
			e.SlotID,
			errMsg,
		)
	}

	fmt.Fprintf(f, "  %s\n\n", sep)

	var failures []BoostLogEntry
	for _, e := range entries {
		if !e.Success {
			failures = append(failures, e)
		}
	}
	if len(failures) > 0 {
		fmt.Fprintf(f, "  CHI TIẾT CÁC LƯỢT THẤT BẠI\n")
		fmt.Fprintf(f, "  %s\n", sep)
		for i, e := range failures {
			actionVN := "boost mới"
			if e.Action == "transfer" {
				actionVN = "chuyển server"
			}
			fmt.Fprintf(f, "  [%d] Tài khoản : %s (ID: %s)\n", i+1, e.Username, e.UserID)
			fmt.Fprintf(f, "      Slot ID   : %s\n", e.SlotID)
			fmt.Fprintf(f, "      Hành động : %s\n", actionVN)
			fmt.Fprintf(f, "      Lỗi       : %s\n\n", e.Error)
		}
	}

	fmt.Fprintf(f, "  Được tạo bởi Luminious v1.0.0\n")
	fmt.Fprintf(f, "  %s\n", sep)

	return filename, nil
}

func writeCSV(entries []BoostLogEntry, logDir, ts string) (string, error) {
	filename := filepath.Join(logDir, fmt.Sprintf("luminious_%s.csv", ts))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	f.WriteString("\xef\xbb\xbf")

	w.Write([]string{
		"Thời gian", "Tài khoản", "User ID", "Guild ID",
		"Slot ID", "Hành động", "Kết quả", "Chi tiết lỗi",
	})

	for _, e := range entries {
		actionVN := "Boost mới"
		if e.Action == "transfer" {
			actionVN = "Chuyển server"
		}
		w.Write([]string{
			e.Timestamp.Format(time.RFC3339),
			e.Username,
			e.UserID,
			e.GuildID,
			e.SlotID,
			actionVN,
			boolLabelVN(e.Success),
			e.Error,
		})
	}

	return filename, nil
}

func boolLabelVN(b bool) string {
	if b {
		return "THÀNH CÔNG"
	}
	return "THẤT BẠI"
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-2] + ".."
}
