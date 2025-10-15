package tests

import (
	"fmt"
	"github.com/Francois-Coleongco/LoadBalancer/types"
	"testing"
)

func TestRobin(t *testing.T) {
	s := types.InitServers()
	s.AddToFront("69.69.69.69", 421)
	s.AddToFront("70.70.70.70", 422)
	s.AddToFront("71.71.71.71", 423)
	s.AddToFront("72.72.72.72", 424)
	s.AddToFront("73.73.73.73", 425)

	fmt.Println(s.Nodes)

	s.TraverseMNodes(10)

	s.DeleteServer("72.72.72.72", 424)
	s.DeleteServer("69.69.69.69", 421)

	fmt.Println("TRAVERSING AGAIN")

	s.TraverseMNodes(0)

}

func TestRobinRemoveAll(t *testing.T) {
	s := types.InitServers()
	s.AddToFront("69.69.69.69", 421)
	s.AddToFront("70.70.70.70", 422)
	s.AddToFront("71.71.71.71", 423)
	s.AddToFront("72.72.72.72", 424)
	s.AddToFront("73.73.73.73", 425)

	fmt.Println(s.Nodes)

	s.TraverseMNodes(10)

	s.DeleteServer("72.72.72.72", 424)
	s.DeleteServer("69.69.69.69", 421)

	fmt.Println("TRAVERSING AGAIN")

	s.TraverseMNodes(0)

}

func TestRobinAddOne(t *testing.T) {
	s := types.InitServers()
	s.AddToFront("69.69.69.69", 421)

	fmt.Println(s.Nodes)

	s.TraverseMNodes(1)

	s.DeleteServer("69.69.69.69", 421)
	fmt.Println(s.Nodes)

	fmt.Println("TRAVERSING AGAIN")

	s.TraverseMNodes(1)

}
