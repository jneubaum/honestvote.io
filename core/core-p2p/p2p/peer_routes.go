package p2p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/jneubaum/honestvote/core/core-database/database"
	"github.com/jneubaum/honestvote/tests/logger"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	for {
		d := json.NewDecoder(conn)
		var write Message
		d.Decode(&write)

		switch write.Message {
		case "connect":
			logger.Println("peer_routes.go", "HandleConn()", "Recieved Connect Message")

			var node database.Node
			json.Unmarshal(write.Data, &node)

			AcceptConnectMessage(node, conn)
		case "get id":
			var node database.Node
			json.Unmarshal(write.Data, &node)
			node.IPAddress = conn.RemoteAddr().String()[0:9]
			if !database.DoesNodeExist(node) {
				database.AddNode(node)
			}
		case "recieve data":
			buffer := bytes.NewBuffer(write.Data)
			DecodeData(buffer)
		case "get data":
			database.MoveDocuments(Nodes, database.DatabaseName, database.CollectionPrefix+database.ElectionHistory)
		case "vote":
			vote := write.Vote
			ReceiveVote(vote)
		case "verify":
			block := new(database.Block)
			json.Unmarshal(write.Data, block)
			logger.Println("peer_routes.go", "HandleConn()", "Verifying")
			VerifyBlock(*block, conn)
		case "sign":
			block := new(database.Block)
			err := json.Unmarshal(write.Data, &block)
			if err == nil {
				ReceiveResponses(block, write.Signature)
			} else {
				fmt.Println(err)
			}
		case "update":
			block := new(database.Block)
			json.Unmarshal(write.Data, block)
			if database.UpdateBlockchain(database.MongoDB, *block) {
				PrevHash = block.Hash
				PrevIndex = block.Index
				fmt.Println(string(PrevIndex) + " " + PrevHash)
			}
		}

	}
}
