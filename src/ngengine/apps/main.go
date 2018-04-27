package main

import (
	"fmt"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", "0.0.0.0:0") // listen on localhost
	port := l.Addr().(*net.TCPAddr).Port
	ip := l.Addr().(*net.TCPAddr).IP
	fmt.Println(ip, port)

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {

		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPAddr:
				fmt.Println(v.IP)
			}

		}
	}

}
