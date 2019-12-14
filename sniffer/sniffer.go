package sniffer 

import (
		"fmt"
		"time"
		"resolver/logger"
		"github.com/google/gopacket"
		"github.com/google/gopacket/pcap"
)


var (
		snapshotLen int32  = 2048
		promiscuous bool   = false
		timeout     time.Duration = -1 * time.Millisecond 
)

type SnifferConf struct{
	Device string
	Filter string
}

// Sniffer is responsible to sniff raw packet from the configured interface.
type Sniffer struct {
	packetChannel chan<- gopacket.Packet
	handle        *pcap.Handle
	conf SnifferConf
	logger *logger.Logger
}

// Creates new sniffer object.
// The sniffer will send every packet in the given channel.
func NewSniffer(packetChannel chan<- gopacket.Packet, conf SnifferConf,logger *logger.Logger ) (*Sniffer) {
	return &Sniffer{
		packetChannel : packetChannel,
		handle : nil,
		conf: conf,
		logger: logger,
	}
}

// Initialize the sniffer by opening a handle to the relevant interface.
// NOTICE: user must call this function before using the Sniffer.
func (s *Sniffer) InitializeSniffer() error{
	// Open device
	handle, err := pcap.OpenLive(s.conf.Device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return err
	}

	err = handle.SetBPFFilter(s.conf.Filter)
	if err != nil {
		return err
	}

	s.handle = handle
	return nil
}

// Start to sniff. send every packet in the channel.
func (s *Sniffer) Sniff(){
    // Use the handle as a packet source to process all packets
    packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err != nil {
			fmt.Printf("ERROR!! %s\n", err.Error())
			continue
		}
		s.packetChannel <- packet
	}
}

