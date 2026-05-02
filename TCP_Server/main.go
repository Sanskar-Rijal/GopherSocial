package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp",":8080");
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		// Reader and writer interfaces very important 
		conn, err := ln.Accept()
		if err != nil {
			//handle Error
		}
		//Go routine to handle multiple connections 
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
		defer conn.Close()
		//creating a. new reader from the connection 
		reader := bufio.NewReader(conn)
		//Read the command line from the connection 
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(conn, "Error reading command: %v\n",err)
			return
		}

		//Trim the newline Character and split the line into command and resouce 
		parts := strings.SplitN(strings.TrimSpace(line)," ",2)
		if len(parts) != 2 {
			fmt.Fprintf(conn, "Invalid command format. Expected format: Command:Resourse\n")
			return 
		}

		command := parts[0]
		resource := parts[1]
		log.Printf("Received Command as: %s %s \n", command, resource)

		//Handle the command 
		switch command {
		case "GET":
			handleGet(conn, resource)
		default:
			fmt.Fprintf(conn, "Unknown command: %s\n", command)
		}
	}

	func handleGet(conn net.Conn, resource string){
		fmt.Fprintf(conn, "GET command received for resource : %s\n", resource);
	}