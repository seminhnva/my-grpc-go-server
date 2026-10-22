package application

import (
	"fmt"
	"math/rand"
	"time"
)

type ResiliencyService struct {
}

func selectStatusCode(statusCodes []uint32, intn func(int) int) uint32 {
	return statusCodes[intn(len(statusCodes))]
}

func (f *ResiliencyService) GenerateResiliency(minDelaySecond int32, maxDelaySercond int32, statusCodes []uint32) (string, uint32) {
	delay := rand.Intn(int(maxDelaySercond-minDelaySecond)) + int(minDelaySecond)
	delaySecond := time.Duration(delay) + time.Second
	time.Sleep(delaySecond)

	// Repeating a status code in statusCodes increases its selection probability.
	statusCode := selectStatusCode(statusCodes, rand.Intn)
	str := fmt.Sprintf("The time now is %v, execution delay for %v seconds, statusCode: %v", time.Now().Format("15:04:05.000"), delay, statusCode)

	return str, statusCode
}
