package discord

import "fmt"

import "time"

type PremiumType int

const (
	PremiumNone    PremiumType = 0
	PremiumClassic PremiumType = 1
	PremiumNitro   PremiumType = 2
	PremiumBasic   PremiumType = 3
)

func (p PremiumType) String() string {
	switch p {
	case PremiumClassic:
		return "Nitro Classic"
	case PremiumNitro:
		return "Nitro"
	case PremiumBasic:
		return "Nitro Basic"
	default:
		return "None"
	}
}

func (p PremiumType) BoostSlots() int {
	switch p {
	case PremiumNitro:
		return 2
	case PremiumClassic, PremiumBasic:
		return 0
	default:
		return 0
	}
}

type SlotStatus string

const (
	SlotReady    SlotStatus = "READY"
	SlotBoosting SlotStatus = "BOOSTING"
	SlotCooldown SlotStatus = "COOLDOWN"
)

type User struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	GlobalName    string      `json:"global_name"`
	Email         string      `json:"email"`
	Verified      bool        `json:"verified"`
	MFAEnabled    bool        `json:"mfa_enabled"`
	PremiumType   PremiumType `json:"premium_type"`
	Phone         string      `json:"phone"`
}

func (u *User) DisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	if u.Discriminator != "0" && u.Discriminator != "" {
		return u.Username + "#" + u.Discriminator
	}
	return "@" + u.Username
}

type Subscription struct {
	ID            string     `json:"id"`
	Type          int        `json:"type"`
	Status        int        `json:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
}

type PremiumGuildSubscription struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
	Ended   bool   `json:"ended"`
}

type BoostSlot struct {
	ID                       string                    `json:"id"`
	SubscriptionID           string                    `json:"subscription_id"`
	Canceled                 bool                      `json:"canceled"`
	CooldownEndsAt           *time.Time                `json:"cooldown_ends_at"`
	PremiumGuildSubscription *PremiumGuildSubscription `json:"premium_guild_subscription"`
}

func (s *BoostSlot) Status() SlotStatus {
	now := time.Now()
	if s.CooldownEndsAt != nil && s.CooldownEndsAt.After(now) {
		return SlotCooldown
	}
	if s.PremiumGuildSubscription != nil && !s.PremiumGuildSubscription.Ended {
		return SlotBoosting
	}
	return SlotReady
}

func (s *BoostSlot) CooldownRemaining() string {
	if s.CooldownEndsAt == nil {
		return ""
	}
	remaining := time.Until(*s.CooldownEndsAt)
	if remaining <= 0 {
		return ""
	}
	days := int(remaining.Hours()) / 24
	hours := int(remaining.Hours()) % 24
	minutes := int(remaining.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

type AccountStatus struct {
	Token          string
	User           *User
	HasNitro       bool
	NitroType      PremiumType
	Slots          []BoostSlot
	ReadySlots     []*BoostSlot
	BusySlots      []*BoostSlot
	CooldownSlots  []*BoostSlot
	Error          string
	CheckedAt      time.Time
}

func (a *AccountStatus) TotalAvailable(includeTransfer bool) int {
	n := len(a.ReadySlots)
	if includeTransfer {
		n += len(a.BusySlots)
	}
	return n
}

type Guild struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	PremiumTier       int    `json:"premium_tier"`
	PremiumSubCount   int    `json:"premium_subscription_count"`
}

type BoostRequest struct {
	SlotIDs []string `json:"user_premium_guild_subscription_slot_ids"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Discord API error %d (HTTP %d): %s", e.Code, e.Status, e.Message)
}
