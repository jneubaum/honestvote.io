package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jneubaum/honestvote/tests/logger"

	"github.com/jneubaum/honestvote/core/core-crypto/crypto"
	"github.com/jneubaum/honestvote/core/core-database/database"
	"github.com/jneubaum/honestvote/core/core-discovery/discovery"
	"github.com/jneubaum/honestvote/core/core-http/http"
	"github.com/jneubaum/honestvote/core/core-p2p/p2p"
	"github.com/joho/godotenv"
)

//defaults
var TCP_PORT string = "7000"  //tcp PORT for peer to peer routes
var UDP_PORT string = "7001"  //udp PORT for node discovery
var HTTP_PORT string = "7002" //tcp PORT for light nodes to http routes

var ROLE string = "producer" //options producer || full || registry
var DATABASE_HOST string = "127.0.0.1"
var COLLECTION_PREFIX string = ""
var REGISTRY_IP string
var REGISTRY_PORT string = "7002"
var REGISTRY bool = false // is producer registry node or not
var INSTITUTION_NAME string
var PUBLIC_KEY string
var PRIVATE_KEY string
var LOGGING_MODE string = "All" // Levels of Debugging All | Debug | Info
var HOSTNAME string = "127.0.0.1"
var EMAIL_ADDRESS string
var EMAIL_PASSWORD string
var REGISTRATION_TYPE string // DEFAULT_WHITELIST | EXTERNAL_WHITELIST | ALL

// external white list election parameters
var WHITELIST_DATABASE_DRIVER string
var WHITELIST_DATABASE_USER string
var WHITELIST_DATABASE_PASSWORD string
var WHITELIST_DATABASE_HOST string
var WHITELIST_DATABASE_PORT string
var WHITELIST_DATABASE_NAME string
var WHITELIST_TABLE_NAME string
var ELIGIBLE_VOTER_FIELD string

