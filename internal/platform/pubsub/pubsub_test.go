package pubsub_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
)

func TestPubSub(t *testing.T) {
	topicA := "topicA"
	topicB := "topicB"
	msg := "Hello"

	ps := pubsub.New()
	chA := ps.Subscribe(topicA)
	chB := ps.Subscribe(topicB)

	timeout := make(chan bool, 1)
	result := make(chan error, 1)

	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	aDone := false
	bDone := false

	go func() {
		for {
			select {
			case m := <-chA:
				if string(m) != msg+" A" {
					result <- fmt.Errorf("Expected [%s],  got [%s]", msg+" A", string(m))
				} else {
					aDone = true
				}

			case m := <-chB:
				if string(m) != msg+" B" {
					result <- fmt.Errorf("Expected [%s],  got [%s]", msg+" B", string(m))
				} else {
					bDone = true
				}

			case <-timeout:
				result <- errors.New("Timeout waiting for message")
			}

			if aDone && bDone {
				result <- nil
			}
		}
	}()

	ps.Send(topicA, []byte(msg+" A"))
	ps.Send(topicB, []byte(msg+" B"))

	err := <-result
	if err != nil {
		t.Error(err)
	}
}
