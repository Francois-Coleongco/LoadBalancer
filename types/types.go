package types

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Node struct {
	prev *Node
	next *Node

	URL  string
	PORT uint16

	TotalConnections uint64
}

type ServersList struct {
	mu    sync.RWMutex
	Nodes map[string]*Node // map string to node so easy lookup for deletion (funeral time)
	First *Node
	Size  uint64
}

type Servers struct {
	Lives ServersList
	Deads ServersList
}

func InitServers() *Servers {
	// will block and doesn't need locking
	s := new(Servers)
	s.Lives.Size = 0
	s.Lives.Nodes = make(map[string]*Node)
	s.Lives.First = nil

	s.Deads.Size = 0
	s.Deads.Nodes = make(map[string]*Node)
	s.Deads.First = nil

	return s
}

func (s *ServersList) DeleteServerFromList(whole string) {
	value, ok := s.Nodes[whole]

	log.Println("removing server: ", whole)

	if !ok {
		log.Println("could not remove server", whole)
		return
	}

	log.Println("got past could not remove server, no deadlock?")

	if value == s.First {
		s.First = s.First.next
	}

	if s.Size > 1 { // removes the server from the round robin schedule
		value.prev.next = value.next
		value.next.prev = value.prev
	} else {
		s.First = nil
	}

	s.Size--

}

func (s *ServersList) GetSize() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Size
}

func (s *ServersList) AddToFront(url string, port uint16) {

	s.mu.Lock()
	defer s.mu.Unlock()

	n := new(Node)
	n.URL = url
	n.PORT = port

	if s.Size == 0 {
		n.prev = n
		n.next = n
		s.First = n
	} else {
		old_prev := s.First.prev
		s.First.prev = n
		n.next = s.First
		old_prev.next = n
		n.prev = old_prev
	}

	whole := url + ":" + strconv.Itoa(int(port))

	s.Nodes[whole] = n

	log.Println("added server: ", whole)
	s.Size++
}

func (s *Servers) DeleteServer(url string, port uint16) {
	{
		s.Lives.mu.Lock()
		defer s.Lives.mu.Unlock()

		whole := url + ":" + strconv.Itoa(int(port))

		s.Lives.DeleteServerFromList(whole)

	}

	// link up s.Nodes[whole] to the DeadNodes

	log.Println("starting in DeadNodes")

	s.Deads.AddToFront(url, port) // safe to use url and port from the args because server was in Lives for certain, and was deleted as per assumption from the above return

}

func (s *ServersList) TraverseMNodes(m uint64) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m == 0 { // 0 passes means you can't even iterate one server
		m = s.Size
	}

	if s.Size == 0 {
		return 0
	}

	curr_node := s.First

	var sum uint64 = 0

	for i := uint64(0); i < m; i++ {
		println(curr_node.URL, curr_node.PORT)
		sum += curr_node.TotalConnections
		curr_node = curr_node.prev
	}

	return sum
}

func (s *ServersList) GetServer() (*Node, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Size == 0 { // no servers in list
		return nil, fmt.Errorf("Server retrieved was nil")
	}

	ret := s.First

	s.First = s.First.next

	fmt.Printf("hitting server: %s on port %d\n", ret.URL, ret.PORT)

	return ret, nil
}

func (s *Servers) ListenAndAddBack() {
	for {
		size := s.Deads.GetSize()
		for i := uint64(0); i < size; i++ {
			server, err := s.Deads.GetServer()
			if err != nil {
				log.Println("GetServer failed in ListenAndAddBack")
				continue
			}

			whole := server.URL + ":" + strconv.Itoa(int(server.PORT))
			_, err = http.Get(whole)

			if err == nil {
				// server is alive
				s.Lives.AddToFront(server.URL, server.PORT)
				s.Deads.DeleteServerFromList(whole)
			}

		}

	}

}
