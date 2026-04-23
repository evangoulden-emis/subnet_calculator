package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
)

const (
	IPv4AddressSize = 32
)

func main() {
	cidr := flag.String("cidr", "10.0.0.0/22", "Enter a CIDR, valid formats are: 10.0.0.0/8 or 10.0.0.0 255.0.0.0")
	subnetwork := flag.Int("sub", 0, "Enter the bitmask to subdivide the CIDR into subnets e.g. 27")
	flag.Parse()

	fmt.Printf("The following options were entered, CIDR: %s, Subnetwork: %d \n", *cidr, *subnetwork)
	ip, ipnet, err := net.ParseCIDR(*cidr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *subnetwork == 0 {
		*subnetwork, _ = ipnet.Mask.Size()
	}

	if *subnetwork > IPv4AddressSize || *subnetwork <= 1 {
		fmt.Println("Invalid Subnetwork, you have entered a value that is either greater than 32 or less than 1")
	}
	u32BaseIp := binary.BigEndian.Uint32(ip.To4())
	maskSize, _ := ipnet.Mask.Size()
	// Calculate the subnet jump size dynamically
	jmpSize := uint32(1) << (IPv4AddressSize - *subnetwork)

	// Calculate the total number of IPs in the original CIDR block
	totalIPs := uint32(1) << (IPv4AddressSize - maskSize)

	// Calculate the upper limit for iteration
	upperLimit := u32BaseIp + totalIPs

	// Calculate the subnet jump size
	var subnets []uint32
	for currentIP := u32BaseIp; currentIP < upperLimit; currentIP += jmpSize {
		subnets = append(subnets, currentIP)
	}

	convertedSub := convertUint32ToIpAddress(subnets)
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
	fmt.Printf("\n%-18s %-30s %-5s %-15s\n", "Subnet CIDR", "Start IP - End IP", "Usable IPs", "Broadcast")
	for _, subnet := range subnets {
		network := net.IPNet{
			IP:   subnet,
			Mask: net.CIDRMask(netmask, 32),
		}
		startIP := binary.BigEndian.Uint32(subnet) + 1
		endIP := binary.BigEndian.Uint32(subnet) + uint32(1<<(32-netmask)) - 2
		broadcastIP := binary.BigEndian.Uint32(subnet) + uint32(1<<(32-netmask)) - 1

		fmt.Printf("%-18s %-30s %-5d %-15s\n",
			network.String(),
			fmt.Sprintf("%s - %s", toIPAddr(startIP), toIPAddr(endIP)), int((endIP-startIP)+1),
			toIPAddr(broadcastIP),
		)

	}
}

func toIPAddr(u32 uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, u32)
	return ip
}
