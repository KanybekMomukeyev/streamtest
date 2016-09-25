package main

import (
	"io"
	pb "github.com/KanybekMomukeyev/testingGRPC/protolocation"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc"
	"flag"
	"fmt"
	"bufio"
	"os"
)

var (
	serverAddr *string
	title      string
)

func Chat(letters ...string) error {
	// get connection for chat
	conn := connect(serverAddr)
	defer conn.Close()
	client := pb.NewChatClient(conn)

	stream, err := client.Chat(context.Background())
	if err != nil {
		grpclog.Println("%v.Chat(_) = _, %v", client, err) // better logging
		return err
	}

	// receive msg
	waitc := make(chan struct{})
	var recevieErr error

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done
				close(waitc)
				return
			}
			if err != nil {
				grpclog.Printf("Failed to receive a msg : %v", err) // need better logging
				recevieErr = err
				return
			}
			grpclog.Printf("client -- server status: %s", in.Content)
		}
	}()

	if recevieErr != nil {
		return recevieErr
	}

	// send msg
	for _, str := range letters {
		grpclog.Printf("client -- send msg: %v", str)
		if err := stream.Send(&pb.Msg{Content: str, Title: title}); err != nil {
			grpclog.Printf("%v.Send(%v) = %v", stream, str, err) // need better logging
			return err
		}
	}

	// close send
	stream.CloseSend()
	<-waitc
	return nil
}

func InitChatClient(t string, srvAddr *string) {
	title = t
	serverAddr = srvAddr
}

func connect(srvAddr *string) *grpc.ClientConn {
	conn, err := grpc.Dial(*srvAddr, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	grpclog.Println("client started...")

	return conn
}

var (
	msgc = make(chan string) // the message channel
	//serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	myTitle    = flag.String("title", "", "The name show to your friend")
)

// an input from command line
func input() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		msgc <- text
	}
}

func main() {
	flag.Parse()
	fmt.Println("start the program")

	for {
		// start the app
		waitc := make(chan struct{}) // a wait lock

		// start the client thread
		go func() {
			for {
				msg := <-msgc // a message to send
				InitChatClient(*myTitle, serverAddr)
				err := Chat(msg)
				if err != nil {
					// restart the client
					fmt.Printf("send Err: %v", err)
				}
			}
		}()

		// start the input thread
		go input()

		<-waitc
		// finished in this round restart the app
		fmt.Println("restart the app")
	}
}
