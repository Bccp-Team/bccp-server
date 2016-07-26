package runners

import (
	"encoding/gob"
	"flag"
	"log"
	"net"
	"sync"
)

var runnerService string
var runnerToken string

var runnerMaps map[uint]*clientInfo

func WaitRunners() {
	flag.StringVar(&runnerService, "runner-service", "127.0.0.1:4243", "the runner service")
	flag.StringVar(&runnerToken, "runner-token", "bccp", "the runner token")

	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", runnerService)

	if err != nil {
		log.Panic(err)
	}

	//FIXME: TLS, toussa toussa
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf(err.Error())
			continue
		}

		go handleClient(conn)
	}
}

type clientInfo struct {
	uid        uint
	currentRun int
	conn       net.Conn
	mut        sync.Mutex
	encoder    *gob.Encoder
	decoder    *gob.Decoder
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	var connection SubscribeRequest
	err := decoder.Decode(&connection)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	if connection.Token != runnerToken {
		log.Printf("bad token receive: %v", connection.Token)
		return
	}

	//FIXME: generate Token and ClientUid"command.hh"
	var uid uint = 0
	answer := SubscribeAnswer{ClientUID: uid}

	err = encoder.Encode(&answer)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	client := &clientInfo{uid: uid, conn: conn, encoder: encoder, decoder: decoder}

	runnerMaps[uid] = client

	for {
		var clientReq ClientRequest
		err = decoder.Decode(&clientReq)

		if err != nil {
			log.Printf(err.Error())
			return
		}

		switch clientReq.Kind {
		case Ack:
		case Finish:
		case Logs:
		case Error:
		default:
		}
	}
}
