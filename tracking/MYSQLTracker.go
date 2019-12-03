package tracking

import (
	"fmt"
	"time"
	"database/sql"
	"resolver/logger"
	_ "github.com/go-sql-driver/mysql"
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

type MYSQLConf struct{
	ConnectionString string
	MySqlTableName string
}

// MYSQLTracker is a tracker object that use MySql Db to manage the tracking data.
type MYSQLTracker struct {
	db *sql.DB
	logger *logger.Logger
	conf MYSQLConf
}

// Creates new MYSQLTracker object.
func NewMYSQLTracker(logger *logger.Logger, conf MYSQLConf) *MYSQLTracker{
	return &MYSQLTracker{
		logger : logger,
		conf : conf,
	}
}

// Initiazlies the MYSQLTracker object. connect to the relevant MySql DB.
// NOTICE: user must call this function before using the Tracker.
func (st *MYSQLTracker) InitializeTracker() error {
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

func (st *MYSQLTracker) DeleteDomainByName(name string) error{
	_, err := st.db.Exec(fmt.Sprintf(mySqlDeleteDomainQuery, st.conf.MySqlTableName), name)
	return err
}
