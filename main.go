package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func extractIPs(r *http.Request) []net.IP {
	// Ưu tiên các header từ proxy/nginx
	headers := []string{"X-Forwarded-For", "X-Real-IP"}
	var result []net.IP

	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			for _, ip := range strings.Split(ips, ",") {
				parsed := net.ParseIP(strings.TrimSpace(ip))
				if parsed != nil {
					result = append(result, parsed)
				}
			}
		}
	}

	// Nếu không có, dùng RemoteAddr
	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := net.ParseIP(ipStr); ip != nil {
			result = append(result, ip)
		}
	}

	return result
}

func getPreferredIP(r *http.Request, preferV6 bool) string {
	ips := extractIPs(r)
	var fallback string

	for _, ip := range ips {
		if ip.To4() != nil {
			if !preferV6 {
				return ip.String() // Trả v4 nếu không ưu tiên v6
			}
			if fallback == "" {
				fallback = ip.String() // lưu v4 làm fallback
			}
		} else {
			if preferV6 {
				return ip.String() // Trả v6 nếu được
			}
		}
	}

	return fallback
}

func ipv4Handler(w http.ResponseWriter, r *http.Request) {
	ip := getPreferredIP(r, false)
	if ip == "" {
		http.Error(w, "No IP found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ip)
}

func ipv6Handler(w http.ResponseWriter, r *http.Request) {
	ip := getPreferredIP(r, true)
	if ip == "" {
		http.Error(w, "No IP found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ip)
}

func main() {
	http.HandleFunc("/", ipv4Handler)
	http.HandleFunc("/v6", ipv6Handler)

	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
