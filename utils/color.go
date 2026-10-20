package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	BgBlack  = "\033[40m"
	BgBlue   = "\033[44m"
	BgCyan   = "\033[46m"
)

func Colorize(color, text string) string {
	return color + text + Reset
}

func Sprintf(color, format string, args ...interface{}) string {
	return color + fmt.Sprintf(format, args...) + Reset
}

func RandomDelay(min, max time.Duration) {
	if max <= min {
		time.Sleep(min)
		return
	}
	diff := max - min
	jitter := time.Duration(rand.Int63n(int64(diff)))
	time.Sleep(min + jitter)
}

func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "—"
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
