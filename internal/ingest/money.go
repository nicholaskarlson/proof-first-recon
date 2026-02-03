package ingest

import (
	"fmt"
	"strings"
)

func ParseAmountCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	sign := int64(1)
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}
	if s == "" {
		return 0, fmt.Errorf("no digits")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("too many dots")
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" {
		intPart = "0"
	}

	var dollars int64
	for _, ch := range intPart {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("invalid digit %q", ch)
		}
		dollars = dollars*10 + int64(ch-'0')
	}

	if len(fracPart) > 2 {
		return 0, fmt.Errorf("more than 2 decimal places")
	}
	var cents int64
	if len(fracPart) == 1 {
		fracPart = fracPart + "0"
	}
	if len(fracPart) == 2 {
		for _, ch := range fracPart {
			if ch < '0' || ch > '9' {
				return 0, fmt.Errorf("invalid digit %q", ch)
			}
			cents = cents*10 + int64(ch-'0')
		}
	}

	return sign * (dollars*100 + cents), nil
}

func FormatCents(c int64) string {
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	dollars := c / 100
	cents := c % 100
	return fmt.Sprintf("%s%d.%02d", sign, dollars, cents)
}
