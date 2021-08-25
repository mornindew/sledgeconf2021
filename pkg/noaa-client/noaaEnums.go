package noaaclient

import (
	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
)

//DataProducts - enum for the specific data sets that Noaa Offers
type DataProduct int

const (
	WaterLevel DataProduct = iota
	AirTemperature
	WaterTemperature
	Wind
	AirPressure
	AirGap
	Conductivity
	Visibility
	Humidity
	Salinity
	HourlyHeight
	HighLow
	DailyMean
	MonthlyMean
	OneMinuteWaterLevel
	Preditions
	Datums
	Currents
	CurrentsPredictions
	MaximumLimit //Must be the last one - makes it very easy to loop though
)

func (enum DataProduct) String() string {
	switch enum {
	case WaterLevel:
		return "water_level"
	case AirTemperature:
		return "air_temperature"
	case WaterTemperature:
		return "water_temperature"
	case Wind:
		return "wind"
	case AirPressure:
		return "air_pressure"
	case AirGap:
		return "air_gap"
	case Conductivity:
		return "conductivity"
	case Visibility:
		return "visibility"
	case Humidity:
		return "humidity"
	case Salinity:
		return "salinity"
	case HourlyHeight:
		return "hourly_height"
	case HighLow:
		return "high_low"
	case DailyMean:
		return "daily_mean"
	case MonthlyMean:
		return "monthly_mean"
	case OneMinuteWaterLevel:
		return "one_minute_water_level"
	case Preditions:
		return "predictions"
	case Datums:
		return "datums"
	case Currents:
		return "currents"
	case CurrentsPredictions:
		return "currents_predictions"
	default:
		//I could/should throw an error here, have a default, or at least log something
		return ""
	}
}

func (enum DataProduct) ConvertToGrpcEnum() sledgconf_demo_proto_v1.DataType {
	switch enum {
	case WaterLevel:
		return sledgconf_demo_proto_v1.DataType_WaterLevel
	case AirTemperature:
		return sledgconf_demo_proto_v1.DataType_AirTemperature
	case WaterTemperature:
		return sledgconf_demo_proto_v1.DataType_WaterTemperature
	case Wind:
		return sledgconf_demo_proto_v1.DataType_Wind
	case AirPressure:
		return sledgconf_demo_proto_v1.DataType_AirPressure
	case AirGap:
		return sledgconf_demo_proto_v1.DataType_AirGap
	case Conductivity:
		return sledgconf_demo_proto_v1.DataType_Conductivity
	case Visibility:
		return sledgconf_demo_proto_v1.DataType_Visibility
	case Humidity:
		return sledgconf_demo_proto_v1.DataType_Humidity
	case Salinity:
		return sledgconf_demo_proto_v1.DataType_Salinity
	case HourlyHeight:
		return sledgconf_demo_proto_v1.DataType_HourlyHeight
	case HighLow:
		return sledgconf_demo_proto_v1.DataType_HighLow
	case DailyMean:
		return sledgconf_demo_proto_v1.DataType_DailyMean
	case MonthlyMean:
		return sledgconf_demo_proto_v1.DataType_MonthlyMean
	case OneMinuteWaterLevel:
		return sledgconf_demo_proto_v1.DataType_OneMinuteWaterLevel
	case Preditions:
		return sledgconf_demo_proto_v1.DataType_Preditions
	case Datums:
		return sledgconf_demo_proto_v1.DataType_Datums
	case Currents:
		return sledgconf_demo_proto_v1.DataType_Currents
	case CurrentsPredictions:
		return sledgconf_demo_proto_v1.DataType_CurrentsPredictions
	default:
		//I could/should throw an error here, have a default, or at least log something
		return sledgconf_demo_proto_v1.DataType_WaterLevel
	}
}

//Datum - enum for the specific datums for Noaa
type Datum int

const (
	CRD Datum = iota
	IGLD
	LWD
	MHHW
	MHW
	MTL
	MSL
	MLW
	MLLW
	NAVD
	STND
)

//showing a different way to write the string function
func (enum Datum) String() string {
	return []string{"CRD", "IGLD", "LWD", "MHHW", "MHW", "MTL", "MSL", "MLW", "MLLW", "NAVD", "STND"}[enum]
}

//Helper Function to convert from string to datum
func ConvertStringDatumToEnum(val string) (Datum, error) {
	switch val {
	case "CRD":
		return CRD, nil
	case "IGLD":
		return IGLD, nil
	case "LWD":
		return LWD, nil
	case "MHHW":
		return MHHW, nil
	case "MHW":
		return MHW, nil
	case "MTL":
		return MTL, nil
	case "MSL":
		return MSL, nil
	case "MLW":
		return MLW, nil
	case "MLLW":
		return MLLW, nil
	case "NAVD":
		return NAVD, nil
	case "STND":
		return STND, nil
	}
	//Handle something not matching
	return -1, customerrors.InvalidData{Msg: "Not a valid datum"}
}

//MeasurementUnit - enum for the measurement unit
type MeasurementUnit int

const (
	English MeasurementUnit = iota
	Metric
)

//showing a different way to write the string function
func (enum MeasurementUnit) String() string {
	return []string{"english", "metric"}[enum]
}