//this file will be responsible for deploying the app

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Loading ENV Failed")
	}

	logger.Mode = LOGGING_MODE
	PRIVATE_KEY, PUBLIC_KEY = crypto.GenerateKeyPair()

	// environmental variables override defaults
	if os.Getenv("TCP_PORT") != "" {
		TCP_PORT = os.Getenv("TCP_PORT")
	}
	if os.Getenv("UDP_PORT") != "" {
		UDP_PORT = os.Getenv("UDP_PORT")
	}
	if os.Getenv("HTTP_PORT") != "" {
		HTTP_PORT = os.Getenv("HTTP_PORT")
	}
	if os.Getenv("ROLE") != "" {
		ROLE = os.Getenv("ROLE")
	}
	if os.Getenv("DATABASE_HOST") != "" {
		ROLE = os.Getenv("DATABASE_HOST")
	}
	if os.Getenv("COLLECTION_PREFIX") != "" {
		COLLECTION_PREFIX = os.Getenv("COLLECTION_PREFIX")
	}
	if os.Getenv("REGISTRY_IP") != "" {
		REGISTRY_IP = os.Getenv("REGISTRY_IP")
	}
	if os.Getenv("REGISTRY_PORT") != "" {
		REGISTRY_PORT = os.Getenv("REGISTRY_PORT")
	}
	if os.Getenv("REGISTRY") != "" {
		REGISTRY, _ = strconv.ParseBool(os.Getenv("REGISTRY"))
	}
	if os.Getenv("PRIVATE_KEY") != "" {
		PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
	}
	if os.Getenv("PUBLIC_KEY") != "" {
		PUBLIC_KEY = os.Getenv("PUBLIC_KEY")
	}
	if os.Getenv("HOSTNAME") != "" {
		HOSTNAME = os.Getenv("HOSTNAME")
	}
	if os.Getenv("INSTITUTION_NAME") != "" {
		INSTITUTION_NAME = os.Getenv("INSTITUTION_NAME")
	}
	if os.Getenv("EMAIL_ADDRESS") != "" {
		EMAIL_ADDRESS = os.Getenv("EMAIL_ADDRESS")
	}
	if os.Getenv("EMAIL_PASSWORD") != "" {
		EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	}
	if os.Getenv("REGISTRATION_TYPE") != "" {
		REGISTRATION_TYPE = os.Getenv("REGISTRATION_TYPE")
	}
	if os.Getenv("WHITELIST_DATABASE_DRIVER") != "" {
		WHITELIST_DATABASE_DRIVER = os.Getenv("WHITELIST_DATABASE_DRIVER")
	}
	if os.Getenv("WHITELIST_DATABASE_USER") != "" {
		WHITELIST_DATABASE_USER = os.Getenv("WHITELIST_DATABASE_USER")
	}
	if os.Getenv("WHITELIST_DATABASE_PASSWORD") != "" {
		WHITELIST_DATABASE_PASSWORD = os.Getenv("WHITELIST_DATABASE_PASSWORD")
	}
	if os.Getenv("WHITELIST_DATABASE_HOST") != "" {
		WHITELIST_DATABASE_HOST = os.Getenv("WHITELIST_DATABASE_HOST")
	}
	if os.Getenv("WHITELIST_DATABASE_PORT") != "" {
		WHITELIST_DATABASE_PORT = os.Getenv("WHITELIST_DATABASE_PORT")
	}
	if os.Getenv("WHITELIST_DATABASE_NAME") != "" {
		WHITELIST_DATABASE_NAME = os.Getenv("WHITELIST_DATABASE_NAME")
	}
	if os.Getenv("WHITELIST_TABLE_NAME") != "" {
		WHITELIST_TABLE_NAME = os.Getenv("WHITELIST_TABLE_NAME")
	}
	if os.Getenv("ELIGIBLE_VOTER_FIELD") != "" {
		ELIGIBLE_VOTER_FIELD = os.Getenv("ELIGIBLE_VOTER_FIELD")
	}

	//this domain is the default host to resolve traffic
	if REGISTRY_IP == "" {
		registry_ip, err := net.LookupIP("registry.honestvote.io")
		if err != nil {
			fmt.Println("Unknown host")
		} else {
			REGISTRY_IP = registry_ip[0].String()
		}
	}

	// accept optional flags that override environmental variables
	for index, element := range os.Args {
		switch element {
		case "--tcp": //Set the default port for node tcp PORT
			TCP_PORT = os.Args[index+1]
		case "--udp":
			UDP_PORT = os.Args[index+1]
		case "--http": //Set the default port for http PORT
			HTTP_PORT = os.Args[index+1]
		case "--role": //Set the role of the node options producer || FULL || REGISTRY
			ROLE = os.Args[index+1]
		case "--database-host":
			DATABASE_HOST = os.Args[index+1]
		case "--collection-prefix": //Collection prefix (useful for starting up multiple nodes with same database)
			COLLECTION_PREFIX = os.Args[index+1]
		case "--registry-host": //Sets the registry node
			REGISTRY_IP = os.Args[index+1]
		case "--registry-port": //Sets the registry node port
			REGISTRY_PORT = os.Args[index+1]
		case "--registry":
			REGISTRY, _ = strconv.ParseBool(os.Args[index+1])
		case "--private-key": //Sets the private key
			PRIVATE_KEY = os.Args[index+1]
		case "--public-key": //Sets the public key
			PUBLIC_KEY = os.Args[index+1]
		case "--hostname": //sets the public ip address
			HOSTNAME = os.Args[index+1]
		case "--institution-name": //sets the institutions name
			INSTITUTION_NAME = os.Args[index+1]
		case "--email-address": //sets the public ip address
			EMAIL_ADDRESS = os.Args[index+1]
		case "--email-password": //sets the institutions name
			EMAIL_PASSWORD = os.Args[index+1]
		case "--registration-type":
			REGISTRATION_TYPE = os.Args[index+1]
		case "--whitelist-database-driver":
			WHITELIST_DATABASE_DRIVER = os.Args[index+1]
		case "--whitelist-database-user":
			WHITELIST_DATABASE_USER = os.Args[index+1]
		case "--whitelist-database-password":
			WHITELIST_DATABASE_PASSWORD = os.Args[index+1]
		case "--whitelist-database-host":
			WHITELIST_DATABASE_HOST = os.Args[index+1]
		case "--whitelist-database-port":
			WHITELIST_DATABASE_PORT = os.Args[index+1]
		case "--whitelist-database-name":
			WHITELIST_DATABASE_NAME = os.Args[index+1]
		case "--eligible-voter-field":
			ELIGIBLE_VOTER_FIELD = os.Args[index+1]
		}
	}

	database.CollectionPrefix = COLLECTION_PREFIX
	database.MongoDB = database.MongoConnect(DATABASE_HOST) // Connect to data store

	port, _ := strconv.Atoi(TCP_PORT)

	p2p.Self = database.Node{
		IPAddress:   HOSTNAME,
		Port:        port,
		Role:        ROLE,
		PublicKey:   PUBLIC_KEY,
		Institution: INSTITUTION_NAME,
	}

	p2p.REGISTRATION_TYPE = REGISTRATION_TYPE
	p2p.Whitelist = database.WhiteListElectionSettings{
		DatabaseDriver:     WHITELIST_DATABASE_DRIVER,
		DatabaseUser:       WHITELIST_DATABASE_USER,
		DatabasePassword:   WHITELIST_DATABASE_PASSWORD,
		DatabaseHost:       WHITELIST_DATABASE_HOST,
		DatabasePort:       WHITELIST_DATABASE_PORT,
		DatabaseName:       WHITELIST_DATABASE_NAME,
		TableName:          WHITELIST_TABLE_NAME,
		EligibleVoterField: ELIGIBLE_VOTER_FIELD,
	}
	p2p.PrivateKey = PRIVATE_KEY
	p2p.PublicKey = PUBLIC_KEY
	p2p.Email_Address = EMAIL_ADDRESS
	p2p.Email_Password = EMAIL_PASSWORD

	if !database.DoesNodeExist(p2p.Self) && REGISTRY {
		database.AddNode(p2p.Self)
	}

	if REGISTRY {
		p2p.ConsensusNodes = 1
	} else {
		discovery.FetchLatestPeers(REGISTRY_IP, REGISTRY_PORT, TCP_PORT)
	}

	// udp PORT that sends connected producer to incoming nodes
	// if ROLE == "registry" {
	// 	registry.ListenConnections(UDP_PORT)
	// }

	logger.Println("main.go", "main", "Collection Prefix: "+COLLECTION_PREFIX)

	go http.CreateServer(HTTP_PORT, ROLE)

	// accept incoming connections and handle p2p
	p2p.HTTP_Port = HTTP_PORT
	p2p.ListenConn(TCP_PORT, ROLE)

}
