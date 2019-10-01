package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"testing"
	"time"
)

var r *StorageModule

func TestMain(m *testing.M) {
	r = NewStorageModule()
	c := m.Run()
	os.Exit(c)
}

func TestDBConn(t *testing.T) {

	err := r.Check()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	} else {
		t.Logf("DB COnnection succeeded")
	}

}

func TestGetShifts(t *testing.T) {

	shift, err := r.GetShifts(uint64(86))
	if err == nil && shift == nil {
		t.Logf("No Rows")
	} else if err != nil && shift == nil {
		t.Errorf(err.Error())
	} else {
		t.Logf(shift.Worker, shift.Id, shift.StartAt, shift.EndAt)
	}

}

func TestGetSwaps(t *testing.T) {

	var Id uint64
	var start_at time.Time
	var end_at time.Time
	var err error

	Id = 96
	/*rfc3339st := strings.Replace("2019-08-01 00:30:00", " ", "T", 1) + "Z"
	rfc3339et := strings.Replace("2019-08-01 01:30:00", " ", "T", 1) + "Z"*/
	rfc3339st := strings.Replace("2015-08-01 00:30:00", " ", "T", 1) + "Z"
	rfc3339et := strings.Replace("2021-08-01 01:30:00", " ", "T", 1) + "Z"

	start_at, err = time.Parse(time.RFC3339, rfc3339st)
	if err != nil {
		//t.Errorf(err.Error())
	}
	end_at, _ = time.Parse(time.RFC3339, rfc3339et)
	shifts, err := r.GetSwaps(&start_at, &end_at, Id)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf("%v", shifts)

}
