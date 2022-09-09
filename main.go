package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	natsserver "github.com/nats-io/nats-server/v2/test"
)

var s *server.Server

func init() {
	// Create server
	s, _ = natsserver.RunServerWithConfig("srv.conf")
}

func main() {
	// Create our customers.
	c1 := natsConnect("u1", "p1")
	defer c1.Close()
	c2 := natsConnect("u2", "p2")
	defer c1.Close()

	// Have each customer subscribe for root requests inbound to them.
	c1.Subscribe("cosmonic.requests", reqHandlerCust1)
	c2.Subscribe("cosmonic.requests", reqHandlerCust2)

	// Create our root user who will send the requests.
	nc := natsConnect("root", "s3cr3t!", nats.CustomInboxPrefix("$R"))
	defer nc.Close()

	// Make a system request to cust1
	sysRequest(nc, "cust1")
	// Make a system request to cust2
	sysRequest(nc, "cust2")
}

// This sends a request and gets a response.
func sysRequest(nc *nats.Conn, cust string) {
	reqSubj := fmt.Sprintf("req.%s", cust)
	r, err := nc.Request(reqSubj, nil, time.Second)
	if err != nil {
		log.Fatalf("Error sending request: %v")
	}
	log.Printf("Got a response for %q: %q", cust, r.Data)
}

// Callback for customers to receive requests.
func reqHandlerCust1(m *nats.Msg) {
	m.Respond([]byte("Cust1 - OK"))
}
func reqHandlerCust2(m *nats.Msg) {
	m.Respond([]byte("Cust2 - OK"))
}

// Helper to connect.
func natsConnect(user, pass string, opts ...nats.Option) *nats.Conn {
	opts = append(opts, nats.UserInfo(user, pass))
	nc, err := nats.Connect(s.ClientURL(), opts...)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	return nc
}
