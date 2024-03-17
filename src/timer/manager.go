package timer

import (
	"container/heap"
	"time"
)

type Manager struct {
	id        uint32
	tq        TimerQueue
	execTimer []interface{}
}

func NewTimerManager() *Manager {
	return &Manager{
		tq: make(TimerQueue, 0, 1024),
	}
}


// AddTimer adds a timer to the timer manager.
// The timer manager will call the TimerOuter's TimeOut method when the timer
// expires. If the interval is not 0, the timer manager will add the timer
// back into the queue after it expires.
//
// timerOuter: The TimerOuter to call TimeOut on when the timer expires.
//
// endTime: The time in Unix time format when the timer should expire.
//
// interval: How often to reschedule the timer in seconds. If 0, the timer
// will not be rescheduled.
//
// Returns: The timer id of the timer that was added. This can be used to
// remove the timer later.
func (this *Manager) AddTimer(timerOuter TimerOuter, endTime uint32, interval uint32) uint32 {

	this.id++

	timer := &Timer{
		id:         this.id,
		TimerOuter: timerOuter,
		endTime:    endTime,
		interval:   interval,
	}

	heap.Push(&this.tq, timer)

	return this.id
}

func (this *Manager) RemoveTimer(timerId uint32) {
	for _, timer := range this.tq {
		if timer.id == timerId {
			heap.Remove(&this.tq, timer.index)
			return
		}
	}
}


// RunTimer runs any timers that have expired.
//
// If there are no timers in the queue, this function does nothing.
//
// Otherwise, this function pulls the timer with the soonest expiration time
// from the queue and adds it to the list of timers to call TimeOut on.
//
// If the timer's interval is not 0, the function reschedules the timer
// by adding it back into the queue with its updated expiration time.
//
// Finally, the function calls TimeOut on all the timers that were pulled
// from the queue.
func (this *Manager) RunTimer() {

	// If there are no timers in the queue, there is nothing to do.
	if this.tq.Len() <= 0 {
		return
	}

	// Loop through the timers in the queue and pull out any timers that
	// have expired.
	for this.tq.Len() > 0 {

		// Get the timer with the soonest expiration time.
		tmp := this.tq[0]

		// If the current time is before the expiration time of the soonest
		// expiring timer, there are no more expired timers in the queue.
		if uint32(time.Now().Unix()) < tmp.endTime {
			break
		}

		// Remove the timer from the queue and add it to the list of timers
		// to call TimeOut on.
		timer := heap.Pop(&this.tq).(*Timer)
		this.execTimer = append(this.execTimer, timer)

		// If the timer's interval is not 0, add the timer back into the
		// queue with its updated expiration time.
		if timer.interval > 0 {
			timer.endTime += timer.interval
			heap.Push(&this.tq, timer)
		}
	}

	// If there are any timers to call TimeOut on, call TimeOut on them.
	if len(this.execTimer) > 0 {
		for _, timer := range this.execTimer {
			timer.(TimerOuter).TimeOut()
		}
	}

	// Reset the list of timers to call TimeOut on so it is ready for the
	// next time RunTimer is called.
	this.execTimer = this.execTimer[:0]
}

