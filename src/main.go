package main

import (
	"fmt"
	"github.com/Daniilmipt/timer_golang/timer"
	"time"
)

// main is the entry point of the program.
// It creates a channel to signal when to close the program,
// creates a timer manager, starts the scheduler, waits a bit,
// adds three timers to the timer manager, waits 20 seconds,
// removes one of the timers, and waits for the program to close.
func main() {

	// This channel is used to signal when to close the program.
	closeChan := make(chan struct{})

	// Create a new timer manager.
	timerManager := timer.NewTimerManager()

	// Start the scheduler to run the timerManager.
	go scheduler(timerManager)

	// Wait a second to let the scheduler start.
	time.Sleep(time.Second * 1)

	// Get the current time in Unix time format.
	now := uint32(time.Now().Unix())

	// Add three timers to the timer manager:
	// 1. A() to run immediately.
	// 2. B() to run in 5 seconds.
	// 3. C() to run in 10 seconds.
	timerManager.AddTimer(&timer.TimerCallback{CallBack: A}, now, 0)
	timerId := timerManager.AddTimer(&timer.TimerCallback{CallBack: B}, now, 5)
	timerManager.AddTimer(&timer.TimerCallback{CallBack: C}, now, 10)

	// Wait 20 seconds to let all the timers run.
	time.Sleep(time.Second * 20)

	// Remove one of the timers from the timer manager by its id.
	timerManager.RemoveTimer(timerId)

	// Wait for the program to close.
	<-closeChan
}

func scheduler(timerManager *timer.Manager) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			timerManager.RunTimer()
		}
	}
}

func A() {
	now := uint32(time.Now().Unix())
	fmt.Printf("%v => aaa\n", now)
}

func B() {
	now := uint32(time.Now().Unix())
	fmt.Printf("%v => bbb\n", now)
}

func C() {
	now := uint32(time.Now().Unix())
	fmt.Printf("%v => ccc\n", now)
}
