package main

import (
	"log"
	"os"
	"strconv"

	a3 "github.com/ehteshamz/assignment03IBC"
)

func main() {
	//it is better to check for arguments length and throw error
	satoshiAddress := os.Args[1]
	//should have used a setter and made the variable private
	a3.Quorum, _ = strconv.Atoi(os.Args[2])

	//The function below launches and initializes the chain and the server
	//It then starts a routine for each connection request received
	//The listening address of each node and their conn info is then stored
	//it is important not to sequentually do things in StartListening routine and
	//rather use channel for communication between routines
	go a3.StartListening(satoshiAddress, "satoshi")

	//this should block satoshi till the quorum is complete
	//Hint: we can read from a channel to block a routine unless some other writes
	a3.WaitForQuorum()
	log.Println("Quorum complete, satoshi unblocked!!")

	//sends each peer, the address of another one to connect to in a ring based topology
	//it also sends the current chain to each peer
	a3.SendChainandConnInfo()

	//blocking the main routine, we have others working
	//it is better to use this or (even reading from a channel never being written to)
	//infinite for loop is not recommended
	select {}
}
