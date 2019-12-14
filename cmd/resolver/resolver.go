package main 

import (
	"fmt"
	"sync"
	"resolver/flows/domainLearning"
	"resolver/flows/userLearning"
	"resolver/flows/maintanance"
	"resolver/flows/injection"
	"resolver/tracking"
	"resolver/logger"
	"github.com/spf13/viper"
)

var (
	confFile string = "conf"
	confPath string = "."
)

// The resolver main configuration struct.
type Configuration struct {
	LoggerConf logger.LoggerConf
	MySqlConf tracking.MYSQLConf
	DomainLearnerConf domainLearning.DomainLearnerConf
	InjectionConf injection.InjectionConf
	MaintananceConf maintanance.MaintananceConf
	UserLearnerConf userLearning.UserLearnerConf
}

// This function is responsible to load the system configuration.
func readConfiguration() (Configuration, error){
	var conf Configuration
	var err error
	viper.SetConfigName(confFile)
	viper.AddConfigPath(confPath)
	err = viper.ReadInConfig()
	if err != nil {
		return conf, err
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}
	return conf, err
}


func main(){
	var conf Configuration
	var err error
	var wg sync.WaitGroup
	learningStartCycle := make(chan struct{})
	learningFinishCycle := make(chan struct{})

	conf, err = readConfiguration()
	if err != nil{
		fmt.Printf("Cannot load configuration %s\n",err.Error())
		return
	}
	logger, err := logger.NewLogger(conf.LoggerConf)
	if err != nil{
		fmt.Printf("Cannot load logger %s\n",err.Error())
		return
	}
	tracker := tracking.NewMYSQLTracker(logger, conf.MySqlConf)
	err = tracker.InitializeTracker()
	if err != nil{
		fmt.Printf("Cannot initialize tracker %s\n",err.Error())
		return 
	}

	wg.Add(4)
	learner := userLearning.NewUserLearner(logger, tracker, conf.UserLearnerConf, wg)
	domainLearner := domainLearning.NewDomainLearner(logger, tracker, conf.DomainLearnerConf, wg, 
		learningStartCycle, learningFinishCycle)
	injection := injection.NewInjection(logger, conf.InjectionConf, wg,
		learningStartCycle, learningFinishCycle)
	maintanance := maintanance.NewMaintanance(logger, tracker, conf.MaintananceConf, wg)

	err = learner.InitializeLearner()
	if err != nil{
		fmt.Printf("Cannot initialize user learner %s\n",err.Error())
		return
	}
	err = domainLearner.InitializeDomainLearner()
	if err != nil{
		fmt.Printf("Cannot initialize domain learner %s\n",err.Error())
		return
	}
	err = injection.InitializeInjection() 
	if err != nil{
		fmt.Printf("Cannot initialize injection %s\n",err.Error())
		return
	}
	fmt.Println("Start!")

	go learner.Learn()
	go domainLearner.Learn()
	go maintanance.Maintan()
	go injection.Inject()

	// we have to keep the main thread alive. 
	wg.Wait()
}
