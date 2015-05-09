# [radio](https://github.com/etcinit/radio) [![GoDoc](https://godoc.org/github.com/etcinit/radio?status.svg)](https://godoc.org/github.com/etcinit/radio)

Broadcast messages to multiple Go channels

[![wercker status](https://app.wercker.com/status/a73f3aa6cee48c69c9737b1417927354/m "wercker status")](https://app.wercker.com/project/bykey/a73f3aa6cee48c69c9737b1417927354)

## Summary

- Useful for broadcasting messages across your app
- You could build a multi-threaded IM/logger server with it
(Emphasis on the _could_).
- It does not block when broadcasting or if one of listeners hasn't read from
a channel yet.

## Quick Example

```go
package main

import (
	"fmt"

	"github.com/etcinit/radio"
)

func main() {
	done := make(chan bool)

    // First, we begin by creating a Radio object. Think of it as your own
    // radio station.
	r1 := radio.NewRadio()

    // Now we add two listeners. We get both a channel and identifier for them.
	ch1, id1 := r1.Listen()
	ch2, id2 := r1.Listen()

    // We make the listener do things on a separate goroutine.
	go func() {
		for message := range ch1 {
			if content, ok := message.(string); ok {
				fmt.Println("CH1: ", content)
			}

			done <- true
		}
	}()

    // Same for the second one.
	go func() {
		for message := range ch2 {
			if content, ok := message.(string); ok {
				fmt.Println("CH2: ", content)
			}

			done <- true
		}
	}()

    // We broadcast our first message.
	r1.Broadcast("hello world")

    // Just for this example, we make sure they actually got it.
	<-done
	<-done
    // By this point, both listeners should have written to stdout.

    // Listener 1 got bored, so we remove it.
	r1.Stop(id1)

    // Broadcast a second message.
	r1.Broadcast("you are listening to radio 1 news")

    // Just for this example, we make sure it actually got it.
	<-done

    // Listener 2 also got bored, so we also remove it.
	r1.Stop(id2)

    // Nothing happens if no one is listening. :(
	r1.Broadcast("sadly, no one listens")
}
```

and the expected output is:

```
CH1:  hello world
CH2:  hello world
CH2:  you are listening to radio 1 news
```
