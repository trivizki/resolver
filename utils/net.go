package utils

import (
	"fmt"
	"net"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Thus function formats net.IP array to string.
func FormatIPS(ips []net.IP)string{
	format := ""
	for _, ip := range ips{
		format = fmt.Sprintf("%s %s", ip.String(), format)
	}
	return format
}

func CreatePacketFromByte(packet []byte)gopacket.Packet{
	return gopacket.NewPacket(packet, layers.LayerTypeEthernet, gopacket.Default)
}

func CreateDNSLayer(id uint16,  opCode layers.DNSOpCode, rd bool, name []byte) *layers.DNS{
	dnsLayer := &layers.DNS{}
	dnsLayer.ID = id  
	dnsLayer.OpCode = opCode 
	dnsLayer.RD = rd 
	dnsLayer.Questions = make([]layers.DNSQuestion, 1)
	dnsLayer.Questions[0] = layers.DNSQuestion{
		Name: name,
		Type: layers.DNSTypeA,
		Class: layers.DNSClassIN,
	}
	return dnsLayer
}

// Takes only ipv4 addresses from given ip list.
func Onlyipv4(ips []net.IP) []net.IP{
	var ipv4 []net.IP
	for _, ip := range ips {
		if ip.To4() == nil{
			// Not ipv4 (probably Ipv6)
			continue
		}
		ipv4 = append(ipv4, ip)
	}
	return ipv4
}


