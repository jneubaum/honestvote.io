package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	coredb "github.com/jneubaum/honestvote.io/core/core-database/src"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var nodes = make(map[int]bool)

var Peers []coredb.Peer

var mongoDB *mongo.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Loading ENV Failed")
	}

	mongoDB = coredb.MongoConnect()

	ignore, _ := strconv.Atoi(os.Getenv("PORT"))
	nodes[ignore] = true

	go listenConn()

	for {
		for port := 7000; port <= 7001; port++ {
			if !nodes[port] {
				fmt.Println("Checking...")
				sPort := strconv.Itoa(port)
				conn, _ := net.Dial("tcp", "127.0.0.1:"+sPort)
				if conn != nil {
					fmt.Println("Dial Successful!")
					tmpPeer := coredb.Peer{
						IPAddress: "127.0.0.1",
						Port:      port,
						Socket:    conn,
					}
					Peers = append(Peers, tmpPeer)
					nodes[port] = true

					conn.Write([]byte("connect " + strconv.Itoa(ignore)))
					go handleConn(conn)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func listenConn() {
	portString := ":" + os.Getenv("PORT")
	listen, err := net.Listen("tcp", portString)

	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte

	for {
		length, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		if string(buf[0:7]) == "connect" {
			port, err := strconv.Atoi(string(buf[8:length]))

			if err == nil {
				nodes[port] = true
				tmpPeer := coredb.Peer{
					IPAddress: "127.0.0.1",
					Port:      port,
					Socket:    conn,
				}
				Peers = append(Peers, tmpPeer)
			}
		} else if string(buf[0:8]) == "get data" {
			moveDocuments()
		}
	}
}

//Send the data to the full/peer node
func moveDocuments() {
	MongoData := coredb.GatherMongoData(mongoDB, bson.M{})
	buffer := new(bytes.Buffer)
	tmpArray := MongoData
	js := json.NewEncoder(buffer)
	err := js.Encode(tmpArray)
	if err == nil {
		for _, socket := range Peers {
			fmt.Println("Sending documents.")
			socket.Socket.Write(append([]byte("recieve data "), buffer.Bytes()...))
		}
	}
}
