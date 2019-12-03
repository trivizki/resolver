package context

import (
	"net"
	"fmt"
	"bytes"
	"testing"
	"encoding/hex"
	"resolver/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestBuildContext(t *testing.T){
	p, _ := hex.DecodeString("a0ab1b571a24f0d5bf41eba008004500003cbe8f40004011f8c8c0a80107c0a80101b1dd003500284b5cf02e010000010000000000000377777706676f6f676c6503636f6d0000010001")
	tests := []struct{
		packet gopacket.Packet
		context *Context
		err error
	}{
		{
			utils.CreatePacketFromByte(p),
			&Context{
				SrcMAC:net.HardwareAddr{0xF0, 0xD5, 0xBF, 0x41, 0xEB, 0xA0},
				DstMAC:net.HardwareAddr{0xA0, 0xAB, 0x1B, 0x57, 0x1A, 0x24},
				SrcIP:net.IP{192, 168, 1, 7},
				DstIP:net.IP{192, 168, 1, 1},
				SrcPort:layers.UDPPort(45533),
				DstPort:layers.UDPPort(53),
				Payload:utils.CreateDNSLayer(0xf02e, layers.DNSOpCodeQuery, true, []byte("www.google.com")),
			},
			nil,
		},
	}
	cb := NewContextBuilder()
	for _, test := range tests {
		c, e := cb.BuildContext(test.packet)
		if (compareContexts(*c, *test.context) != true){
			t.Errorf("got diffrent context")
		}
		if (e != test.err){
			t.Errorf("got different error")
			fmt.Println(e.Error())
		}
	}
}

// This function compare between two diffrent context objects.
func compareContexts(c1 Context, c2 Context) bool{
	var err error 
	if (bytes.Compare(c1.SrcMAC,c2.SrcMAC) != 0) {return false}
	if (bytes.Compare(c1.DstMAC,c2.DstMAC) != 0) {return false}
	if (!(c1.SrcIP.Equal(c2.SrcIP))) {return false}
	if (!(c1.DstIP.Equal(c2.DstIP))) {return false}
	if (c1.SrcPort != c2.SrcPort) {return false}
	if (c1.DstPort != c2.DstPort) {return false}

	// Serialize DNS layer in order to compare between byte arrayes
	options := gopacket.SerializeOptions{
		FixLengths: true,
	}
	buffer1 := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer1, options, c1.Payload)
	if (err != nil){return false}
	buffer2 := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer2, options, c2.Payload)
	if (err != nil){return false}
	if (bytes.Compare(buffer1.Bytes(),buffer2.Bytes()) != 0) {return false}

	return true
}
