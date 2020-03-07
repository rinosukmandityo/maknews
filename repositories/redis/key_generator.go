package redis

import (
	"fmt"
)

func generateKey(code string) string {
	return fmt.Sprintf("news<>%s", code)
}

func generateKeyOffsetLimit(offset, limit int) string {
	return fmt.Sprintf("news<>%d<>%d", offset, limit)
}
