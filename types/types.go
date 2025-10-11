package types

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

type Node struct {
	prev *Node
	next *Node

	URL  string
	PORT uint16

	TotalConnections uint32

	alive bool
}

func (n *Node) Available() bool {
	return n.alive
}

type Servers struct {
	mu    sync.RWMutex
	Size  uint32
	Nodes map[string]*Node // map string to node so easy lookup for deletion
	First *Node            // could use first to add new nodes
}

func InitServers() *Servers {
	// will block and doesn't need locking
	s := new(Servers)

	s.Size = 0
	s.Nodes = make(map[string]*Node)
	s.First = nil

	return s
}

func (s *Servers) AddToFront(url string, port uint16) {

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

func (s *Servers) DeleteNode(url string, port uint16) {

	s.mu.Lock()
	defer s.mu.Unlock()

	whole := url + ":" + strconv.Itoa(int(port))

	value, ok := s.Nodes[whole]

	log.Println("removing server: ", url, port, whole)

	if !ok {
		log.Println("could not remove server", url, port)
		return
	}

	if value == s.First {
		s.First = s.First.next
	}

	if s.Size > 1 {
		value.prev.next = value.next
		value.next.prev = value.prev
	} else {
		s.First = nil
	}

	delete(s.Nodes, whole)

	s.Size--
}

func (s *Servers) TraverseMNodes(m uint32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m == 0 { // 0 passes means you can't even iterate one server
		m = s.Size
	}

	if s.Size == 0 {
		return
	}

	curr_node := s.First

	for i := uint32(0); i < m; i++ {
		println(curr_node.URL, curr_node.PORT)
		curr_node = curr_node.prev
	}

}

func (s *Servers) GetServer() (*Node, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Size == 0 { // no servers in list
		return nil, fmt.Errorf("Server retrieved was nil")
	}

	for s.First.Available() == false {
		s.DeleteNode(s.First.URL, s.First.PORT)
	}

	s.First = s.First.next

	return s.First, nil
}
