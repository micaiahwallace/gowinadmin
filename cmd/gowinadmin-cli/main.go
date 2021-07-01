package main

import (
	"flag"
	"log"

	"github.com/micaiahwallace/gowinadmin"
)

func main() {

	// Get cli arguments
	var host string
	flag.StringVar(&host, "host", "127.0.0.1", "remote host to execute")
	flag.Parse()

	// Create and run the command
	request := gowinadmin.QuserRequest{Server: host}
	sessions, runErr := request.RunQuser()
	if runErr != nil {
		log.Fatalf("Quser failed: %v\n", runErr)
	}

	// Display results
	log.Println(sessions)
}
