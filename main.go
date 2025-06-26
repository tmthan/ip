package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func extractIPs(r *http.Request) []net.IP {
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

	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := net.ParseIP(ipStr); ip != nil {
			result = append(result, ip)
		}
	}

	return result
}

func getIPv4(ips []net.IP) string {
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

func getIPv6(ips []net.IP) string {
	for _, ip := range ips {
		// IP không phải IPv4, nhưng là hợp lệ thì là IPv6
		if ip.To4() == nil && ip.To16() != nil {
			return ip.String()
		}
	}
	return ""
}

func ipv4Handler(w http.ResponseWriter, r *http.Request) {
	ip := getIPv4(extractIPs(r))
	if ip == "" {
		http.Error(w, "IPv4 not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ip)
}

func ipv6Handler(w http.ResponseWriter, r *http.Request) {
	ips := extractIPs(r)
	ipv6 := getIPv6(ips)
	if ipv6 == "" {
		ipv6 = getIPv4(ips)
	}
	if ipv6 == "" {
		http.Error(w, "No IP found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ipv6)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	ips := extractIPs(r)
	res := map[string]string{
		"v4": getIPv4(ips),
		"v6": getIPv6(ips),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/", ipv4Handler)
	http.HandleFunc("/v6", ipv6Handler)
	http.HandleFunc("/json", jsonHandler)

	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
