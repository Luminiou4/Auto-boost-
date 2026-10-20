package discord

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetGuildInfo(guildID string) (*Guild, error) {
	body, _, err := c.Get(fmt.Sprintf("/guilds/%s?with_counts=false", guildID))
	if err != nil {
		return nil, fmt.Errorf("GetGuildInfo: %w", err)
	}
	var guild Guild
	if err := json.Unmarshal(body, &guild); err != nil {
		return nil, fmt.Errorf("GetGuildInfo decode: %w", err)
	}
	return &guild, nil
}

func (c *Client) ApplyBoost(guildID string, slotIDs []string) error {
	if len(slotIDs) == 0 {
		return nil
	}
	payload := BoostRequest{SlotIDs: slotIDs}
	_, _, err := c.Put(fmt.Sprintf("/guilds/%s/premium/subscriptions", guildID), payload)
	if err != nil {
		return fmt.Errorf("ApplyBoost to guild %s: %w", guildID, err)
	}
	return nil
}

func (c *Client) ApplyBoostSingle(guildID string, slotID string) error {
	return c.ApplyBoost(guildID, []string{slotID})
}

func (c *Client) UnboostSlot(currentGuildID string, pgSubID string) error {
	_, _, err := c.Delete(
		fmt.Sprintf("/guilds/%s/premium/subscriptions/%s", currentGuildID, pgSubID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("UnboostSlot from guild %s sub %s: %w", currentGuildID, pgSubID, err)
	}
	return nil
}

func (c *Client) TransferBoost(slot *BoostSlot, targetGuildID string) error {
	if slot.PremiumGuildSubscription == nil {
		return fmt.Errorf("slot %s is not currently boosting any guild", slot.ID)
	}

	currentGuildID := slot.PremiumGuildSubscription.GuildID
	pgSubID := slot.PremiumGuildSubscription.ID

	if err := c.UnboostSlot(currentGuildID, pgSubID); err != nil {
		return fmt.Errorf("TransferBoost unboost step: %w", err)
	}

	if err := c.ApplyBoostSingle(targetGuildID, slot.ID); err != nil {
		return fmt.Errorf("TransferBoost reboost step: %w", err)
	}

	return nil
}
