package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// IPAddresses holds the IPv4 and IPv6 addresses
type IPAddresses struct {
	V4 string `json:"v4"`
	V6 string `json:"v6"`
}

func getLocalIPs() (string, string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
		return "", ""
	}

	var ipv4, ipv6 string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
			} else if ipnet.IP.To16() != nil {
				ipv6 = ipnet.IP.String()
			}
		}
	}
	return ipv4, ipv6
}

func main() {
	ipv4, ipv6 := getLocalIPs()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", ipv4)
	})

	http.HandleFunc("/v6", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", ipv6)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ipAddrs := IPAddresses{
			V4: ipv4,
			V6: ipv6,
		}
		json.NewEncoder(w).Encode(ipAddrs)
	})

	port := "8081"
	fmt.Printf("Server starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}