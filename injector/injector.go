package injector 

import (
		"time"
		"github.com/google/gopacket/pcap"
)


var (
		snapshotLen int32  = 1024
		promiscuous bool   = false
		timeout     time.Duration = 30 * time.Second
)

type InjectorConf struct {
	Device string
}

// Injector is responsible to inject packet to the configured interface.
type Injector struct {
	handle        *pcap.Handle
	conf InjectorConf
}

func NewInjector(conf InjectorConf) (*Injector) {
	return &Injector{
		conf : conf,
	}
}

// Initializes the injector object
// NOTICE: User must call this function beforse any usage.
func (s *Injector) InitializeInjector() error{

	// Open device
	handle, err := pcap.OpenLive(s.conf.Device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return err
	}

	s.handle = handle
	return nil
}

// Injects the given packet.
func (s *Injector) Inject(packet []byte) error {
	return s.handle.WritePacketData(packet)
}

