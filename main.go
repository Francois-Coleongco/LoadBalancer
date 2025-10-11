package main

import (
	"bufio"
	"crypto"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

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

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			server, err := s.GetServer()
			if err != nil {
				// only reason for GetServer to not work is if there are no servers
				return
			}
			target := server.URL + ":" + strconv.Itoa(int(server.PORT))
			u, _ := url.Parse(target)
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host

		},

		ModifyResponse: func(resp *http.Response) error {

			// use current time + ip addr as hash information
			currTime := time.Now().UTC().String()
			addrPortTime := resp.Request.RemoteAddr + currTime
			h := crypto.SHA256.New()
			h.Write([]byte(addrPortTime))
			hash := hex.EncodeToString(h.Sum(nil))

			tracking := &http.Cookie{
				HttpOnly: true,
				Name:     "LB_Tracker",
				Value:    hash,
			}

			resp.Header.Add("Set-Cookie", tracking.String())
			// save cookie value for session in redis
			// should store the LB_Tracker as key and session data as value AS WELL AS the current time to track time of last request to determine a dead or alive connection

			return nil

		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy Error: %v\n", err)
			http.Error(w, "Backend unavailable", http.StatusBadGateway)
		},
	}

	log.Println("Starting proxy")

	http.ListenAndServe(":"+strconv.Itoa(*lb_port), &proxy)

}
