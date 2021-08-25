package utils

import (
	"testing"
	"time"
)

func TestTimeUtils(t *testing.T) {
	val := time.Date(2021, time.August, 23, 0, 0, 0, 0, time.UTC)
	stringVal := ConvertTimeToyyyyMMdd(val)
	if stringVal != "20210823" {
		t.Error("Incorrect Value")
	}
}
