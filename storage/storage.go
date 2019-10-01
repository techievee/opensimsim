package storage

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/techievee/opensimsim/config"
	"github.com/techievee/opensimsim/models"
	"time"
)

type StorageModule struct {
	db  *sql.DB
	cfg *config.MySql
}

func NewStorageModule(cfg *config.MySql) *StorageModule {

	dbDriver := cfg.DbDriver
	dbUser := cfg.DbUser
	dbPass := cfg.DbPass
	dbName := cfg.DbName
	dbHost := cfg.DbHost
	dbPort := cfg.DbPort

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return &StorageModule{db: db, cfg: cfg}
}

func (s *StorageModule) Check() error {
	err := s.db.Ping()
	if err != nil {
		return err // proper error handling instead of panic in your app
	} else {
		return nil
	}
}

func (s *StorageModule) Close() {
	s.db.Close()
	return
}

// All the DB relation function for Swap Module
//Function queries the DB to find the shiftid match for the worker
//Return Shift, error
//If there is norows, it returns nil, nil
//If error during the Query, it returns nil, error
//If success fetch, it return shift, error
func (s *StorageModule) GetShifts(shift_id uint64) (*models.Shifts, error) {

	var id uint64
	var worker string
	var start_at mysql.NullTime
	var end_at mysql.NullTime

	stmtOut, err := s.db.Prepare("SELECT id,worker,start_at,end_at  FROM shifts s WHERE s.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(shift_id).Scan(&id, &worker, &start_at, &end_at)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		sh := &models.Shifts{id, worker, start_at.Time, end_at.Time}
		fmt.Println(worker)
		return sh, err
	default:
		return nil, err
	}

}

// All the DB relation function for Swap Module
//Function queries the DB to find the shiftid match for the worker
//Return Shift, error
//If there is norows, it returns nil, nil
//If error during the Query, it returns nil, error
//If success fetch, it return shift, error
func (s *StorageModule) GetSwaps(startdate *time.Time, enddate *time.Time, id uint64) (*[]models.JsonShifts, error) {

	var shifts []models.JsonShifts

	stmtOut, err := s.db.Prepare("SELECT id,worker,start_at,end_at  FROM shifts s WHERE (s.start_at >= ? or s.end_at<= ?)  and s.id<>?")
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(enddate, startdate, id)
	defer rows.Close()

	count := 0

	for rows.Next() {
		var id uint64
		var worker string
		var start_at mysql.NullTime
		var end_at mysql.NullTime

		count += 1

		err = rows.Scan(&id, &worker, &start_at, &end_at)
		if err != nil {
			return nil, err
		}

		//Convert the RFC3339 to without T and Z for the output
		layout := "2006-01-02 15:04:05"
		iso8601st := start_at.Time.Format(layout)
		iso8601et := end_at.Time.Format(layout)

		sh := models.JsonShifts{id, worker, iso8601st, iso8601et}
		shifts = append(shifts, sh)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		//Send error and nil
		return nil, err
	}

	//If no row exist, then return nil for both return values
	if count == 0 {

		return nil, nil
	}

	return &shifts, nil

}
