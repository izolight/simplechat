//go:generate protoc -I ../chat --go_out=plugins=grpc:../chat ../chat/chat.proto

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/izolight/simplechat/chat"

	"google.golang.org/grpc"
)

const (
	port = ":12345"
)

type chatServer struct{
	clientMessages []chan *pb.ChatMessage
}

func (c *chatServer)sendToAll(me chan *pb.ChatMessage, msg *pb.ChatMessage) {
	for _, mc := range c.clientMessages {
		if mc != me {
			mc <- msg
		}
	}
}

func newServerMessage(text string) *pb.ChatMessage {
	msg := &pb.ChatMessage{}
	msg.User = "Server"
	msg.Timestamp = time.Now().Unix()
	msg.Text = text
	return msg
}

func (c *chatServer) SendMessage(stream pb.Chat_SendMessageServer) error {
	messages := make(chan *pb.ChatMessage)
	c.clientMessages = append(c.clientMessages, messages)
	errors := make(chan error)
	newUser := true
	user := ""
	// receive all messages and send to all servers
	go func () {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				msg := newServerMessage(fmt.Sprintf("%s has left the server", user))
				go c.sendToAll(messages, msg)
				errors <- nil
				break
			}
			if err != nil {
				errors <- err
				break
			}
			if newUser {
				msg := newServerMessage(fmt.Sprintf("%s joined the server", in.User))
				go c.sendToAll(messages, msg)
				user = in.User
				newUser = false
			}
			in.Timestamp = time.Now().Unix()
			go c.sendToAll(messages, in)
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
