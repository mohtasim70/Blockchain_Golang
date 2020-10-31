package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	a2 "github.com/mohtasim70/assignment02IBC"
)

//var chainHead *a2.Block
var Quorum int

type Client struct {
	conn             net.Conn
	listeningAddress string
}

func handleConnection(conn net.Conn, addchan chan Client) {
	newClient := Client{
		conn: conn,
	}
	var listeningAddress string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&listeningAddress)
	if err != nil {
		//handle error
	}
	Quorum--
	newClient.listeningAddress = listeningAddress
	addchan <- newClient

}
func handlePeer(conn net.Conn) {

	buf := make([]byte, 50)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		fmt.Println("Error in handPeer")
	}
	fmt.Println("Recieved in handle: ", string(buf))

}

//The function below launches the server, uses different second argument
//It then starts a routine for each connection request received
//2 parts
//The function below launches and initializes the chain and the server
//It then starts a routine for each connection request received
//The listening address of each node and their conn info is then stored
//it is important not to sequentually do things in StartListening routine and
//rather use channel for communication between routines
func StartListening(listeningAddress string, node string) {
	if node == "satoshi" {
		ln, err := net.Listen("tcp", listeningAddress)
		if err != nil {
			log.Fatal(err)
		}
		clientsSlice := make([]Client, Quorum)
		addchan := make(chan Client)
		var chainHead *a2.Block
		chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn, addchan)
			clientsSlice = append(clientsSlice, <-addchan)
			chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
		}
	} else {
		ln, err := net.Listen("tcp", listeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			go handlePeer(conn)

		}

	}
}

//this should block satoshi till the quorum is complete
//Hint: we can read from a channel to block a routine unless some other writes

//sends each peer, the address of another one to connect to in a ring based topology
//it also sends the current chain to each peer
func SendChainandConnInfo() {

}

//For Satoshiiiiii above below are for nodes

//Satoshi is there waiting for our address, it stores it somehow
func WriteString(conn net.Conn, myListeningAddress string) {
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(myListeningAddress)
	if err != nil {
		log.Println("In Write String: ", err)
	}
}

//once the satoshi unblocks on Quorum completion it sends peer to connect to
func ReadString(conn net.Conn) string {
	var client string
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&client)
	if err != nil {
		log.Println(err)
	}
	return client
}

//Then satoshi sends the chain with x+1 blocks
func ReceiveChain(conn net.Conn) *a2.Block {
	var block *a2.Block
	gobEncoder := gob.NewDecoder(conn)
	err := gobEncoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	return block
}

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
	go StartListening(myListeningAddress, "others")

	log.Println("Sending my listening address to Satoshi")
	//Satoshi is there waiting for our address, it stores it somehow
	WriteString(conn, myListeningAddress)

	//once the satoshi unblocks on Quorum completion it sends peer to connect to
	log.Println("receiving peer to connect to ... ")
	receivedString := ReadString(conn)
	log.Println("Recieved Client: ", receivedString)
	//
	// //Then satoshi sends the chain with x+1 blocks
	log.Println("receiving Chain")
	chainHead := ReceiveChain(conn)
	log.Println(a2.CalculateBalance("Satoshi", chainHead))

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
