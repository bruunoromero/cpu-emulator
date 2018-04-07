package bus

type bus struct {
	cpuChannel chan int
	ramChannel chan []int
}

// Instance is the interface of the bus type
type Instance interface {
	SendToCPU(int)
	SendToRAM([]int)
	ReceiveAtCPU() int
	ReceiveAtRAM() []int
}

// New returns a new instance of bus
func New() Instance {
	return &bus{
		cpuChannel: make(chan int),
		ramChannel: make(chan []int),
	}
}

func (bus *bus) SendToCPU(payload int) {
	go func() {
		bus.cpuChannel <- payload
	}()
}

func (bus *bus) SendToRAM(payload []int) {
	go func() {
		bus.ramChannel <- payload
	}()
}

func (bus *bus) ReceiveAtCPU() int {
	return <-bus.cpuChannel
}

func (bus *bus) ReceiveAtRAM() []int {
	return <-bus.ramChannel
}
