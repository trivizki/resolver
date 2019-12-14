package userLearning

import (
	"time"
	"sync"
	"github.com/google/gopacket"
	"resolver/sniffer"
	"resolver/context"
	"resolver/tracking"
	"resolver/logger"
)

var (
	flowName string = "userLearner"
)

type UserLearnerConf struct {
	SnifferConf sniffer.SnifferConf
}

// UserLearner is responsoble to learn the local user behavior.
// Means, record every dns query.
type UserLearner struct {
	packetChannel chan gopacket.Packet
	sniffer *sniffer.Sniffer
	contextBuilder *context.ContextBuilder
	tracker tracking.Tracker
	logger *logger.Logger
	wg sync.WaitGroup
	notifyStartCycle <-chan struct{}
	notifyFinishCycle <-chan struct{}
}

// create new user learner.
// NOTICE: One has to call InitializeLearner function before start using an UserLearner object.
func NewUserLearner(logger *logger.Logger, tracker tracking.Tracker, conf UserLearnerConf, wg sync.WaitGroup, notifyStartCycle <-chan struct{}, notifyFinishCycle <-chan struct {}) *UserLearner{
	packetChannel := make(chan gopacket.Packet, 30)
	sniffer := sniffer.NewSniffer(packetChannel, conf.SnifferConf, logger)
	contextBuilder := context.NewContextBuilder()
	return &UserLearner{
		packetChannel : packetChannel,
		sniffer : sniffer,
		contextBuilder : contextBuilder,
		tracker : tracker,
		logger : logger,
		wg : wg,
		notifyStartCycle : notifyStartCycle,
		notifyFinishCycle : notifyFinishCycle,
	}
}

// InitializeLearner.
func (ul *UserLearner) InitializeLearner() error{
	var err error
	err = ul.sniffer.InitializeSniffer()
	if err != nil{
		ul.logger.Error(flowName, "cannot initialize sniffer %s", err.Error())
		return err
	}
	return err 
}

// Start recording dns queries.
// This is a blocking function. means the is noot intended to finish the job so it is recommended
// to call this function in diffrent thread (gorutine).
func (ul *UserLearner) Learn(){
	go ul.sniffer.Sniff()
	ul.listen()
	ul.wg.Done()
}

func (ul *UserLearner) listen(){
	for {
		select {
		case packet := <- ul.packetChannel:
				go ul.handlePacket(packet)
			case <- ul.notifyStartCycle:
				ul.logger.Debug(flowName, "action=pause")
				ul.pause()
		}
	}

}

func (ul *UserLearner) pause(){
	for {
		select {
			case <- ul.packetChannel:
				continue
			case <- ul.notifyFinishCycle:
				ul.logger.Debug(flowName, "action=resume")
				ul.listen()
		}
	}
}



// handle one sniffed packet.
// buid context from packet, recording the query.
func (ul *UserLearner)handlePacket(packet gopacket.Packet){
		start := time.Now()
		var err error
		var context *context.Context
		context, err = ul.contextBuilder.BuildContext(packet)
		if err != nil{
			ul.logger.Error(flowName, "cannot build context %s", err.Error())
			return
		} 
		domain := string(context.Payload.Questions[0].Name)
		if !(context.IsQuery()) {
			ul.logger.Debug(flowName, "action=no_query domain=%s", domain)
			return
		}
		err = ul.tracker.RecordDomain(domain, context.Timestamp)
		if err != nil {
			ul.logger.Error(flowName, "cannot record domain %s %s",domain, err.Error())
			return
		}
		t := time.Now()
		ul.logger.Info(flowName, "action=record domain=%s time=%s", domain, t.Sub(start))
}
