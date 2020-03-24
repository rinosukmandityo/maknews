package redis

import (
	"fmt"
	"strconv"

	m "github.com/rinosukmandityo/maknews/models"
)

func generateKey(code string) string {
	return fmt.Sprintf("news<>%s", code)
}

func generateKeyScore(data m.News) (string, float64) {
	return generateKey(strconv.Itoa(data.ID)), getScore(data)
}

func getScore(data m.News) float64 {
	return float64(data.Created.UnixNano())
}
