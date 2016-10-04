package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "github.com/KanybekMomukeyev/streamtest/protolocation"
)

var (
	port       = flag.Int("port", 10000, "The server port")
	grpcServer *grpc.Server
)

type chat_server struct {
}

func (*chat_server) Chat(stream pb.Chat_ChatServer) error {
	for {
		in, err := stream.Recv()

		// end of the streaming
		if err == io.EOF {
			grpclog.Println("server -- finished stream")
			return nil
		}

		if err != nil {
			grpclog.Printf("returned with error %v", err)
			return err
		}

		content := in.Content
		title := in.Title

		grpclog.Printf("server -- received message:\n%v: %v", title, content)
		revMsg := "message from server content received"

		stream.Send(&pb.Msg{Content: revMsg, Title:"server title"})
	}
}

func Shutdown() {
	grpclog.Println("shutdown server")
	grpcServer.Stop()
}

func main() {
	flag.Parse()

	fmt.Println("start the server")
	grpclog.Println("start server...")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatal("failed to listen: %v", err)
	}

	grpcServer = grpc.NewServer()
	pb.RegisterChatServer(grpcServer, new(chat_server))
	grpcServer.Serve(lis)

	grpclog.Println("server shutdown...")
}

//package main
//
//import (
//	"fmt"
//	"time"
//	"math/rand"
//	"bufio"
//	"os"
//	"flag"
//)
//
////func boring(msg string) {
////	for i := 0; ; i++ {
////		fmt.Println(msg, i)
////		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
////	}
////}
////
////func main() {
////	go boring("boring!")
////	fmt.Println("I'm listening.")
////	time.Sleep(2 * time.Second)
////	fmt.Println("You're boring; I'm leaving.")
////}
//
//func boring(msg string, c chan string) {
//	for i := 0; ; i++ {
//		c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any suitable value.
//		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
//	}
//}
//
//func boring2(msg string) <-chan string { // Returns receive-only channel of strings.
//	c := make(chan string)
//	go func() { // We launch the goroutine from inside the function.
//		for i := 0; ; i++ {
//			c <- fmt.Sprintf("%s %d", msg, i)
//			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
//		}
//	}()
//	return c // Return the channel to the caller.
//}
//
//func fanIn(input1, input2 <-chan string) <-chan string {
//	c := make(chan string)
//	go func() { for { c <- <-input1 } }()
//	go func() { for { c <- <-input2 } }()
//	return c
//}
//
//func fake_fibonacci(n int, c chan int) {
//	for i := 0; i < n; i++ {
//		c <- i
//	}
//	//defer close(c)
//}
//
//func fibonacci(c, quit chan int) {
//	x, y := 0, 1
//	for {
//		select {
//		case c <- x:
//			x, y = y, x+y
//		case <-quit:
//			fmt.Println("quit")
//			return
//		}
//	}
//}
//
//func call_fibonacci()  {
//		c := make(chan int)
//		quit := make(chan int)
//		go func() {
//			for i := 0; i < 10; i++ {
//				fmt.Println(<-c)
//			}
//			quit <- 0
//		}()
//		fibonacci(c, quit)
//}
//
////func main() {
////
////	call_fibonacci()
////
////	//c := make(chan int, 10)
////	//go fake_fibonacci(cap(c), c)
////	//
////	//for i := 0; i < cap(c); i++ {
////	//	fmt.Println(<-c)
////	//}
////
////	//for i := range c {
////	//	fmt.Println(i)
////	//}
////}
//var (
//	msgc = make(chan string) // the message channel
//)
//
//func main() {
//	flag.Parse()
//
//	fmt.Println("start the program")
//	for {
//		// start the app
//		waitc := make(chan struct{}) // a wait lock
//
//		// start the client thread
//		go func() {
//			for {
//				msg := <-msgc // a message to send
//				print(msg)
//
//				if msg == "break" {
//					break
//				}
//
//			}
//		}()
//
//		// start the input thread
//		go func() {
//			for {
//				reader := bufio.NewReader(os.Stdin)
//				text, _ := reader.ReadString('\n')
//				msgc <- text
//			}
//		}()
//
//		<-waitc
//
//		// finished in this round restart the app
//		fmt.Println("restart the app")
//	}
//}
//
//
//func someMethod()  {
//	ch := make(chan int, 2)
//	ch <- 1
//	ch <- 2
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//
//
//	cc := fanIn(boring2("Joe"), boring2("Ann"))
//	for i := 0; i < 10; i++ {
//		fmt.Println(<-cc)
//	}
//	fmt.Println("You're both boring; I'm leaving.")
//
//
//	joe := boring2("Joe")
//	ann := boring2("Ann")
//
//	for i := 0; i < 5; i++ {
//		fmt.Println(<-joe)
//		fmt.Println(<-ann)
//	}
//
//	fmt.Println("You're both boring; I'm leaving.")
//
//	// create channel
//	c := make(chan string)
//
//	// call goroutine
//	go boring("boring!", c)
//
//	fmt.Printf("1_You say: %q\n", <-c) // Receive expression is just a value.
//	fmt.Printf("2_You say: %q\n", <-c) // Receive expression is just a value.
//	fmt.Printf("3_You say: %q\n", <-c) // Receive expression is just a value.
//	fmt.Printf("4_You say: %q\n", <-c) // Receive expression is just a value.
//	fmt.Printf("5_You say: %q\n", <-c) // Receive expression is just a value.
//	fmt.Printf("6_You say: %q\n", <-c) // Receive expression is just a value.
//
//	//// receive from goroutine
//	//for i := 0; i < 15; i++ {
//	//	fmt.Printf("You say: %q\n", <-c) // Receive expression is just a value.
//	//}
//
//	time.Sleep(2 * time.Second)
//
//	fmt.Println("You're boring; I'm leaving.")
//
//}