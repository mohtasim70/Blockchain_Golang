package assignment03IBC

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"

	a2 "github.com/mohtasim70/assignment02IBC"
)

var chainHead *a2.Block
var wg sync.WaitGroup // 1

var Quorum int

type Client struct {
	conn             net.Conn
	ListeningAddress string
}

//var ListeningAddress []string
var clientsSlice = make([]Client, Quorum)
var rwchan = make(chan string)

func handleConnection(conn net.Conn, addchan chan Client) {
	newClient := Client{
		conn: conn,
	}
	var ListeningAddress string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&ListeningAddress)
	if err != nil {
		//handle error
	}
	Quorum--

	newClient.ListeningAddress = ListeningAddress
	fmt.Println("inHandle: ", newClient.ListeningAddress)
	addchan <- newClient
	//WaitForQuorum()

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
func StartListening(ListeningAddress string, node string) {
	if node == "satoshi" {
		ln, err := net.Listen("tcp", ListeningAddress)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Faital")
		}

		addchan := make(chan Client)

		chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err, "Yooooo")
				continue
			}
			if Quorum >= 0 {
				go handleConnection(conn, addchan)
				clientsSlice = append(clientsSlice, <-addchan)
				chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
			}

		}
	} else {
		ln, err := net.Listen("tcp", ListeningAddress)
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
func WaitForQuorum() {
	for {
		if Quorum == 0 {
			break
		}
	}

}

//sends each peer, the address of another one to connect to in a ring based topology
//it also sends the current chain to each peer
func SendChainandConnInfo() {
	for i := 0; i < len(clientsSlice); i++ {
		if i == len(clientsSlice)-1 {
			gobEncoder := gob.NewEncoder(clientsSlice[i].conn)
			err := gobEncoder.Encode(clientsSlice[0].ListeningAddress)
			if err != nil {
				//	log.Println(err)
			}
		} else {
			gobEncoder := gob.NewEncoder(clientsSlice[i].conn)
			err := gobEncoder.Encode(clientsSlice[i+1].ListeningAddress)
			if err != nil {
				//	log.Println(err)
			}
		}
		gobEncoder := gob.NewEncoder(clientsSlice[i].conn)
		err := gobEncoder.Encode(chainHead)
		if err != nil {
			//	log.Println(err)
		}
	}

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
