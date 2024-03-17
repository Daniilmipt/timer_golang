package timer

type Timer struct {
	TimerOuter      
	id       uint32
	endTime  uint32
	interval uint32
	index    int
}