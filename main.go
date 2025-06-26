package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func getClientIPV6(r *http.Request) string {
	// Check các header thường được dùng khi đi qua proxy hoặc load balancer
	headers := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			// Lấy IP đầu tiên nếu có nhiều IP
			return strings.Split(ips, ",")[0]
		}
	}

	// Nếu không có header, lấy từ RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func handlerV6(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIPV6(r)
	fmt.Fprintf(w, clientIP)
}


func getClientIPV4(r *http.Request) string {
	// Check các header có thể chứa IP
	headers := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			for _, ip := range strings.Split(ips, ",") {
				ip = strings.TrimSpace(ip)
				if parsed := net.ParseIP(ip); parsed != nil && parsed.To4() != nil {
					return ip // Trả về IP v4 đầu tiên hợp lệ
				}
			}
		}
	}

	// Nếu không có header, lấy từ RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	parsed := net.ParseIP(ip)
	if parsed != nil && parsed.To4() != nil {
		return ip // IP v4
	}
	return "" // Nếu không tìm thấy IPv4
}

func handlerV4(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIPV4(r)
	if clientIP != "" {
		fmt.Fprint(w, clientIP)
	} else {
		http.Error(w, "IPv4 address not found", http.StatusInternalServerError)
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	v4 := getClientIPV4(r)
	v6 := getClientIPV6(r)

	response := map[string]string{
		"v4": v4,
		"v6": v6,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handlerV4)
	http.HandleFunc("/v6", handlerV6)
	http.HandleFunc("/json", jsonHandler)

	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
