package delay

import (
	"math/rand"
	"time"
)

var (
	Randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func PretendHeavyOperation() {
	sleepyTime := Randomizer.Intn(500) + 250 // Anywhere between 250 and 750 ms
	time.Sleep(time.Duration(sleepyTime) * time.Millisecond)
}
