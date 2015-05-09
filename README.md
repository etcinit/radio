# [radio](https://github.com/etcinit/radio)

Broadcast messages to multiple Go channels

## Quick Example

```go
package main

import (
	"fmt"

	"github.com/etcinit/radio"
)

func main() {
	done := make(chan bool)

	r1 := radio.NewRadio()

	ch1, id1 := r1.Listen()
	ch2, id2 := r1.Listen()

	go func() {
		for message := range ch1 {
			if content, ok := message.(string); ok {
				fmt.Println("CH1: ", content)
			}

			done <- true
		}
	}()

	go func() {
		for message := range ch2 {
			if content, ok := message.(string); ok {
				fmt.Println("CH2: ", content)
			}

			done <- true
		}
	}()

	r1.Broadcast("hello world")

	<-done
	<-done

	r1.Stop(id1)

	r1.Broadcast("you are listening to radio 1 news")

	<-done

	r1.Stop(id2)

	r1.Broadcast("sadly, no one listens")
}
```

and the expected output is:

```
CH1:  hello world
CH2:  hello world
CH2:  you are listening to radio 1 news
```
