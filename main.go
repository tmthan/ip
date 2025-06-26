package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			return strings.Split(ips, ",")[0]
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	fmt.Fprintf(w, "Your IP address is: %s\n", clientIP)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
