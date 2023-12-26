package ipaddress

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

//  its Working For local
func GetLocalIP() (string, error) {
	cmd := exec.Command("hostname", "-I")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	ipAddresses := strings.Fields(string(output))
	if len(ipAddresses) > 0 {
		return ipAddresses[0], nil
	}

	return "", fmt.Errorf("no local network IP address found")
}

//  it will work based on Request  on Event service still Need to verify
func handleRemote(w http.ResponseWriter, r *http.Request) {
	log.Println("handleRemote (+)")
	log.Println("r.RemoteAddr", r.RemoteAddr)
	ipAddress := r.RemoteAddr
	fmt.Println("Client IP Address:", ipAddress)
	// Your handling logic here
	log.Println("handleRemote (-)")
}

//  if Device Address is comming with http
func handleHttp(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	}

	fmt.Println("Client IP Address:", ipAddress)

	// Your handling logic here
}
