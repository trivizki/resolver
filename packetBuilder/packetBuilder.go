package packetBuilder

import (
	"net"
	"errors"
	"resolver/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	responseTTL uint32 = 20
)

// PacketBuilder is responsible to build dns response packet.
type PacketBuilder struct{}

func NewPacketBuilder() *PacketBuilder{
	return &PacketBuilder{}
}

// This function builds a dns response packet from the given data.
// @ srcMac - the respsonse packet source mac
// @ dstMac - the respsonse packet dest mac
// @ srcIP - the respsonse packet souce IP 
// @ dstIP - the respsonse packet dest IP 
// @ srcPort - the respsonse packet source port (UDP)
// @ dstPort - the respsonse packet dest port (UDP)
// @ dnsRequest - the DNS query that we want to response to 
// @ answers - the ip addresses associated with given quried domain. 
func (p *PacketBuilder) BuildPacket(srcMac net.HardwareAddr, dstMac net.HardwareAddr,
						srcIP net.IP, dstIP net.IP, srcPort layers.UDPPort,
						dstPort layers.UDPPort, dnsRequest *layers.DNS, answers []net.IP) ([]byte, error){
  // Lets fill out some information
    ipLayer := &layers.IPv4{
	  SrcIP: srcIP, 
	  DstIP: dstIP, 
	  Version: 4,
	  Protocol: layers.IPProtocolUDP,
	  Flags: layers.IPv4DontFragment,
	  TTL: 64,
	}
	ethernetLayer := &layers.Ethernet{
	      SrcMAC: srcMac,
	      DstMAC: dstMac, 
		  EthernetType: layers.EthernetTypeIPv4,
	}
	udpLayer := &layers.UDP{
	      SrcPort: srcPort, 
	      DstPort: dstPort, 
	}

	ip4answers := utils.Onlyipv4(answers)
	if (len(ip4answers) < 1){
		return nil, errors.New("no ipv4 answer")
	}
	dnsLayer := p.buildDnsResponse(dnsRequest, ip4answers)
	udpLayer.SetNetworkLayerForChecksum(ipLayer)

	// Create the packet with the layers
	options := gopacket.SerializeOptions{
		FixLengths: true,
		ComputeChecksums: true,
	}

	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, options,
	    ethernetLayer,
	    ipLayer,
		udpLayer,
		dnsLayer,
	)
	return buffer.Bytes(), err
}
// Builds DNS response payload according to the given dnsRequrst and the ip addreess
// that assocciated with the query.
func (p *PacketBuilder) buildDnsResponse(dnsRequest *layers.DNS, answers []net.IP) *layers.DNS {
	dnsResponse := &layers.DNS{}
	dnsResponse.ID = dnsRequest.ID
	dnsResponse.QR = true 
	dnsResponse.OpCode = dnsRequest.OpCode
	dnsResponse.RD = dnsRequest.RD
	dnsResponse.RA = true 
	dnsResponse.Questions = make([]layers.DNSQuestion, 1)
	dnsResponse.Questions[0] = dnsRequest.Questions[0]
	dnsResponse.Answers = make([]layers.DNSResourceRecord, len(answers))

	for i, ip := range answers {
		dnsResponse.Answers[i] = layers.DNSResourceRecord{}
		dnsResponse.Answers[i].Name = dnsRequest.Questions[0].Name
		dnsResponse.Answers[i].Type = layers.DNSTypeA
		dnsResponse.Answers[i].Class = layers.DNSClassIN
		dnsResponse.Answers[i].TTL = responseTTL 
		dnsResponse.Answers[i].IP = ip
	}
	return dnsResponse
}
