package pubsub

type PubSub struct {
	topics map[string][]chan []byte
}

func New() *PubSub {
	return &PubSub{
		topics: make(map[string][]chan []byte),
	}
}

func (p *PubSub) Subscribe(topic string) chan []byte {
	ch := make(chan []byte)
	_, ok := p.topics[topic]
	if ok {
		p.topics[topic] = append(p.topics[topic], ch)
	} else {
		p.topics[topic] = []chan []byte{ch}
	}
	return ch
}

func (p *PubSub) Unsubscribe(ch chan []byte) {

	for topic, channels := range p.topics {
		for i, c := range channels {
			if c == ch {
				close(ch)
				channels = append(channels[:i], channels[i+1:]...)
				if len(channels) == 0 {
					delete(p.topics, topic)
				}
			}
		}
	}
}

func (p *PubSub) Send(topic string, msg []byte) {
	go func() {
		t, ok := p.topics[topic]
		if ok {
			for _, ch := range t {
				ch <- msg
			}
		}
	}()
}
