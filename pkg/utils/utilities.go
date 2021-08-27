package utils

import (
	"encoding/json"
	"time"
)

//ConvertTimeToyyyyMMdd will converts to a valid string
//
//Note:  I chose not to pass in a pointer here so that the utility function don't have to handle nil values
func ConvertTimeToyyyyMMdd(val time.Time) string {
	return val.Format("20060102")
}

//MarshalDataToInterface - Utility function to marshall data to struct.  If you don't pass in a pointer then no value will be returned
func MarshalDataToInterface(data []byte, pointerOfStructToMarshalInto interface{}) error {

	err := json.Unmarshal(data, pointerOfStructToMarshalInto)
	if err != nil {
		return err
	}
	return nil
}
