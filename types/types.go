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
	size  uint32
	nodes map[string]*Node // map string to node so easy lookup for deletion
	first *Node            // could use first to add new nodes
}

func (s *Servers) AddToFront(url string, port uint16) {

	n := new(Node)
	n.url = url
	n.port = port

	if s.size == 0 {
		n.prev = n
		n.next = n
	}

	n.prev = s.first.prev
	n.next = s.first
	s.first.prev = n

	whole := url + strconv.Itoa(int(port))

	s.nodes[whole] = n

	s.size++
}

func (s *Servers) DeleteNode(url string, port uint16) {
	whole := url + strconv.Itoa(int(port))

	value, ok := s.nodes[whole]

	if !ok {
		fmt.Println("could not remove server", value.url, value.port)
		return
	}

	value.prev.next = value.next
	value.next.prev = value.prev
}
