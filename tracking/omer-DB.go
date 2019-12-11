package tracking

import (
	"time"
)

var (
	MySqlDriverName string = "mysql"
	mySqlCreateTableQuery string = "CREATE TABLE IF NOT exists `resolver`.`%s`"+
				" (`name` VARCHAR(50) NOT NULL, `amount` INT NULL, `first_query` DATETIME NULL," +
				"`last_query` DATETIME NULL,  PRIMARY KEY (`name`)," +
				"UNIQUE INDEX `name_UNIQUE` (`name` ASC));"
	mySqlInsertNewDomainQuery string = "INSERT INTO `resolver`.`%s` VALUES (?, 1, ?, ?)"
	mySqlUpdateDomainQuery string = "UPDATE `resolver`.`%s` SET amount = ?, last_query = ? WHERE name = ?"
	mySqlGetDomainAmountQuery string = "SELECT amount FROM `resolver`.`%s` WHERE name=?"
	mySqlGetDomainByAmountQuery string ="SELECT name FROM `resolver`.`%s` WHERE amount > ?"
	mySqlGetOldDomainByDateQuery string ="SELECT name FROM `resolver`.`%s` WHERE last_query < ?"
	mySqlDeleteDomainQuery string ="DELETE FROM `resolver`.`%s` WHERE name=?"
)

type Tracker interface{
		
	type myDB struct{
		ConnectionString string
		MyDBTableName string
	}//close myDB

	// MYSQLTracker is a tracker object that use MySql Db to manage the tracking data.
	type DBTracker struct {
		db *sql.DB
		logger *logger.Logger
		conf MYSQLConf
	}
	
	//create new DBtracker object.
	func NewDBTracker(logger *logger.Logger, conf MYSQLConf) *MYSQLTracker{
		return &MYSQLTracker{
			logger : logger,
			conf : conf,
		}
	}

	// Initialzie the Tracker object, i.e. connecting to DB
	// We use this function instead of the New* function in order to avoid errors in constructors.
	// According to 'clean code' best practices we know its better to do so.
	func (st *DBTracker) InitializeTracker() error {
		var err error
		var db *sql.DB
		db, err = sql.Open(MySqlDriverName, st.conf.ConnectionString)
		if err != nil {
			return err
		}
		_, err = db.Exec(fmt.Sprintf(mySqlCreateTableQuery, st.conf.MySqlTableName))
		if err != nil {
			return err
		}
		// In order to prevent Database overload.
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(time.Minute)
		st.db = db
		return err 
	}

	// Add documentation about user's request.
	// Actually this function adds an entry of the given domain that occured at specified time.
	func RecordDomain(name string, timestamp time.Time) error {
		
	
		func (st *MYSQLTracker) RecordDomain(name string, timestamp time.Time) error{
			var err error
			var results *sql.Rows
			results, err = st.db.Query(fmt.Sprintf(mySqlGetDomainAmountQuery, st.conf.MySqlTableName), name)
			if err != nil {
				results.Close()
				return err
			}
			if(results.Next() == false){
				results.Close()
				_, err = st.db.Exec(fmt.Sprintf(mySqlInsertNewDomainQuery,st.conf.MySqlTableName), name, timestamp, timestamp)
				return err
			}
			err = results.Err()
			if err != nil {
				results.Close()
				return err
			}
			var amount int
			err = results.Scan(&amount)
			if (err != nil){
				results.Close()
				return err
			}
			results.Close()
			_, err = st.db.Exec(fmt.Sprintf(mySqlUpdateDomainQuery, st.conf.MySqlTableName), amount+1, timestamp, name)
			return err
		}	

	}

	// Gets recorded domains that occured more than the given amount of times.
	GetDomainsByAmount(amount int) ([]string, error)

	// Get domains that their last query was before the given date.
	GetOldDomainsByDate(date time.Time)([]string, error)

	// Delete from table the given domain.
	DeleteDomainByName(name string) error


	func (st *MYSQLTracker) DeleteDomainByName(name string) error{
		_, err := st.db.Exec(fmt.Sprintf(mySqlDeleteDomainQuery, st.conf.MySqlTableName), name)
		return err
	}


	// Initiazlies the MYSQLTracker object. connect to the relevant MySql DB.
	// NOTICE: user must call this function before using the Tracker.

	
	func (st *MYSQLTracker) GetDomainsByAmount(amount int)([]string, error){
		domains := []string{}
		results, err := st.db.Query(fmt.Sprintf(mySqlGetDomainByAmountQuery, st.conf.MySqlTableName), amount)
		if (err != nil){
			return domains, err
		}
		defer results.Close()
		for results.Next(){
			var domain string
			results.Scan(&domain)
			domains = append(domains, domain)
		}
		return domains, err
	}
	
	func (st *MYSQLTracker) GetOldDomainsByDate(date time.Time)([]string, error){
		domains := []string{}
		results, err := st.db.Query(fmt.Sprintf(mySqlGetOldDomainByDateQuery, st.conf.MySqlTableName), date)
		if (err != nil){
			return domains, err
		}
		defer results.Close()
		for results.Next(){
			var domain string
			results.Scan(&domain)
			domains = append(domains, domain)
		}
		return domains, err
	}


}// close tracker