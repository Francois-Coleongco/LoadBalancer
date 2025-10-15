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

	TotalConnections uint64

	alive bool
}

func (n *Node) Available() bool {
	return n.alive
}

type Servers struct {
	mu    sync.RWMutex
	Size  uint64
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
	n.alive = true

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

	s.Nodes[whole].alive = false // server is dead.

	s.Size--
}

func (s *Servers) TraverseMNodes(m uint64) uint64 {
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

func (s *Servers) GetServer() (*Node, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Size == 0 { // no servers in list
		return nil, fmt.Errorf("Server retrieved was nil")
	}

	for s.First.Available() == false {
		s.DeleteServer(s.First.URL, s.First.PORT)
	}

	ret := s.First

	s.First = s.First.next

	fmt.Printf("hitting server: %s on port %d\n", ret.URL, ret.PORT)

	return ret, nil
}

func (s *Servers) GetMean() (uint64, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()
	// this function might not be needed if i optimize the mean calculation on the fly during addition of connections (GetServer) and deletion of servers (DeleteServer)

	if s.Size == 0 {
		return 0, fmt.Errorf("couldn't get mean because there are no servers!")
	}

	all_connections := s.TraverseMNodes(0) // go for the whole length

	mean := all_connections / s.Size

	return mean, nil

}

// func (s *Servers) GetServerStandardDeviation(mean uint64, server_connection_count uint64) {
//
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
//
// 	diff := (server_connection_count - mean)
// 	(diff * diff) / (s.Size - 1)
// }
