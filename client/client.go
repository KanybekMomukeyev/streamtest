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
	"github.com/KanybekMomukeyev/streamtest/database"
)

var (
	serverAddr = "192.168.1.204:8080"
	//serverAddr = "localhost:10000"
	title = "constTitle"
)

func Chat(letters ...string) error {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	grpclog.Println("client started...")
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
			messageReceived, err := stream.Recv()
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

			grpclog.Printf("client -- server status: %s", messageReceived.Content)
			grpclog.Printf("client -- Title %s", messageReceived.Title)
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

var (
	msgc = make(chan string) // the message channel
)

func main() {
	flag.Parse()
	//database.SomeMethodBleve()
	fmt.Println("start the program")
	for {
		// start the app
		waitc := make(chan struct{}) // a wait lock

		// start the client thread
		go func() {
			for {
				msg := <-msgc // a message to send

				err := Chat(msg)

				if err != nil {
					fmt.Printf("send Err: %v", err)
				}
			}
		}()

		// start the input thread
		go func() {
			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				msgc <- text
			}
		}()

		<-waitc

		// finished in this round restart the app
		fmt.Println("restart the app")
	}
}
