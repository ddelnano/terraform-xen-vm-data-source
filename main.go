package main

import (
	"encoding/json"
	"log"
	"os"

	xenapi "github.com/terra-farm/go-xen-api-client"
)

type ExternalProgramProtocol struct {
	Query map[string]string
}

func main() {
	var program ExternalProgramProtocol
	output := map[string]interface{}{}

	decoder := json.NewDecoder(os.Stdin)
	decoder.Decode(&program)

	username := os.Getenv("XAPI_USERNAME")
	password := os.Getenv("XAPI_PASSWORD")
	host := os.Getenv("XAPI_HOST")

	if username == "" || password == "" || host == "" {
		log.Fatalf("XAPI_HOST, XAPI_PASSWORD and XAPI_HOST variables must be set")
	}

	client, err := xenapi.NewClient("https://"+host, nil)

	if err != nil {
		log.Fatalf("failed to create client with error: %v", err)
	}

	session, err := client.Session.LoginWithPassword(username, password, "1.0", "terraform")
	if err != nil {
		log.Fatalf("failed to login with error: %v", err)
	}

	for resourceName, uuid := range program.Query {

		instance, err := client.VM.GetByUUID(session, uuid)

		if err != nil {
			log.Fatalf("failed to get vm by ID with error: %v", err)
		}

		// log.Printf("Found the following VM: %v", instance)

		m, err := client.VM.GetGuestMetrics(session, instance)

		if err != nil {
			log.Fatalf("failed to retrieve guest metrics with error: %v", err)
		}

		metrics, err := client.VMGuestMetrics.GetRecord(session, m)

		if err != nil {
			log.Fatalf("failed to retrieve guest metrics record with error: %v", err)
		}

		// log.Printf("Found the following metrics: %v", metrics)

		networks := metrics.Networks
		nets := map[string]string{}

		// TODO: These keys can be parsed more intelligently
		// The last number is the device order (network 0, 1 , etc)
		if ip, ok := networks["0/ip"]; ok {
			if ip != "" {
				nets["ip_address"] = ip
			}
		}

		if ip, ok := networks["0/ipv4/0"]; ok {
			if ip != "" {
				nets["ipv4_address"] = ip
			}
		}

		if ip, ok := networks["0/ipv6/0"]; ok {
			if ip != "" {
				nets["ipv6_address"] = ip
			}
		}

		output[resourceName] = nets
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(output)
}
