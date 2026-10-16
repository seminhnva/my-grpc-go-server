package application

import (
	"fmt"
	"math/rand"
	"time"
)

type ResiliencyService struct {
}

func (f *ResiliencyService) GenerateResiliency(minDelaySecond int32, maxDelaySercond int32, statusCodes []uint32) (string, uint32) {
	delay := rand.Intn(int(maxDelaySercond-minDelaySecond)) + int(minDelaySecond)
	delaySecond := time.Duration(delay) * time.Second
	time.Sleep(delaySecond)

	idx := rand.Intn(len(statusCodes))
	str := fmt.Sprintf("The time now is %v, execution delay for %v seconds", time.Now().Format("15:04:05.000"), delay)

	return str, statusCodes[idx]
}
