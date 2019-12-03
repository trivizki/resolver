package context

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
	"net"
	"fmt"
)


// Context holds all the relevant data on the packet in the current flow
// This is the main data structure of the system that supports the system logic.
type Context struct {
	Timestamp time.Time
	SrcMAC net.HardwareAddr
	DstMAC net.HardwareAddr
	SrcIP net.IP
	DstIP net.IP
	SrcPort layers.UDPPort
	DstPort layers.UDPPort
	Payload *layers.DNS
}

func (c *Context) IsQuery() bool{
	return c.Payload.QR == false
}

func (c *Context) IsAQuery() bool{
	if !(c.IsQuery()){
		return false
	}
	if c.Payload.QDCount != 1{
		return false
	}
	q := c.Payload.Questions[0]
	if q.Type == layers.DNSTypeA{
		return true
	}
	return false
}

// ContextBuilder is responsible to create The packet's context.
type ContextBuilder struct {
}

// Creates new context builder object.
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}

/*
	BuildContext create the UDP packet's context from original packet.
	Basiclly this function parse udp packet.
*/
func (cb *ContextBuilder) BuildContext (packet gopacket.Packet) (*Context, error){
	c := &Context{}
	c.Timestamp = packet.Metadata().CaptureInfo.Timestamp

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer == nil {
		return c, fmt.Errorf("Could not find ethernet layer")
	}
	ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	c.SrcMAC = ethernetPacket.SrcMAC
	c.DstMAC = ethernetPacket.DstMAC

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return c, fmt.Errorf("Could not find ip layer")
	}

	ip, _ := ipLayer.(*layers.IPv4)
	c.SrcIP = ip.SrcIP
	c.DstIP = ip.DstIP

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return c, fmt.Errorf("Could not find udp layer")
	}
	udp, _ := udpLayer.(*layers.UDP)
	c.SrcPort = udp.SrcPort
	c.DstPort = udp.DstPort

	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer == nil {
		return c, fmt.Errorf("Could not find dns layer")
	}
	dns, _ := dnsLayer.(*layers.DNS)
	c.Payload = dns

	return c, nil
}
