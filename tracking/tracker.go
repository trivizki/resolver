package tracking

import (
	"time"
)

type Tracker interface{
	// Initialzie the Tracker object, i.e. connecting to DB
	// We use this function instead of the New* function in order to avoid errors in constructors.
	// According to 'clean code' best practices we know its better to do so.
	InitializeTracker() error

	// Add documentation about user's request.
	// Actually this function adds an entry of the given domain that occured at specified time.
	RecordDomain(name string, timestamp time.Time) error

	// Gets recorded domains that occured more than the given amount of times.
	GetDomainsByAmount(amount int) ([]string, error)

	// Get domains that their last query was before the given date.
	GetOldDomainsByDate(date time.Time)([]string, error)

	// Delete from table the given domain.
	DeleteDomainByName(name string) error
}

