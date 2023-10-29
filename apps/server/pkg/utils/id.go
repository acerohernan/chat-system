package utils

import "github.com/lithammer/shortuuid/v4"

const guidSize = 12

const (
	UserPrefix = "US_"
)

func NewGuid(prefix string) string {
	return prefix + shortuuid.New()[:guidSize]
}
