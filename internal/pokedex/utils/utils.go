package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomDurations(count int) []time.Duration {
	var durations []time.Duration

	for i := 0; i < count; i++ {
		// Generate a random duration between 200ms and 2000ms
		duration := time.Duration(GenerateRandomNum(200,1500)) * time.Millisecond
		durations = append(durations, duration)
	}

	return durations
}

func DummyLoading(){
	durations := generateRandomDurations(3)
    for _, duration := range durations {
		fmt.Println("...")
        time.Sleep(duration)
    }
}

func GenerateRandomNum(min,max int)int{
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	
	return rng.Intn(max-min)+min
}