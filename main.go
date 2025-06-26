package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func getClientIPv4(r *http.Request) string {
	// Ưu tiên lấy từ các header thường dùng khi có proxy
	headers := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			for _, ip := range strings.Split(ips, ",") {
				ip = strings.TrimSpace(ip)
				parsed := net.ParseIP(ip)
				if parsed != nil && parsed.To4() != nil {
					return ip // Trả IP v4 đầu tiên hợp lệ
				}
			}
		}
	}

	// Fallback: lấy từ RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	parsed := net.ParseIP(ip)
	if parsed != nil && parsed.To4() != nil {
		return ip
	}

	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip := getClientIPv4(r)
	if ip == "" {
		http.Error(w, "IPv4 address not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ip)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
