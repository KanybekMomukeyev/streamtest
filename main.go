package main

import (
	"fmt"
	"time"
	"math/rand"
)

//func boring(msg string) {
//	for i := 0; ; i++ {
//		fmt.Println(msg, i)
//		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
//	}
//}
//
//func main() {
//	go boring("boring!")
//	fmt.Println("I'm listening.")
//	time.Sleep(2 * time.Second)
//	fmt.Println("You're boring; I'm leaving.")
//}

func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any suitable value.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func boring2(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() { for { c <- <-input1 } }()
	go func() { for { c <- <-input2 } }()
	return c
}

func main() {

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)


	cc := fanIn(boring2("Joe"), boring2("Ann"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-cc)
	}
	fmt.Println("You're both boring; I'm leaving.")


	joe := boring2("Joe")
	ann := boring2("Ann")

	for i := 0; i < 5; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}

	fmt.Println("You're both boring; I'm leaving.")




	// create channel
	c := make(chan string)

	// call goroutine
	go boring("boring!", c)

	fmt.Printf("1_You say: %q\n", <-c) // Receive expression is just a value.
	fmt.Printf("2_You say: %q\n", <-c) // Receive expression is just a value.
	fmt.Printf("3_You say: %q\n", <-c) // Receive expression is just a value.
	fmt.Printf("4_You say: %q\n", <-c) // Receive expression is just a value.
	fmt.Printf("5_You say: %q\n", <-c) // Receive expression is just a value.
	fmt.Printf("6_You say: %q\n", <-c) // Receive expression is just a value.

	//// receive from goroutine
	//for i := 0; i < 15; i++ {
	//	fmt.Printf("You say: %q\n", <-c) // Receive expression is just a value.
	//}

	time.Sleep(2 * time.Second)

	fmt.Println("You're boring; I'm leaving.")
}
