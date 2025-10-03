package main

import (
	"bufio"
	"flag"
	"fmt"
	// "net/http/httputil"
	// "net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Francois-Coleongco/LoadBalancer/types"
)

func main() {

	var file_name *string = flag.String("-f", "", "please enter your server file after the -f")

	flag.Parse()

	s := types.InitServers()

	file, err := os.Open(*file_name)

	if err != nil {
		fmt.Println("error opening file with name: ", *file_name)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		read_in_server := strings.Split(line, ":")
		read_in_port, err := strconv.ParseInt(read_in_server[1], 10, 16)

		if err != nil {
			fmt.Println("could not parse port of server: ", read_in_server)
			continue
		}

		s.AddToFront(read_in_server[0], uint16(read_in_port))
	}

	if s.Size == 0 {
		fmt.Println("no servers read in :( how can we possibly balance now! NOOOOOOOOO")
		return
	}

	// start the http proxy
	// read in new requests passing off to servers based on an injected load balancing session cookie for backend agnosticism
	// set this cookie in the response then the client will send it every time automatically. then just route using that cookie
	// every time the user sends a request, the user will have their session info (specific to the backend whatever it may be) stored in a Redis instance.
	// if the load balancer decides it is optimal, they may pipe a subsequent request by this user to another server safely as the session data is stored in redis, not local to a backend server

	// u := new(url.URL)

	// httputil.NewSingleHostReverseProxy()

}
