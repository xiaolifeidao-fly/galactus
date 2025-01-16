package common

import (
	"encoding/hex"
	"github.com/lifei6671/gorand"
	"math/rand"
	"time"
)

func GetUUId() string {
	return gorand.NewUUID4().String()
}

func RandomString(n int) string {
	b := make([]byte, (n+1)/2)
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	if _, err := src.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)[:n]
}
