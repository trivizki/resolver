package maintanance

import (
	"time"
	"sync"
	"resolver/tracking"
	"resolver/logger"
)

var (
	flowName string = "maintanance"
)

type MaintananceConf struct {
	Period time.Duration
	// The duration past since last query.
	// Every domain name that was queried before will be deleted.
	MaxLastQuery time.Duration
}

// Maintanance is responsible to maintan the tracking database in order to avoid overload.
type Maintanance struct {
	tracker tracking.Tracker
	logger *logger.Logger
	conf MaintananceConf
	wg sync.WaitGroup
}

// Creates new maintanance object.
func NewMaintanance(logger *logger.Logger, tracker tracking.Tracker, conf MaintananceConf, wg sync.WaitGroup) *Maintanance{
	return &Maintanance{
		tracker : tracker,
		logger : logger,
		conf : conf,
		wg : wg,
	}
}

// Start the maintance flow
// This is a blocking function. means the is noot intended to finish the job so it is recommended
// to call this function in diffrent thread (gorutine).
func (dl *Maintanance) Maintan(){
	ticker := time.NewTicker(dl.conf.Period)
	for _ = range ticker.C {
		dl.logger.Debug(flowName, "action=start_cycle")
		startingDate := time.Now().Add(dl.conf.MaxLastQuery)
		domains, err := dl.tracker.GetOldDomainsByDate(startingDate)
		if err != nil{
			dl.logger.Error(flowName, "cannot query tracker %s", err.Error())
			continue
		}
		dl.logger.Debug(flowName, "action=got_domains, domains=%s", domains)
		dl.deleteDomains(domains)
		dl.logger.Debug(flowName, "action=done_cycle")
	}
	dl.wg.Done()
}

// Delete the given domain from DB.
func (dl *Maintanance) deleteDomains(domains []string){
	var err error
	for _, domain := range domains{
		err = dl.tracker.DeleteDomainByName(domain)
		if (err != nil){
			dl.logger.Error(flowName, "Error while delete domain %s. %s", domain, err.Error())
		} else {
			dl.logger.Info(flowName, "action=delete domain=%s",domain)
		}
	}
}
