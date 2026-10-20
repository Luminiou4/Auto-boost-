package core

import (
	"autoboost/discord"
)

type SlotWork struct {
	Account    *discord.AccountStatus
	Slot       *discord.BoostSlot
	NeedsXfer  bool
}

type BoostPlan struct {
	TargetGuildID  string
	Work           []SlotWork
	TotalScheduled int
	SkippedCooldown int
}

func BuildPlan(
	accounts []*discord.AccountStatus,
	guildID string,
	amount int,
	includeTransfer bool,
) *BoostPlan {
	plan := &BoostPlan{TargetGuildID: guildID}

	for _, acc := range accounts {
		if acc.Error != "" || !acc.HasNitro {
			continue
		}
		for _, slot := range acc.ReadySlots {
			plan.Work = append(plan.Work, SlotWork{
				Account:   acc,
				Slot:      slot,
				NeedsXfer: false,
			})
		}
	}

	if includeTransfer {
		for _, acc := range accounts {
			if acc.Error != "" || !acc.HasNitro {
				continue
			}
			for _, slot := range acc.BusySlots {

				if slot.PremiumGuildSubscription != nil &&
					slot.PremiumGuildSubscription.GuildID == guildID {
					continue
				}
				plan.Work = append(plan.Work, SlotWork{
					Account:   acc,
					Slot:      slot,
					NeedsXfer: true,
				})
			}
		}
	}

	for _, acc := range accounts {
		plan.SkippedCooldown += len(acc.CooldownSlots)
	}

	if amount > 0 && amount < len(plan.Work) {
		plan.Work = plan.Work[:amount]
	}

	plan.TotalScheduled = len(plan.Work)
	return plan
}

func AvailableCount(accounts []*discord.AccountStatus, includeTransfer bool) int {
	n := 0
	for _, acc := range accounts {
		if acc.Error != "" || !acc.HasNitro {
			continue
		}
		n += len(acc.ReadySlots)
		if includeTransfer {
			n += len(acc.BusySlots)
		}
	}
	return n
}
