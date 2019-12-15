package tracking

import (
	"time"
)

var (
	MAX_CAPACITY int = 1000
)

		
type myDB struct{
	ConnectionString string
	MyDBTableName string	
}

// 
type DBTracker struct {
	myDBTbleByTime *[][]string
	myDBTableByName *[][]string
	logger *logger.Logger
}
	
//create new DBtracker object.
func NewDBTracker(logger *logger.Logger) *DBTracker{
	return &DBTracker{
		myDBTableByName : make([][]string, 0, MAX_CAPACITY)
		myDBTbleByTime : make([][]string,0, MAX_CAPACITY)
		logger: logger
	}

}

// Initialzie the Tracker object, i.e. connecting to DB
// We use this function instead of the New* function in order to avoid errors in constructors.
// According to 'clean code' best practices we know its better to do so.
func (st *DBTracker) InitializeTracker() error {

}

// Add documentation about user's request.
// Actually this function adds an entry of the given domain that occured at specified time.
func (st *DBTracker ) RecordDomain(name string, timestamp time.Time) error {
		
	//updtae the time table - the time table arranged by time of arrival(AKA TOA)
	*myDBTableBytime = append(myDBTableBytime, timestamp)
	*myDBTableBytime[len(myDBTbleByTime-1)][2] = name
		
	//updae the name table - the name table shall be orginaized lexi...
	myDBTableByName = append(myDBTableBytime, name)
	myDBTableByName [len(myDBTableByName)-1][timestamp] 
    sort.Slice(myDBTableByName)
	}		

	// Gets recorded domains that occured more than the given amount of times.
	func (st *MyDBTableByName)GetDomainsByAmount(amount int) ([]string, error){
		
		result := map[string]int

		for index name := range myDBTableByName{
			if value > 0, := dict[name]; {
				result[name] += 1
			} else {
				result[name] = 1
			}

		}



	}

	// Get domains that their last query was before the given date.
	func (st *DBTracker)GetOldDomainsByDate(date time.Time)([]string, error){


	}

	// Delete from table the given domain.

	func (st *DBTracker) DeleteDomainByName(name string) error{
	
	}