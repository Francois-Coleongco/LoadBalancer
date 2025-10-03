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
		s.First = n
	} else {
		old_prev := s.First.prev
		s.First.prev = n
		n.next = s.First
		old_prev.next = n
		n.prev = old_prev
	}

	whole := url + strconv.Itoa(int(port))

	s.Nodes[whole] = n

	s.Size++
}

func (s *Servers) DeleteNode(url string, port uint16) {
	whole := url + strconv.Itoa(int(port))

	value, ok := s.Nodes[whole]

	fmt.Println("removing server: ", url, port)

	if !ok {
		fmt.Println("could not remove server", value.url, value.port)
		return
	}

	value.prev.next = value.next
	value.next.prev = value.prev

	delete(s.Nodes, whole)

	if value == s.First {
		s.First = s.First.next
	}
}

func (s *Servers) TraverseMNodes(m uint32) {
	if m == 0 { // 0 passes means you can't even iterate one server
		m = s.Size
	}

	curr_node := s.First

	for i := uint32(0); i < m; i++ {
		println(curr_node.url, curr_node.port)
		curr_node = curr_node.prev
	}

}
