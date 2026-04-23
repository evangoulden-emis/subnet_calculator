package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"net"
)

const (
	TOTAL_BITS = 32
)

func main() {
	cidr := flag.String("cidr", "192.168.0.0/23", "Enter a CIDR, valid formats are: 10.0.0.0/8 or 10.0.0.0 255.0.0.0")
	subnetwork := flag.Int("sub", 27, "Enter the bitmask to subdivide the CIDR into subnets e.g. 27")
	flag.Parse()

	fmt.Printf("The following options were entered, CIDR: %s, Subnetwork: %s \n", *cidr, *subnetwork)
	ip, _, err := net.ParseCIDR(*cidr)
	if err != nil {
		fmt.Println(err)
		return
	}
	u32BaseIp := binary.BigEndian.Uint32(ip.To4())
	// Calculate the subnet jump size

	jmpSize := math.Pow(2, float64(TOTAL_BITS-*subnetwork))
	origJmpSize := jmpSize

	count := 0
	var sub []uint32
	//subnets := make([]net.IPMask, 4)
	for {
		//if jmpSize > 255 {
		//	break
		//}

		if count == 0 {
			// First Itteration
			sub = append(sub, u32BaseIp)
			count++
		} else {
			sub = append(sub, u32BaseIp+uint32(jmpSize))
			count++
			jmpSize = jmpSize + origJmpSize
		}
	}
	convertedSub := convertUint32ToIpAddress(sub)
	printSubnets(convertedSub, *subnetwork)

}

func convertUint32ToIpAddress(u []uint32) []net.IP {
	ips := make([]net.IP, 0, len(u))
	for _, i := range u {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}
	return ips
}

func printSubnets(subnets []net.IP, netmask int) {
	for _, subnet := range subnets {
		network := net.IPNet{
			IP:   subnet,
			Mask: net.CIDRMask(netmask, 32),
		}
		fmt.Println(network.String())
	}
}
