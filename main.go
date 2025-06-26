package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// Lấy danh sách IP (net.IP) từ các header phổ biến và RemoteAddr
func extractIPs(r *http.Request) []net.IP {
	headers := []string{"X-Forwarded-For", "X-Real-IP"}
	var ips []net.IP

	for _, header := range headers {
		raw := r.Header.Get(header)
		if raw != "" {
			for _, ipStr := range strings.Split(raw, ",") {
				ipStr = strings.TrimSpace(ipStr)
				if ip := net.ParseIP(ipStr); ip != nil {
					ips = append(ips, ip)
				}
			}
		}
	}

	// Thêm RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := net.ParseIP(host); ip != nil {
			ips = append(ips, ip)
		}
	}

	return ips
}

// Tìm IPv4 đầu tiên trong danh sách
func getIPv4(ips []net.IP) string {
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

// Tìm IPv6 đầu tiên trong danh sách
func getIPv6(ips []net.IP) string {
	for _, ip := range ips {
		if ip.To4() == nil && ip.To16() != nil {
			return ip.String()
		}
	}
	return ""
}

// Handler trả IPv4 hoặc lỗi nếu không có
func ipv4Handler(w http.ResponseWriter, r *http.Request) {
	ips := extractIPs(r)
	ipv4 := getIPv4(ips)
	if ipv4 == "" {
		http.Error(w, "IPv4 address not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ipv4)
}

// Handler trả IPv6 hoặc fallback IPv4, lỗi nếu không có
func ipv6Handler(w http.ResponseWriter, r *http.Request) {
	ips := extractIPs(r)
	ipv6 := getIPv6(ips)
	if ipv6 == "" {
		ipv6 = getIPv4(ips)
	}
	if ipv6 == "" {
		http.Error(w, "No IPv6 or IPv4 address found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ipv6)
}

// Handler trả JSON {v4: "...", v6: "..."}
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	ips := extractIPs(r)
	resp := map[string]string{
		"v4": getIPv4(ips),
		"v6": getIPv6(ips),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", ipv4Handler)
	http.HandleFunc("/v6", ipv6Handler)
	http.HandleFunc("/json", jsonHandler)

	log.Println("Server listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
