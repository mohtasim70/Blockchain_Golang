package main

import (
	"fmt"
	"log"
	"net"
	"os"

	a2 "github.com/ehteshamz/assignment02IBC"
	a3 "github.com/ehteshamz/assignment03IBC"
)

func main() {
	//it is better to check for arguments length and throw error
	satoshiAddress := os.Args[1]
	myListeningAddress := os.Args[2]

	conn, err := net.Dial("tcp", satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}
	//The function below launches the server, uses different second argument
	//It then starts a routine for each connection request received
	go a3.StartListening(myListeningAddress, "others")

	log.Println("Sending my listening address to Satoshi")
	//Satoshi is there waiting for our address, it stores it somehow
	a3.WriteString(conn, myListeningAddress)

	//once the satoshi unblocks on Quorum completion it sends peer to connect to
	log.Println("receiving peer to connect to ... ")
	receivedString := a3.ReadString(conn)
	log.Println(receivedString)

	//Then satoshi sends the chain with x+1 blocks
	log.Println("receiving Chain")
	chainHead := a3.ReceiveChain(conn)
	log.Println(a2.CalculateBalance("Satoshi", &chainHead))

	//Each node then connects to the other peer info received from satoshi
	//The topology eventually becomes a ring topology
	//Each node then both writes the hello world to the connected peer
	//and also receives the one from another peer

	log.Println("connecting to the other peer ... ")
	peerConn, err := net.Dial("tcp", receivedString)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Fprintf(peerConn, "Hello From %v to %v\n", myListeningAddress, receivedString)
	}
	select {}
}
