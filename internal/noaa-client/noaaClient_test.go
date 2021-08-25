package noaaclient

import (
	"fmt"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {

	client := NewNoaaClient(CRD, "metric")
	if client == nil {
		t.Error("empty client")
	}
}
func TestDatumEnums(t *testing.T) {

	val := MHHW.String()
	if val != "MHHW" {
		t.Error("Value didn't match")
	}
}

func TestMakeCallsForAllProducts(t *testing.T) {

	client := NewNoaaClient(CRD, "metric")
	if client == nil {
		t.Error("empty client")
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, -1, 0)
	stationID := "8454000"

	for product := DataProduct(0); product < MaximumLimit; product++ {
		val, err := client.RetreiveDataVariable(&startTime, &endTime, product, &stationID)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if val != nil {
			fmt.Println("Got Val for : " + product.String())
		} else {
			t.Error("empty val")
		}
	}
}
