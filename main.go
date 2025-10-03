package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"os"
	"strconv"
	"strings"

	"github.com/Francois-Coleongco/LoadBalancer/types"
)

func main() {

	var file_name *string = flag.String("f", "", "please enter your server file after the -f")

	var lb_port *int = flag.Int("p", 6969, "please enter a valid port 0-65535")

	if *lb_port < 0 || 65535 < *lb_port {
		log.Println("please provide a valid port 0-65535")
	}

	flag.Parse()

	s := types.InitServers()

	file, err := os.Open(*file_name)

	if err != nil {
		log.Println("error opening file with name: ", *file_name)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		colon_location := strings.LastIndex(line, ":")
		read_in_server := line[:colon_location]
		read_in_port, err := strconv.ParseInt(line[colon_location+1:], 10, 16) // plus one cuz colon is included

		if err != nil {
			log.Println("could not parse port of server: ", read_in_server)
			continue
		}

		s.AddToFront(read_in_server, uint16(read_in_port))
	}

	if s.Size == 0 {
		log.Println("no servers read in :( how can we possibly balance now! NOOOOOOOOO")
		return
	}

	// start the http proxy
	// read in new requests passing off to servers based on an injected load balancing session cookie for backend agnosticism
	// set this cookie in the response then the client will send it every time automatically. then just route using that cookie
	// every time the user sends a request, the user will have their session info (specific to the backend whatever it may be) stored in a Redis instance.
	// if the load balancer decides it is optimal, they may pipe a subsequent request by this user to another server safely as the session data is stored in redis, not local to a backend server

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			server := s.GetServer()
			target := server.URL + ":" + strconv.Itoa(int(server.PORT))
			u, _ := url.Parse(target)
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy Error: %v\n", err)
			http.Error(w, "Backend unavailable", http.StatusBadGateway)
		},
	}

	log.Println("Starting proxy")

	http.ListenAndServe(":"+strconv.Itoa(*lb_port), &proxy)

}
