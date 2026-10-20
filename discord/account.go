package discord

import (
	"encoding/json"
	"fmt"
	"time"
)

func (c *Client) GetMe() (*User, error) {
	body, _, err := c.Get("/users/@me")
	if err != nil {
		return nil, fmt.Errorf("GetMe: %w", err)
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("GetMe decode: %w", err)
	}
	return &user, nil
}

func (c *Client) GetSubscriptions() ([]Subscription, error) {
	body, _, err := c.Get("/users/@me/billing/subscriptions")
	if err != nil {
		return nil, fmt.Errorf("GetSubscriptions: %w", err)
	}
	var subs []Subscription
	if err := json.Unmarshal(body, &subs); err != nil {
		return nil, fmt.Errorf("GetSubscriptions decode: %w", err)
	}
	return subs, nil
}

func (c *Client) GetBoostSlots() ([]BoostSlot, error) {
	body, _, err := c.Get("/users/@me/guilds/premium/subscription-slots")
	if err != nil {
		return nil, fmt.Errorf("GetBoostSlots: %w", err)
	}
	var slots []BoostSlot
	if err := json.Unmarshal(body, &slots); err != nil {
		return nil, fmt.Errorf("GetBoostSlots decode: %w", err)
	}
	return slots, nil
}

func (c *Client) CheckAccount(token string) *AccountStatus {
	status := &AccountStatus{
		Token:     token,
		CheckedAt: time.Now(),
	}

	user, err := c.GetMe()
	if err != nil {
		status.Error = fmt.Sprintf("GetMe failed: %v", err)
		return status
	}
	status.User = user
	status.NitroType = user.PremiumType
	status.HasNitro = user.PremiumType != PremiumNone

	if status.HasNitro {
		slots, err := c.GetBoostSlots()
		if err != nil {

			status.Error = fmt.Sprintf("GetBoostSlots failed: %v", err)
		} else {
			status.Slots = slots
			classifySlots(status)
		}
	}

	return status
}

func classifySlots(s *AccountStatus) {
	s.ReadySlots = nil
	s.BusySlots = nil
	s.CooldownSlots = nil

	for i := range s.Slots {
		slot := &s.Slots[i]
		switch slot.Status() {
		case SlotReady:
			s.ReadySlots = append(s.ReadySlots, slot)
		case SlotBoosting:
			s.BusySlots = append(s.BusySlots, slot)
		case SlotCooldown:
			s.CooldownSlots = append(s.CooldownSlots, slot)
		}
	}
}
