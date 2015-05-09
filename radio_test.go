package radio

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRadio(t *testing.T) {
	NewRadio()
}

func TestNoise(t *testing.T) {
	radio := NewRadio()

	assert.True(t, radio.IsLive())

	err1 := radio.Broadcast(Noise{})
	assert.Nil(t, err1)

	assert.False(t, radio.IsLive())

	err2 := radio.Broadcast("hi there")
	assert.NotNil(t, err2)
}

func TestWakeUp(t *testing.T) {
	radio := NewRadio()

	_, id1 := radio.Listen()

	radio.Broadcast(Noise{})
	assert.False(t, radio.IsLive())

	ch2, _ := radio.Listen()

	err1 := radio.WakeUp()
	assert.Nil(t, err1)
	assert.True(t, radio.IsLive())

	ch3, _ := radio.Listen()

	err2 := radio.WakeUp()
	assert.NotNil(t, err2)

	err3 := radio.Call(id1, "too late")
	assert.NotNil(t, err3)

	radio.Broadcast("nice")

	assert.Equal(t, "nice", <-ch2)
	assert.Equal(t, "nice", <-ch3)
}

func TestBroadcastAndListen(t *testing.T) {
	radio := NewRadio()

	ch1, _ := radio.Listen()
	ch2, _ := radio.Listen()

	radio.Broadcast("hello")

	assert.Equal(t, "hello", <-ch1)
	assert.Equal(t, "hello", <-ch2)

	radio.Broadcast("nice")

	assert.Equal(t, "nice", <-ch1)
	assert.Equal(t, "nice", <-ch2)
}

func TestStop(t *testing.T) {
	radio := NewRadio()

	ch1, id1 := radio.Listen()
	ch2, _ := radio.Listen()

	radio.Stop(id1)

	radio.Broadcast("hello")

	assert.Equal(t, "hello", <-ch2)

	channelClosed := false
	select {
	case x, ok := <-ch1:
		if ok {
			fmt.Printf("Value %d was read.\n", x)
		} else {
			channelClosed = true
		}
	default:
		fmt.Println("No value ready, moving on.")
	}

	assert.True(t, channelClosed)

	assert.NotNil(t, radio.Stop(1337))
}

func TestCall(t *testing.T) {
	radio := NewRadio()

	ch1, id1 := radio.Listen()
	ch2, id2 := radio.Listen()

	err1 := radio.Call(id1, "hello there")
	err2 := radio.Call(id2, "wow omg")
	err3 := radio.Call(1337, "not happening dear")

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotNil(t, err3)

	assert.Equal(t, "hello there", <-ch1)
	assert.Equal(t, "wow omg", <-ch2)
}

func TestAlienate(t *testing.T) {
	radio := NewRadio()

	_, id1 := radio.Listen()
	_, id2 := radio.Listen()

	assert.Equal(t, 2, radio.Count())

	radio.Alienate()

	assert.Equal(t, 0, radio.Count())

	err1 := radio.Call(id1, "hello there")
	err2 := radio.Call(id2, "wow omg")

	assert.NotNil(t, err1)
	assert.NotNil(t, err2)

	// Now, check that broadcasting still works afterwards
	ch3, id3 := radio.Listen()

	assert.Equal(t, 1, radio.Count())

	radio.Broadcast("ok, this should work")
	radio.Call(id3, "this should too")

	assert.Equal(t, "ok, this should work", <-ch3)
	assert.Equal(t, "this should too", <-ch3)
}
