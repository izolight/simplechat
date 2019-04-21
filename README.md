# simpleChat

simpleChat is a basic Chatserver and Client implementation in go and grpc

## Features

- multiple users
- nickname (not changeable/unique)
- relaying messages to all except yourself
- timestamps
- message for leave/join

## Todo

- name changing
- unique names
- channels
- permissions
- authentication
- better input handling

## How to use

Start a server ```go run server/server.go``` and connect with a client ```go run client/client.go```.