package domainLearning

import (
	"time"
	"net"
	"sync"
	"resolver/tracking"
	"resolver/caching"
	"resolver/resolving"
	"resolver/logger"
	"resolver/utils"
)

var (
	flowName string = "domainLearning"
)

type DomainLearnerConf struct {
	Period time.Duration // time between cycles.
	MinAmount int
	CacheExpiration time.Duration
}

// UserLearner is responsoble to learn the local user behavior.
// Means, record every dns query.
type DomainLearner struct {
	tracker tracking.Tracker
	resolver *resolving.Resolver
	cache caching.Cacher
	logger *logger.Logger
	conf DomainLearnerConf
	// In order to notify the main when finish.
	wg sync.WaitGroup
	// notify injection and userLearner when starting a cycle.
	notifyStartCycle chan<- struct{}
	notifyFinishCycle chan<- struct{}
}

// Creates new DomainLearner object.
func NewDomainLearner(logger *logger.Logger, tracker tracking.Tracker, conf DomainLearnerConf, wg sync.WaitGroup, notifyStartCycle chan<- struct{}, notifyFinishCycle chan<- struct {}) *DomainLearner{
	resolver := resolving.NewResolver()
	cache := caching.NewRedisCache()
	return &DomainLearner{
		tracker : tracker,
		resolver : resolver,
		cache : cache,
		logger : logger,
		conf : conf,
		wg : wg, 
		notifyStartCycle : notifyStartCycle,
		notifyFinishCycle : notifyFinishCycle,
	}
}

func (dl *DomainLearner) InitializeDomainLearner()error{
	var err error
	
	err = dl.cache.InitializeCache()
	if err != nil{
		dl.logger.Error(flowName, "cannot initialize tracker %s", err.Error())
	}
	return err
}

// Starts the DomainLearner
// This is a blocking function. means the is noot intended to finish the job so it is recommended
// to call this function in diffrent thread (gorutine).
func (dl *DomainLearner) Learn(){
	ticker := time.NewTicker(dl.conf.Period)
	for _ = range ticker.C {
		dl.logger.Debug(flowName, "action=start_cycle")
		dl.notifyStartCycle <- struct{}{}
		domains, err := dl.tracker.GetDomainsByAmount(dl.conf.MinAmount)
		if err != nil{
			//report
			dl.logger.Error(flowName, "cannot query tracker %s", err.Error())
			continue
		}
		dl.logger.Debug(flowName, "action=got_domains, domains=%s", domains)
		dl.learnDomains(domains)
		dl.logger.Debug(flowName, "action=done_cycle")
		dl.notifyFinishCycle <- struct{}{}
	}
	dl.wg.Done()
}

// resolver each domain and updates the cache.
func (dl *DomainLearner) learnDomains(domains []string){
	for _, domain := range domains{
		ips, err := dl.resolver.Resolve(domain)
		if err != nil{
			//report
			dl.logger.Error(flowName, "cannot resolve domain %s  - %s",domain,  err.Error())
			return
		}
		if (len(ips) > 0){
			dl.updateDomainIps(domain, ips)
		}
	}
}

func (dl *DomainLearner)updateDomainIps(domain string, ips []net.IP){
	err := dl.cache.UpdateDomain(domain, ips, dl.conf.CacheExpiration)
	if err != nil {
		//report 
		dl.logger.Error(flowName, "cannot update cache for domain %s  - %s",domain,  err.Error())
		return
	}
	dl.logger.Info(flowName, "action=cache-domain domain=%s ips=%s", domain, utils.FormatIPS(ips))
}
