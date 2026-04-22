package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
)

const (
	TOTAL_BITS = 32
)

func main() {
	cidr := flag.String("cidr", "192.168.1.0/24", "Enter a CIDR, valid formats are: 10.0.0.0/8 or 10.0.0.0 255.0.0.0")
	subnetwork := flag.Int("sub", 27, "Enter the bitmask to subdivide the CIDR into subnets e.g. 27")
	flag.Parse()

	fmt.Printf("The following options were entered, CIDR: %s, Subnetwork: %s \n", *cidr, *subnetwork)
	ip, _, err := net.ParseCIDR(*cidr)
	if err != nil {
		fmt.Println(err)
		return
	}
	u32BaseIp := binary.BigEndian.Uint32(ip.To4())
	fmt.Println(u32BaseIp)
	// Calculate the subnet jump size

	jmpSize := math.Pow(2, float64(TOTAL_BITS-*subnetwork))
	origJmpSize := jmpSize

	if jmpSize > 255 {

		os.Exit(1)
	}
	fmt.Println(jmpSize)
	count := 0
	//subnets := make([]net.IPMask, 4)
	for {
		if jmpSize > 255 {
			break
		}
		var sub []uint32
		if count == 0 {
			// First Itteration
			sub = append(sub, u32BaseIp)
			count++
		} else {
			sub = append(sub, u32BaseIp+uint32(jmpSize))
			count++
			jmpSize = jmpSize + origJmpSize
		}

		fmt.Println(sub)

	}
}
