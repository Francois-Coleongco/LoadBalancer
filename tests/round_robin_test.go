package tests

import (
	"fmt"
	"github.com/Francois-Coleongco/LoadBalancer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRobin(t *testing.T) {
	s := types.InitServers()
	s.Lives.AddToFront("69.69.69.69", 421)
	s.Lives.AddToFront("70.70.70.70", 422)
	s.Lives.AddToFront("71.71.71.71", 423)
	s.Lives.AddToFront("72.72.72.72", 424)
	s.Lives.AddToFront("73.73.73.73", 425)

	fmt.Println(s.Lives.Nodes)

	s.Lives.TraverseMNodes(10)

	s.DeleteServer("72.72.72.72", 424)
	s.DeleteServer("69.69.69.69", 421)

	fmt.Println("TRAVERSING AGAIN")

	s.Lives.TraverseMNodes(0)

	assert.Equal(t, s.Lives.Size, uint64(3), "should be equal Lives")
	assert.Equal(t, s.Deads.Size, uint64(2), "should be equal Deads")

}

func TestRobinServerBecameHealthy(t *testing.T) {
	// the following is sequential so it is safe to access s.Lives.Size
	s := types.InitServers()

	s.Lives.AddToFront("69.69.69.69", 421)
	s.Lives.AddToFront("70.70.70.70", 422)
	s.Lives.AddToFront("71.71.71.71", 423)
	s.Lives.AddToFront("72.72.72.72", 424)
	s.Lives.AddToFront("73.73.73.73", 425)

	fmt.Println(s.Lives.Nodes)

	s.DeleteServer("72.72.72.72", 424)
	assert.Equal(t, uint64(1), s.Deads.GetSize(), "deleted a single server, should have 1 in Deads")

	fmt.Println("TRAVERSING AGAIN")

	s.Lives.TraverseMNodes(0)

	fmt.Println("TRAVERSING Deadds")
	s.Lives.TraverseMNodes(0)

	s.Lives.AddToFront("72.72.72.72", 424)
	s.Deads.DeleteServerFromList("72.72.72.72:424")

	fmt.Println("TRAVERSING AGAIN")
	s.Lives.TraverseMNodes(0)

	assert.Equal(t, uint64(5), s.Lives.GetSize())
	assert.Equal(t, uint64(0), s.Deads.GetSize())

}

func TestRobinAddOne(t *testing.T) {
	s := types.InitServers()
	s.Lives.AddToFront("69.69.69.69", 421)

	fmt.Println(s.Lives.Nodes)

	s.Lives.TraverseMNodes(1)

	s.DeleteServer("69.69.69.69", 421)
	fmt.Println(s.Lives.Nodes)

	fmt.Println("TRAVERSING AGAIN")

	s.Lives.TraverseMNodes(1)

	assert.Equal(t, uint64(0), s.Lives.GetSize())

}
