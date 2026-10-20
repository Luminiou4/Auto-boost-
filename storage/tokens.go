package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadTokens(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("tokens file not found: %s\nCreate it with one Discord token per line.", filePath)
		}
		return nil, fmt.Errorf("open tokens file: %w", err)
	}
	defer f.Close()

	seen := make(map[string]bool)
	var tokens []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !ValidateTokenFormat(line) {
			continue
		}

		if !seen[line] {
			seen[line] = true
			tokens = append(tokens, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read tokens file: %w", err)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no valid tokens found in %s", filePath)
	}

	return tokens, nil
}

func ValidateTokenFormat(token string) bool {

	if len(token) < 50 {
		return false
	}

	if strings.ContainsAny(token, " \t\r\n") {
		return false
	}

	if strings.Count(token, ".") < 1 {
		return false
	}
	return true
}
