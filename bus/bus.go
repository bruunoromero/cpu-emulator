package bus

type bus struct {
	channels map[string]chan []int8
}

// Instance is the interface of the bus type
type Instance interface {
	MakeChannel(string)
	SendTo(string, []int8)
	ReceiveFrom(string) []int8
}

// New returns a new instance of bus
func New() Instance {
	return &bus{
		channels: make(map[string]chan []int8),
	}
}

func (bus *bus) MakeChannel(channel string) {
	bus.channels[channel] = make(chan []int8)
}

func (bus *bus) SendTo(channel string, payload []int8) {
	go func() {
		bus.channels[channel] <- payload
	}()
}

func (bus *bus) ReceiveFrom(channel string) []int8 {
	return <-bus.channels[channel]
}
