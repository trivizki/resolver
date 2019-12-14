package tracking

import (
	"time"
)

var (

)

		
type myDB struct{
	ConnectionString string
	MyDBTableName string	
}

// 
type DBTracker struct {
	*MyDBTbleByTime [][2]string
	*MyDBTableByName [][2]string
	logger *logger.Logger
}
	
//create new DBtracker object.
func NewDBTracker(logger *logger.Logger) *DBTracker{
	*MyDBTableByName = make([][2]string, )
	*MyDBTableByTime = make([][2]string, )
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
	*MyDBTableBytime = append(MyDBTableByName, {timestamp, name})
		
	//updae the name table - the name table shall be orginaized lexi...
	MyDBTableByName

	}		

	// Gets recorded domains that occured more than the given amount of times.
	func (st *MyDBTableByName)GetDomainsByAmount(amount int) ([]string, error){
		
		var result string[]

	}

	// Get domains that their last query was before the given date.
	GetOldDomainsByDate(date time.Time)([]string, error)

	// Delete from table the given domain.

	func (st *MYSQLTracker) DeleteDomainByName(name string) error{
		func SearchStrings(a []string, name string) int
}

// Initiazlies the MYSQLTracker object. connect to the relevant MySql DB.
// NOTICE: user must call this function before using the Tracker.

	
func (st *MYSQLTracker) GetOldDomainsByDate(date time.Time)([]string, error){
		


	}