package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/jneubaum/honestvote/core/core-consensus/consensus"
	"github.com/jneubaum/honestvote/core/core-crypto/crypto"
	"github.com/jneubaum/honestvote/core/core-database/database"
	"github.com/jneubaum/honestvote/tests/logger"
)

//Send a block out to be verified by other peers
func ProposeBlock(block database.Block, peers []net.Conn) {
	j, err := json.Marshal(block)

	write := new(Message)
	write.Message = "verify"
	write.Data = j

	jWrite, err := json.Marshal(write)

	if err == nil {
		fmt.Println(len(peers))
		for _, peer := range peers {
			peer.Write(jWrite)
		}
	}
}

//Decide if the block sent is valid
func VerifyBlock(block database.Block, conn net.Conn) {
	var valid bool

	if consensus.VerifyHash(PrevIndex, PrevHash, block) {
		valid = true
	} else {
		valid = false
	}

	j, err := json.Marshal(block)

	write := new(Message)
	write.Message = "sign"
	write.Data = []byte(strconv.FormatBool(valid))
	write.Signature = make(map[string]string)
	write.Signature[PublicKey], err = crypto.Sign(j, PrivateKey)

	if err != nil {
		write.Signature[PublicKey] = "Signature"
	}

	jWrite, err := json.Marshal(write)

	if err == nil {
		logger.Println("peer_routes.go", "HandleConn()", "Sending response")
		conn.Write(jWrite)
	}
}

//Go through all responses from other peers and see the result
func CheckResponses(size int) {
	counter := size

	checkBlock, err := json.Marshal(ProposedBlock)

	if err != nil {
		return
	}

	/*
		Iterate through nested map that holds boolean as first arg and
		a map[string]string as its second
	*/
	for b, v1 := range SignatureMap {
		for k, v2 := range v1 {
			valid, err := crypto.Verify(checkBlock, k, v2)
			if valid && err == nil && b {
				ProposedBlock.Signatures[k] = v2
			} else {
				counter--
			}
		}
	}

	if size == counter {
		j, err := json.Marshal(ProposedBlock)

		write := new(Message)
		write.Message = "update"
		write.Data = j

		jWrite, err := json.Marshal(write)

		if err == nil {
			if database.UpdateBlockchain(database.MongoDB, ProposedBlock) {
				PrevHash = ProposedBlock.Hash
				PrevIndex = ProposedBlock.Index
				fmt.Println(string(PrevIndex) + " " + PrevHash)
			}

			for _, node := range Nodes {
				node.Write(jWrite)
			}
		}
	} else {
		logger.Println("data_exchange.go", "CheckResponses", "Someone is a bad actor or this block is wrong.")
	}
}

func CreateElection(election database.Election) {

	database.UpdateElections(database.MongoDB, election)

	j, err := json.Marshal(election)

	write := new(Message)
	write.Message = "make election"
	write.Data = j

	jWrite, err := json.Marshal(write)

	if err == nil {
		for _, node := range Nodes {
			node.Write(jWrite)
		}
	}
}
