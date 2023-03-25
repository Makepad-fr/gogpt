package gogpt

import (
	"math/rand"
	"time"
)

const baseURL = "https://chat.openai.com"


func randomTimeOut() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return float64(r.Intn(10000) + 1000)
}