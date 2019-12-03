package tracking

import (
	"fmt"
	"time"
	"database/sql"
	"resolver/logger"
	_ "github.com/mattn/go-sqlite3"
)

var (
	driverName string = "sqlite3"
	dbName string = "tracking.db"
	tableName string = "tracking"
	createTableQuery string = "CREATE TABLE IF NOT EXISTS %s ('name' TEXT, 'firstQuery' TEXT,"+
								"'lastQuery' TEXT, 'amount' INTEGER, PRIMARY KEY('name'));"
	insertNewDomainQuery string = "INSERT INTO %s VALUES ('%s', '%s', '%s', 1);"
	updateDomainQuery string = "update '%s' SET amount=%d WHERE name='%s'"
	getDomainAmountQuery string = "SELECT amount FROM '%s' WHERE name='%s'"
	getDomainByAmountQuery string ="SELECT name FROM '%s' WHERE amount > %d"
)

type SQLTracker struct {
	db *sql.DB
	logger *logger.Logger
}

func NewSQLTracker(logger *logger.Logger) *SQLTracker{
	return &SQLTracker{
		logger : logger,
	}
}

func (st *SQLTracker) InitializeTracker() error {
	var err error
	var db *sql.DB
	db, err = sql.Open(driverName, fmt.Sprintf("./%s",dbName))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(createTableQuery, tableName))
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	st.db = db
	return err 
}

func (st *SQLTracker) RecordDomain(name string, timestamp time.Time) error{
	var err error
	var results *sql.Rows
	st.logger.Debug("SQLTracker", "action=start_record domain=%s",name)
	results, err = st.db.Query(fmt.Sprintf(getDomainAmountQuery, tableName, name))
	st.logger.Debug("SQLTracker", "action=got_query_results domain=%s",name)
	if err != nil {
		return err
	}
	if(results.Next() == false){
		st.logger.Debug("SQLTracker", "action=insert_new domain=%s",name)
		results.Close()
		_, err = st.db.Exec(fmt.Sprintf(insertNewDomainQuery, tableName, name, timestamp, timestamp))
		st.logger.Debug("SQLTracker", "action=finish_insert_new domain=%s",name)
		return err
	}
	st.logger.Debug("SQLTracker", "action=scan_amount domain=%s",name)
	var amount int
	err = results.Scan(&amount)
	results.Close()
	if (err != nil){
		fmt.Printf("error scan amount %s \n", name)
		return err
	}
	st.logger.Debug("SQLTracker", "action=finish_scan_amount domain=%s amount=%d",name, amount)
	st.logger.Debug("SQLTracker", "action=updating domain=%s",name)
	_, err = st.db.Exec(fmt.Sprintf(updateDomainQuery, tableName, amount+1, name))
	st.logger.Debug("SQLTracker", "action=finish_updating domain=%s",name)
	fmt.Printf("updated %s \n", name)
	return err
}

func (st *SQLTracker) GetDomainsByAmount(amount int)([]string, error){
	domains := []string{}
	results, err := st.db.Query(fmt.Sprintf(getDomainByAmountQuery, tableName, amount))
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
