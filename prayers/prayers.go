package prayers

import (
	"fmt"
	"math"

	"github.com/CapregSoft/go-prayer-time/commons"
	"github.com/CapregSoft/go-prayer-time/constants"
)

type Prayer struct {
	MethodParams   map[int][]float64
	CalcMethod     int // caculation method
	AsrJuristic    int // Juristic method for Asr
	DhuhrMinutes   int // minutes after mid-day for Dhuhr
	AdjustHighLats int // adjusting method for higher latitudes

	TimeFormat int // time format

	Lat      float64 // latitude
	Lng      float64 // longitude
	TimeZone int     // time-zone
	JDate    float64 // Julian date

	//Times   []int
	Offsets []int
	// Time Names
	TimeName []string
	//--------------------- Technical Settings --------------------

	NumIterations int // number of iterations needed to compute times

	//------------------- Calc Method Parameters --------------------

}

type PrayerData struct {
	Fajar   string `json:"fajar"`
	Sunrise string `json:"sunrise"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Sunset  string `json:"sunset"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

func (p *Prayer) Init() {

	p.CalcMethod = 3
	//p.AsrJuristic = 0
	p.DhuhrMinutes = 0
	p.AdjustHighLats = 1
	p.TimeFormat = 0
	p.NumIterations = 1
	//p.Times = make([]int, 7)
	p.Offsets = []int{0, 0, 0, 0, 0, 0, 0}
	p.TimeName = []string{
		"Fajr",
		"Sunrise",
		"Dhuhr",
		"Asr",
		"Sunset",
		"Maghrib",
		"Isha",
	}
	p.MethodParams = map[int][]float64{
		constants.JAFARI:  {16, 0, 4, 0, 14},
		constants.KARACHI: {18, 1, 0, 0, 18},
		constants.ISNA:    {15, 1, 0, 0, 15},
		constants.MWL:     {18, 1, 0, 0, 17},
		constants.MAKKAH:  {18.5, 1, 0, 1, 90},
		constants.EGYPT:   {19.5, 1, 0, 0, 17.5},
		constants.TEHRAN:  {17.7, 0, 4.5, 0, 14},
		constants.CUSTOM:  {18, 1, 0, 0, 17},
	}

}

// compute mid-day (Dhuhr, Zawal) time
func (p *Prayer) computeMidDay(t float64) float64 {
	T := commons.EquationOfTime(p.JDate + t)
	Z := commons.FixHour(12 - T)
	return Z
}

// compute time for a given angle G
func (p *Prayer) computeTime(G float64, t float64) float64 {

	D := commons.SunDeclination(p.JDate + t)
	Z := p.computeMidDay(t)
	B := -commons.Dsin(G) - commons.Dsin(D)*commons.Dsin(p.Lat)
	M := commons.Dcos(D) * commons.Dcos(p.Lat)
	V := commons.Darccos(B/M) / 15.0

	if G > 90 {
		Z = Z + (-V)
	} else {
		Z = Z + V
	}
	return Z
}

// compute the time of Asr
// Shafii: step=1, Hanafi: step=2
func (p *Prayer) computeAsr(step int, t float64) float64 { // Shafii: step=1, Hanafi: step=2

	D := commons.SunDeclination(p.JDate + t)
	G := -commons.Darccot(float64(step) + commons.Dtan(math.Abs(p.Lat-D)))
	return p.computeTime(G, t)
}

// -------------------- Interface Functions --------------------
// return prayer times for a given date

func (p *Prayer) getDatePrayerTimes(year int, month int, day int, latitude float64, longitude float64, tZone int) []string {
	p.Lat = latitude
	p.Lng = longitude
	p.TimeZone = tZone
	p.JDate = commons.JulianDate(year, month, day) - (longitude / (15.0 * 24.0))
	return p.computeDayTimes()
}
func (p *Prayer) GetPrayerTimes(year int, month int, day int, latitude float64, longitude float64, tZone int) []string {

	return p.getDatePrayerTimes(year, month, day, latitude, longitude, tZone)
}

func (p *Prayer) GetPrayerTimesAsObject(year int, month int, day int, latitude float64, longitude float64, tZone int) PrayerData {

	pray := p.getDatePrayerTimes(year, month, day, latitude, longitude, tZone)
	return PrayerData{
		Fajar:   pray[0],
		Sunrise: pray[1],
		Dhuhr:   pray[2],
		Asr:     pray[3],
		Sunset:  pray[4],
		Maghrib: pray[5],
		Isha:    pray[6],
	}
}

// set custom values for calculation parameters
func (p *Prayer) setCustomParams(params []float64) {

	for i := 0; i < 5; i++ {
		if params[i] == -1 {
			params[i] = p.MethodParams[p.CalcMethod][i]
			p.MethodParams[constants.CUSTOM] = params
		} else {
			p.MethodParams[constants.CUSTOM][i] = params[i]
		}
	}
	p.CalcMethod = constants.CUSTOM
}

// compute prayer times at given julian date
func (p *Prayer) computeTimes(times []float64) []float64 {

	t := p.dayPortion(times)
	Fajr := p.computeTime(180-p.MethodParams[p.CalcMethod][0], t[0])
	Sunrise := p.computeTime(180-0.833, t[1])
	Dhuhr := p.computeMidDay(t[2])
	Asr := p.computeAsr(1+p.AsrJuristic, t[3])
	Sunset := p.computeTime(0.833, t[4])
	Maghrib := p.computeTime(p.MethodParams[p.CalcMethod][2], t[5])
	Isha := p.computeTime(p.MethodParams[p.CalcMethod][4], t[6])

	return []float64{Fajr, Sunrise, Dhuhr, Asr, Sunset, Maghrib, Isha}
}

func (p *Prayer) computeDayTimes() []string {
	times := []float64{5, 6, 12, 13, 18, 18, 18} //default times

	for i := 1; i <= p.NumIterations; i++ {
		times = p.computeTimes(times)
	}
	times = p.adjustTimes(times)
	return p.adjustTimesFormat(times)

}
func (p *Prayer) adjustTimes(times []float64) []float64 {
	for i := 0; i < len(times); i++ {
		times[i] += float64(p.TimeZone) - p.Lng/15.0
	}
	times[2] += float64(p.DhuhrMinutes) / 60.0 //Dhuhr
	if p.MethodParams[p.CalcMethod][1] == 1 {  // Maghrib
		times[5] = times[4] + p.MethodParams[p.CalcMethod][2]/60.0
	}
	if p.MethodParams[p.CalcMethod][3] == 1 { // Isha
		times[6] = times[5] + p.MethodParams[p.CalcMethod][4]/60.0
	}
	if p.AdjustHighLats != constants.NONE {
		times = p.adjustHighLatTimes(times)
	}

	return times
}
func (p *Prayer) adjustTimesFormat(times []float64) []string {
	formatted := make([]string, len(times))
	if p.TimeFormat == constants.FLOATING {
		for i := 0; i < len(times); i++ {
			formatted[i] = fmt.Sprintf("%f", times[i])
		}
		return formatted
	}
	for i := 0; i < 7; i++ {
		if p.TimeFormat == constants.TIME_12 {
			formatted[i] = commons.FloatToTime12(times[i], false)
		} else if p.TimeFormat == constants.TIME_12_NS {
			formatted[i] = commons.FloatToTime12NS(times[i])
		} else {
			formatted[i] = commons.FloatToTime24(times[i])
		}
	}
	return formatted
}
func (p *Prayer) adjustHighLatTimes(times []float64) []float64 {

	var nightTime float64 = commons.TimeDiff(times[4], times[1]) // sunset to sunrise

	// Adjust Fajr
	var FajrDiff float64 = p.nightPortion(p.MethodParams[p.CalcMethod][0]) * nightTime
	if math.IsNaN(times[0]) || commons.TimeDiff(times[0], times[1]) > FajrDiff {
		times[0] = times[1] - FajrDiff
	}

	// Adjust Isha
	ishaAngle := 0.0
	if p.MethodParams[p.CalcMethod][3] == 0 {
		ishaAngle = p.MethodParams[p.CalcMethod][4]

	} else {
		ishaAngle = 18
	}

	var IshaDiff float64 = p.nightPortion(ishaAngle) * nightTime
	if math.IsNaN(times[6]) || commons.TimeDiff(times[4], times[6]) > IshaDiff {
		times[6] = times[4] + IshaDiff
	}

	// Adjust Maghrib
	var MaghribAngle float64 = 0.0
	if p.MethodParams[p.CalcMethod][1] == 0 {
		MaghribAngle = p.MethodParams[p.CalcMethod][2]
	} else {
		MaghribAngle = 4
	}
	MaghribDiff := p.nightPortion(MaghribAngle) * nightTime
	if math.IsNaN(times[5]) || commons.TimeDiff(times[4], times[5]) > MaghribDiff {
		times[5] = times[4] + MaghribDiff
	}
	return times
}

// the night portion used for adjusting times in higher latitudes
func (p *Prayer) nightPortion(angle float64) float64 {
	val := float64(0.0)
	if p.AdjustHighLats == constants.ANGLE_BASED {
		val = 1.0 / 60.0 * angle
	}
	if p.AdjustHighLats == constants.MID_NIGHT {
		val = 1.0 / 2.0
	}
	if p.AdjustHighLats == constants.ONE_SEVENTH {
		val = 1.0 / 7.0
	}
	return val
}

//var methodParams2 = map[int][]float64{JAFARI: []float64{16, 0, 4, 0, 14}}
func (p *Prayer) dayPortion(times []float64) []float64 {
	for i := 0; i < len(times); i++ {
		times[i] /= 24
	}
	return times
}
