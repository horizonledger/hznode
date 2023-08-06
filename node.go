package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	log "github.com/sirupsen/logrus"
)

var lg = log.New()

type Subscriber interface {
	Notify(string)
}

type Publisher struct {
	subs   map[Subscriber]struct{}
	mu     sync.RWMutex
	closed bool
}

func NewPublisher() *Publisher {
	return &Publisher{
		subs: make(map[Subscriber]struct{}),
	}
}

func (p *Publisher) Subscribe(s Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.closed {
		p.subs[s] = struct{}{}
	}
}

func (p *Publisher) Unsubscribe(s Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subs, s)
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true
}

func (p *Publisher) Publish(msg string) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.closed {
		for s := range p.subs {
			s.Notify(msg)
		}
	}
}

type Vertex struct {
	Connection net.Conn
}

func (v *Vertex) Notify(msg string) {
	err := wsutil.WriteServerText(v.Connection, []byte(msg))
	if err != nil {
		lg.Println("Write error: ", err)
	}
}

type Node struct {
	Publisher *Publisher
	Handler   *http.ServeMux
}

func (n *Node) handleRead(v *Vertex) {
	defer v.Connection.Close()

	for {
		msg, _, err := wsutil.ReadClientData(v.Connection)
		if err != nil {
			// handle error
			if err == io.EOF {
				lg.Println("Connection closed")
			} else {
				lg.Println("Read error: ", err)
			}
			//n.Publisher.Unsubscribe(v)
			break
		}

		lg.Info("msg: ", string(msg))

		if string(msg) == "ping" {
			err := wsutil.WriteServerText(v.Connection, []byte("pong"))
			if err != nil {
				// handle error
				lg.Println("Write error: ", err)
				break
			}
		}
	}
}

func (n *Node) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
		lg.Println(err)
		return
	}

	vertex := &Vertex{Connection: conn}
	n.Publisher.Subscribe(vertex)

	go n.handleRead(vertex)

	// go func() {
	// 	defer conn.Close()

	// 	for {
	// 		_, _, err := wsutil.ReadClientData(conn)
	// 		if err != nil {
	// 			// handle error
	// 			if err == io.EOF {
	// 				lg.Println("Connection closed")
	// 			} else {
	// 				lg.Println("Read error: ", err)
	// 			}
	// 			n.Publisher.Unsubscribe(vertex)
	// 			break
	// 		}
	// 	}
	// }()
}

func NewNode() *Node {
	n := &Node{
		Publisher: NewPublisher(),
		Handler:   http.NewServeMux(),
	}
	n.routes()
	return n
}

func (n *Node) routes() {
	fs := http.FileServer(http.Dir("./clientstatic"))
	n.Handler.Handle("/", fs)
	n.Handler.HandleFunc("/ws", n.handleWS)
}

func main() {
	node := NewNode()

	go func() {
		for !node.Publisher.closed {
			time.Sleep(time.Second)
			msg := fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC3339))
			node.Publisher.Publish(msg)
		}
	}()

	lg.Println("Listening on localhost:8000")
	http.ListenAndServe("localhost:8000", node.Handler)
}
