package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// Hàm gốc của bạn để lấy IP từ header hoặc RemoteAddr
func getRawClientIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range headers {
		ips := r.Header.Get(header)
		if ips != "" {
			return strings.TrimSpace(strings.Split(ips, ",")[0])
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Phân loại IP v4 / v6
func classifyIP(rawIP string) (v4 string, v6 string) {
	parsedIP := net.ParseIP(rawIP)
	if parsedIP == nil {
		return "", ""
	}

	if parsedIP.To4() != nil {
		return parsedIP.String(), ""
	}
	return "", parsedIP.String()
}

func ipv4Handler(w http.ResponseWriter, r *http.Request) {
	rawIP := getRawClientIP(r)
	v4, _ := classifyIP(rawIP)
	if v4 == "" {
		http.Error(w, "IPv4 not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, v4)
}

func ipv6Handler(w http.ResponseWriter, r *http.Request) {
	rawIP := getRawClientIP(r)
	_, v6 := classifyIP(rawIP)
	if v6 == "" {
		// fallback: trả v4 nếu không có v6
		v4, _ := classifyIP(rawIP)
		if v4 == "" {
			http.Error(w, "No valid IP found", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, v4)
		return
	}
	fmt.Fprint(w, v6)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	rawIP := getRawClientIP(r)
	v4, v6 := classifyIP(rawIP)

	response := map[string]string{
		"v4": v4,
		"v6": v6,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", ipv4Handler)
	http.HandleFunc("/v6", ipv6Handler)
	http.HandleFunc("/json", jsonHandler)

	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
