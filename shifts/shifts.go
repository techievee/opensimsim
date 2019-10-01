package shift

import (
	"github.com/labstack/echo/v4"
	"github.com/techievee/opensimsim/models"
	"github.com/techievee/opensimsim/server"
	"github.com/techievee/opensimsim/storage"
	"net/http"
	"strconv"
)

type ShiftModule struct {
	server  *server.Server
	backend *storage.StorageModule
}

func NewShiftModule(s *server.Server, st *storage.StorageModule) *ShiftModule {

	return &ShiftModule{server: s, backend: st}
}

func (sm *ShiftModule) SetupRoutes() {

	sm.server.E.GET("/shift/:shift_id/swaps", sm.swaps)

}

//Result that need to be sent
type Result struct {
	Shifts *[]models.JsonShifts `json:"shifts"`
}

//To return empty array, instead of null value
type EmptyResult struct {
	Shifts []int64 `json:"shifts"`
}

// Swap Handler
func (sm *ShiftModule) swaps(c echo.Context) error {

	shiftId := c.Param("shift_id")

	//Validate param for number, Since number can be of of Bigint in DB, use ParseInt and parse with 64bit
	//Since DB value is of BigInt, 64 bit, use big int
	shiftId64, err := strconv.ParseUint(shiftId, 10, 64)
	if err != nil {
		c.Logger().Error("Invalid shift id ,", shiftId, err.Error())
		//Not a number, return notfound
		return echo.NewHTTPError(http.StatusNotFound)
	}

	//Fetch the shift ID, from the DB
	c.Logger().Debug("Fetching Shift id : ", shiftId64)
	shift, err := sm.backend.GetShifts(shiftId64)
	if err == nil && shift == nil {
		//Error is nil and shift also nil - No record in DB
		c.Logger().Debug("No Rows returned for Shift id : ", shiftId64)
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil && shift == nil {
		// Error while fetching in DB
		//Can send 500 , but to make the enduser experience better, sending 404
		c.Logger().Errorf("Error while fetchingShift id : ", shiftId64, err.Error())
		return echo.NewHTTPError(http.StatusNotFound)
	}

	//----------------------------------------------------------------------------

	//Shift ID Exist, find the possible swap for this shiftid
	c.Logger().Debug("Fetching Swaps for Shift id : ", shiftId64)

	res := &Result{}
	res.Shifts, err = sm.backend.GetSwaps(&shift.StartAt, &shift.EndAt, shift.Id)
	if err == nil && res.Shifts == nil {
		//Error is nil and shift also nil - No swaps available
		c.Logger().Debug("No Rows returned for Shift id : ", shiftId64)
		//return empty array
		//Workaround to send empty array instead of null
		shi := make([]int64, 0)
		return c.JSON(http.StatusOK, &EmptyResult{Shifts: shi})

	} else if err != nil && res.Shifts == nil {
		//Error while fetching - Error while fetching
		//Can send 500 , but to make the enduser experience better, sending 404
		c.Logger().Errorf("Error while fetching swap for shift id : ", shiftId64, err.Error())
		return echo.NewHTTPError(http.StatusNotFound)
	}

	c.Logger().Debug("Returning swaps for shift id : ", shiftId64)
	//

	return c.JSON(http.StatusOK, res)

}
