package core

import (
	"autoboost/config"
	"autoboost/discord"
	"autoboost/utils"
	"sync"
)

func CheckAllAccounts(tokens []string, cfg *config.Config, proxyFn func() string) []*discord.AccountStatus {
	results := make([]*discord.AccountStatus, len(tokens))

	sem := make(chan struct{}, cfg.MaxConcurrent)
	var wg sync.WaitGroup

	for i, token := range tokens {
		wg.Add(1)
		go func(idx int, tok string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

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

			client := discord.New(tok, opts...)
			status := client.CheckAccount(tok)
			results[idx] = status

			utils.RandomDelay(cfg.CheckDelay, cfg.CheckDelay+500*1000*1000)
		}(i, token)
	}

	wg.Wait()
	return results
}
