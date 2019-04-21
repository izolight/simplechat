package main

import (
	"io"
	"log"
	"net"
	"time"

	pb "simpleChat/chat"

	"google.golang.org/grpc"
)

const (
	port = ":12345"
)

type chatServer struct{
	clientMessages []chan *pb.ChatMessage
}

func (c *chatServer) SendMessage(stream pb.Chat_SendMessageServer) error {
	messages := make(chan *pb.ChatMessage)
	c.clientMessages = append(c.clientMessages, messages)
	errors := make(chan error)

	// receive all messages and send to all servers
	go func () {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errors <- err
				break
			}
			in.Timestamp = time.Now().Unix()
			go func() {
				for _, mc := range c.clientMessages {
					if mc != messages {
						mc <- in
					}
				}
			}()
		}
	}()

	for {
		select {
		case msg := <- messages:
			if err := stream.Send(msg); err != nil {
				return err
			}
		case err := <-errors:
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &chatServer{make([]chan *pb.ChatMessage,0)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
