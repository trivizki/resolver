package packetBuilder

import (
	"bytes"
	"testing"
	"net"
	"encoding/hex"
	"resolver/utils"
	"github.com/google/gopacket/layers"
)

func TestBuildPacket(t *testing.T){
	p, _ := hex.DecodeString("f0d5bf41eba0a0ab1b571a2408000500005a000040004011f73ac0a80101c0a801070035d9f700464b7ff637818000010001000000000377777706676f6f676c6503636f6d00000100010377777706676f6f676c6503636f6d0000010001000000140004acd91664")
	tests := []struct{
		srcMac net.HardwareAddr
		dstMac net.HardwareAddr
		srcIP net.IP
		dstIP net.IP
		srcPort layers.UDPPort
		dstPort layers.UDPPort
		dnsRequest *layers.DNS
		answers []net.IP
		packet []byte 
		err error
	}{
		{
			srcMac:net.HardwareAddr{0xA0, 0xAB, 0x1B, 0x57, 0x1A, 0x24},
			dstMac:net.HardwareAddr{0xF0, 0xD5, 0xBF, 0x41, 0xEB, 0xA0},
			srcIP:net.IP{192, 168, 1, 1},
			dstIP:net.IP{192, 168, 1, 7},
			srcPort:layers.UDPPort(53),
			dstPort:layers.UDPPort(55799),
			dnsRequest:utils.CreateDNSLayer(0xf637, layers.DNSOpCodeQuery, true, []byte("www.google.com")),
			answers:[]net.IP{net.IP{172, 217, 22,100}},
			packet: p,
			err: nil,
		},
	}
	for _, test := range tests{
		builder := NewPacketBuilder() 
		r,e := builder.BuildPacket(test.srcMac, test.dstMac, test.srcIP, test.dstIP,
										test.srcPort, test.dstPort, test.dnsRequest,
											test.answers)
		if (bytes.Compare(r,test.packet) != 0){
			t.Errorf("got different packet\n %s", hex.Dump(r))
		}
		if (e != test.err){
			t.Errorf("got different error")
		}
	}
}
