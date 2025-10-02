package types

import (
	"fmt"
	"strconv"
)

type Node struct {
	prev *Node
	next *Node

	url  string
	port uint16

	alive bool
}

func (n *Node) Available() bool {
	return n.alive
}

type Servers struct {
	Size  uint32
	Nodes map[string]*Node // map string to node so easy lookup for deletion
	First *Node            // could use first to add new nodes
}

func InitServers() *Servers {
	s := new(Servers)

	s.Size = 0
	s.Nodes = make(map[string]*Node)
	s.First = nil

	return s
}

func (s *Servers) AddToFront(url string, port uint16) {

	n := new(Node)
	n.url = url
	n.port = port

	if s.Size == 0 {
		n.prev = n
		n.next = n
	}

	n.prev = s.First.prev
	n.next = s.First
	s.First.prev = n

	whole := url + strconv.Itoa(int(port))

	s.Nodes[whole] = n

	s.Size++
}

func (s *Servers) DeleteNode(url string, port uint16) {
	whole := url + strconv.Itoa(int(port))

	value, ok := s.Nodes[whole]

	if !ok {
		fmt.Println("could not remove server", value.url, value.port)
		return
	}

	value.prev.next = value.next
	value.next.prev = value.prev
}
