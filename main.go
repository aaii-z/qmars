package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Function to scan a given IP and port
func scanPort(ip string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Function to get local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// Check the address type and ensure it is not a loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// Function to generate a list of IPs in the subnet
func generateIPs(subnet string) []string {
	ips := []string{}
	for i := 40; i < 255; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", subnet, i))
	}
	return ips
}

func main() {
	// Define the port to search for
	port := "5457"
	
	// Get the local IP address and subnet
	localIP := getLocalIP()
	if localIP == "" {
		fmt.Println("Unable to determine local IP address.")
		return
	}
	fmt.Printf("Local IP: %s\n", localIP)

	// Extract the subnet by trimming the last octet of the IP address
	subnet := localIP[:strings.LastIndex(localIP, ".")]

	// Generate list of IPs in the subnet
	ips := generateIPs(subnet)

	// Scan the subnet for open port 5457
	found := false
	for i, ip := range ips {
		fmt.Printf("Scanning %s...\n", ip) // Added verbose logging
		if scanPort(ip, port) {
			fmt.Printf("Qmars%d: %s has port %s open\n", i+1, ip, port)
			found = true
		}
		time.Sleep(5 * time.Millisecond) // Added delay of 500 milliseconds between scans
	}

	// If no devices found, start listening on port 5457
	if !found {
		fmt.Printf("No Qmars found. Starting server on port %s...\n", port)
		ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ln.Close()

		// Wait for other Qmars to connect
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("A Qmars has joined!")
			conn.Close()
		}
	}
}

