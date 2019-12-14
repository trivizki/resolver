package injection

import (
	"net"
	"sync"
	"resolver/sniffer"
	"resolver/caching"
	"resolver/packetBuilder"
	"resolver/injector"
	"resolver/context"
	"resolver/logger"
	"resolver/utils"
	"github.com/google/gopacket"
)

var (
	flowName string = "injection"
)

type InjectionConf struct {
	SnifferConf sniffer.SnifferConf
	InjectorConf injector.InjectorConf
}

// Injection is the flow responsible to inject DNS reponses to the user.
type Injection struct {
	sniffer *sniffer.Sniffer
	contextBuilder *context.ContextBuilder
	cache caching.Cacher
	injector *injector.Injector
	packetBuilder *packetBuilder.PacketBuilder
	packetChannel chan gopacket.Packet
	logger *logger.Logger
	wg sync.WaitGroup
	notifyStartCycle <-chan struct{}
	notifyFinishCycle <-chan struct{}
}

func NewInjection(logger *logger.Logger, conf InjectionConf, wg sync.WaitGroup,
			notifyStartCycle <-chan struct{}, notifyFinishCycle <-chan struct {}) *Injection{
	packetChannel := make(chan gopacket.Packet, 100)
	sniffer := sniffer.NewSniffer(packetChannel, conf.SnifferConf, logger)
	contextBuilder := context.NewContextBuilder()
	cacher := caching.NewRedisCache()
	injector := injector.NewInjector(conf.InjectorConf)
	packetBuilder := packetBuilder.NewPacketBuilder()
	return &Injection{
		sniffer : sniffer,
		contextBuilder : contextBuilder,
		cache : cacher,
		injector : injector,
		packetBuilder : packetBuilder,
		packetChannel : packetChannel,
		logger : logger,
		wg : wg,
		notifyStartCycle : notifyStartCycle,
		notifyFinishCycle : notifyFinishCycle,
	}
}

func (i *Injection) InitializeInjection() error{
	var err error
	err = i.sniffer.InitializeSniffer()
	if err != nil{
		i.logger.Error(flowName, "cannot initialize sniffer %s",err.Error())
		return err
	}

	err = i.cache.InitializeCache()
	if err != nil{
		i.logger.Error(flowName, "cannot initialize cache%s",err.Error())
		return err
	}

	err = i.injector.InitializeInjector()
	if err != nil{
		i.logger.Error(flowName, "cannot initialize injector%s",err.Error())
		return err
	}
	return err 
}

// Start the flow
// This is a blocking function. means the is noot intended to finish the job so it is recommended
// to call this function in diffrent thread (gorutine).
func (i *Injection) Inject(){
	go i.sniffer.Sniff()
	i.listen()	
	i.wg.Done()
}

func (i *Injection) listen(){
	for {
		select {
		case packet := <- i.packetChannel:
				go i.handlePacket(packet)
			case <-i.notifyStartCycle:
				i.logger.Debug(flowName, "action=pause")
				i.pause()
		}
	}
}

func (i *Injection) pause(){
	for {
		select {
			case <- i.packetChannel:
				continue
			case <- i.notifyFinishCycle:
				i.logger.Debug(flowName, "action=resume")
				i.listen()
		}
	}
}

// handle one sniffed packet.
// check whether we know the domain's ip, if we know - inject response.
func (i *Injection) handlePacket(packet gopacket.Packet){
		var err error
		var context *context.Context
		var ips []net.IP
		var response []byte
		context, err = i.contextBuilder.BuildContext(packet)
		if err != nil{
			i.logger.Error(flowName, "cannot build context %s",err.Error())
			return
		} 
		if !(context.IsAQuery()) {
			return
		}
		domain := string(context.Payload.Questions[0].Name)
		ips, err = i.cache.GetIPS(domain)
		if err != nil {
			i.logger.Error(flowName, "cannot query cache %s",err.Error())
		}
		if (len(ips) == 0){
			return
		}
		response, err = i.buildInjectPacket(*context, ips)
		if err != nil {
			i.logger.Error(flowName, "cannot build injection packet %s",err.Error())
		}
		err = i.injector.Inject(response)
		if err != nil {
			i.logger.Error(flowName, "cannot inject packet %s",err.Error())
		}
		i.logger.Info(flowName, "action=injection domain=%s ips=%s ",domain, utils.FormatIPS(ips))
}

// Builds the dns response packet.
func (i Injection) buildInjectPacket(context context.Context, answers []net.IP) ([]byte, error){
	return i.packetBuilder.BuildPacket(context.DstMAC, context.SrcMAC, context.DstIP,
	context.SrcIP, context.DstPort, context.SrcPort, context.Payload, answers)
}
