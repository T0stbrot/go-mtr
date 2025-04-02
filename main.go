package main

import (
	"os"
	"fmt"
	"net/netip"
	"encoding/json"
	"net/http"
	ping "github.com/t0stbrot/go-ping"
)

type Details struct {
	IP string `json:"ip"`
	RTT string `json:"rtt"`
	ASN int `json:"asn,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

func details(addr string) (details Details) {
	url := fmt.Sprintf("https://t0stbrot.net/pub-api/ip/%v", addr)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching from API")
		os.Exit(1)
	}
	defer res.Body.Close()

	var response Details
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding API JSON")
		os.Exit(1)
	}

	return response
}

func ipIs(addr string) (is int) {
	ip, _ := netip.ParseAddr(addr)
	if ip.Is6() {
		return 6
	} else if ip.Is4() {
		return 4
	} else {
		return 0
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("You need to provide a Target Address")
		os.Exit(2)
	} else {
		target := args[1]

		if ipIs(target) == 0 {
			fmt.Println("You need to provide a Valid Target Address")
		} else {

			msg := fmt.Sprintf("Starting MTR to %v with 255 maximum hops", target)
			fmt.Println(msg)

			hops := 0
				
			for hops <= 255 {
				var res ping.PingResult
				var info Details
				if ipIs(target) == 6 {
					res = ping.Ping6(target, hops, 1000)
				} else {
					res = ping.Ping4(target, hops, 1000)
				}
				hops++
				ip, _ := netip.ParseAddr(res.LastHop)
				if ip.IsValid() {
					if ip.IsPrivate() {
						info = Details{
							IP: res.LastHop,
							RTT: res.RTT+"ms",
						}
					} else {
						info = details(res.LastHop)
						info.RTT = res.RTT+"ms"
					}
				} else {
					info = Details{
						IP: res.Error,
						RTT: "Timeout",
					}
				}
				var host string
				if info.Hostname != "" {
					host = "[" + info.Hostname + "]"
				}
				msg = fmt.Sprintf("[%2d] [%10s] [AS%6v] %v %v", hops-1, info.RTT, info.ASN, info.IP, host)
				fmt.Println(msg)
				if res.LastHop == target {
					break
					os.Exit(0)
				}
			}
		}
	}

}