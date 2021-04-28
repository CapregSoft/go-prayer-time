package prayers

import (
	"fmt"
	"math"
)

const (
	JAFARI  int = 0
	KARACHI int = 1
	ISNA    int = 2
	MWL     int = 3
	MAKKAH  int = 4
	EGYPT   int = 5
	CUSTOM  int = 6
	TEHRAN  int = 7
)

const (
	SHAFII int = 0
	HANAFI int = 1
)

const (
	NONE        int = 0
	MID_NIGHT   int = 1
	ONE_SEVENTH int = 2
	ANGLE_BASED int = 3
)

const (
	TIME_24    int = 0
	TIME_12    int = 1
	TIME_12_NS int = 2
	FLOATING   int = 3
)

var names = [7]string{
	"Fajr",
	"Sunrise",
	"Dhuhr",
	"Asr",
	"Sunset",
	"Maghrib",
	"Isha",
}

var calcMethod int = 3     // caculation method
var asrJuristic int        // Juristic method for Asr
var dhuhrMinutes int = 0   // minutes after mid-day for Dhuhr
var adjustHighLats int = 1 // adjusting method for higher latitudes

var timeFormat int = 0 // time format

var lat float64   // latitude
var lng float64   // longitude
var timeZone int  // time-zone
var JDate float64 // Julian date

var times []int

//--------------------- Technical Settings --------------------

var numIterations int = 1 // number of iterations needed to compute times

//------------------- Calc Method Parameters --------------------

var methodParams = map[int][]float64{
	JAFARI:  []float64{16, 0, 4, 0, 14},
	KARACHI: []float64{18, 1, 0, 0, 18},
	ISNA:    []float64{15, 1, 0, 0, 15},
	MWL:     []float64{18, 1, 0, 0, 17},
	MAKKAH:  []float64{18.5, 1, 0, 1, 90},
	EGYPT:   []float64{19.5, 1, 0, 0, 17.5},
	TEHRAN:  []float64{17.7, 0, 4.5, 0, 14},
	CUSTOM:  []float64{18, 1, 0, 0, 17},
}

const InvalidTime string = "----"

func FloatToTime24(time float64) string {
	if time < 0 {
		return InvalidTime
	}
	time = FixHour(time + 0.5/60) // add 0.5 minutes to round
	hours := math.Floor(time)
	minutes := math.Floor((time - hours) * 60)
	return fmt.Sprintf("%s%s%s", twoDigitsFormat(int(hours)), ":", twoDigitsFormat(int(minutes)))
}

// convert float hours to 12h format
func FloatToTime12(time float64, noSuffix bool) string {
	if time < 0 {
		return InvalidTime
	}
	time = FixHour(time + 0.5/60) // add 0.5 minutes to round
	hours := math.Floor(time)
	minutes := math.Floor((time - hours) * 60)

	var suffix string = ""
	if hours > 12 {
		suffix = "pm"
	} else {
		suffix = "am"
	}
	hours = float64((int(hours)+12-1)%12 + 1)
	temp := fmt.Sprintf("%d%s%s", int(hours), ":", twoDigitsFormat(int(minutes)))
	if !noSuffix {
		temp = fmt.Sprintf("%s%s", temp, suffix)
	}
	return temp
}

// convert float hours to 12h format with no suffix
func FloatToTime12NS(time float64) string {
	return FloatToTime12(time, true)
}

//---------------------- Compute Prayer Times -----------------------

// return prayer times for a given date
func getDatePrayerTimes(year int, month int, day int, latitude float64, longitude float64, tZone int) []string {
	lat = latitude
	lng = longitude
	timeZone = tZone
	JDate = JulianDate(year, month, day) - longitude/(15*24)

	return computeDayTimes()
}

// compute declination angle of sun and equation of time
func sunPosition(jd float64) []float64 {
	D := jd - 2451545.0
	g := FixAngle(357.529 + 0.98560028*D)
	q := FixAngle(280.459 + 0.98564736*D)
	L := FixAngle(q + 1.915*dsin(g) + 0.020*dsin(2*g))

	//R := 1.00014 - 0.01671*dcos(g) - 0.00014*dcos(2*g)
	e := 23.439 - 0.00000036*D

	d := darcsin(dsin(e) * dsin(L))
	RA := darctan2(dcos(e)*dsin(L), dcos(L)) / 15
	RA = FixHour(RA)
	EqT := q/15.0 - RA

	return []float64{d, EqT}
}

// compute equation of time
func equationOfTime(jd float64) float64 {

	return sunPosition(jd)[1]
}

// compute declination angle of sun
func sunDeclination(jd float64) float64 {
	return sunPosition(jd)[0]
}

// compute mid-day (Dhuhr, Zawal) time
func computeMidDay(t float64) float64 {
	T := equationOfTime(JDate + t)
	Z := FixHour(12 - T)
	return Z
}

// compute time for a given angle G
func computeTime(G float64, t float64) float64 {
	//System.out.println("G: "+G)

	D := sunDeclination(JDate + t)
	Z := computeMidDay(t)
	V := (1.0 / 15.0) * darccos((-dsin(G)-dsin(D)*dsin(lat))/
		(dcos(D)*dcos(lat)))

	if G > 90 {
		Z = Z + (-V)
	} else {
		Z = Z + V
	}
	return Z
}

// compute the time of Asr
func computeAsr(step int, t float64) float64 { // Shafii: step=1, Hanafi: step=2

	D := sunDeclination(JDate + t)
	G := -darccot(float64(step) + dtan(math.Abs(lat-D)))
	return computeTime(G, t)
}

//---------------------- Compute Prayer Times -----------------------

// compute prayer times at given julian date
func computeTimes(times []float64) []float64 {

	t := dayPortion(times)

	Fajr := computeTime(180-methodParams[calcMethod][0], t[0])
	Sunrise := computeTime(180-0.833, t[1])
	Dhuhr := computeMidDay(t[2])
	Asr := computeAsr(1+asrJuristic, t[3])
	Sunset := computeTime(0.833, t[4])
	Maghrib := computeTime(methodParams[calcMethod][2], t[5])
	Isha := computeTime(methodParams[calcMethod][4], t[6])

	return []float64{Fajr, Sunrise, Dhuhr, Asr, Sunset, Maghrib, Isha}
}

func adjustHighLatTimes(times []float64) []float64 {

	var nightTime float64 = GetTimeDifference(times[4], times[1]) // sunset to sunrise

	// Adjust Fajr
	var FajrDiff float64 = nightPortion(methodParams[calcMethod][0]) * nightTime
	if GetTimeDifference(times[0], times[1]) > FajrDiff {
		times[0] = times[1] - FajrDiff
	}

	// Adjust Isha
	ishaAngle := 0.0
	if methodParams[calcMethod][3] == 0 {
		ishaAngle = methodParams[calcMethod][4]

	} else {
		ishaAngle = 18
	}

	var IshaDiff float64 = nightPortion(ishaAngle) * nightTime
	if GetTimeDifference(times[4], times[6]) > IshaDiff {
		times[6] = times[4] + IshaDiff
	}

	// Adjust Maghrib
	var MaghribAngle float64 = 0.0
	if methodParams[calcMethod][1] == 0 {
		MaghribAngle = methodParams[calcMethod][2]
	} else {
		MaghribAngle = 4
	}
	MaghribDiff := nightPortion(MaghribAngle) * nightTime
	if GetTimeDifference(times[4], times[5]) > MaghribDiff {
		times[5] = times[4] + MaghribDiff
	}
	return times
}

// the night portion used for adjusting times in higher latitudes
func nightPortion(angle float64) float64 {
	val := float64(0.0)
	if adjustHighLats == ANGLE_BASED {
		val = 1.0 / 60.0 * angle
	}
	if adjustHighLats == MID_NIGHT {
		val = 1.0 / 2.0
	}
	if adjustHighLats == ONE_SEVENTH {
		val = 1.0 / 7.0
	}
	return val
}

//var methodParams2 = map[int][]float64{JAFARI: []float64{16, 0, 4, 0, 14}}
func dayPortion(times []float64) []float64 {
	for i := 0; i < len(times); i++ {
		times[i] /= 24
	}
	return times
}
func computeDayTimes() []string {
	times := []float64{5, 6, 12, 13, 18, 18, 18} //default times

	for i := 0; i < numIterations; i++ {
		times = computeTimes(times)
	}

	times = adjustTimes(times)
	return adjustTimesFormat(times)
}
func adjustTimes(times []float64) []float64 {
	for i := 0; i < 7; i++ {
		times[i] += float64(timeZone) - lng/15.0
	}
	times[2] += float64(dhuhrMinutes) / 60.0 //Dhuhr
	if methodParams[calcMethod][1] == 1 {    // Maghrib
		times[5] = times[4] + methodParams[calcMethod][2]/60.0
	}
	if methodParams[calcMethod][3] == 1 { // Isha
		times[6] = times[5] + methodParams[calcMethod][4]/60.0
	}
	if adjustHighLats != NONE {
		times = adjustHighLatTimes(times)
	}

	return times
}
func adjustTimesFormat(times []float64) []string {
	formatted := make([]string, len(times))
	if timeFormat == FLOATING {
		for i := 0; i < len(times); i++ {
			formatted[i] = fmt.Sprintf("%f", times[i])
		}
		return formatted
	}
	for i := 0; i < 7; i++ {
		if timeFormat == TIME_12 {
			formatted[i] = FloatToTime12(times[i], true)
		} else if timeFormat == TIME_12_NS {
			formatted[i] = FloatToTime12NS(times[i])
		} else {
			formatted[i] = FloatToTime24(times[i])
		}
	}
	return formatted
}

func GetTimeDifference(c1 float64, c2 float64) float64 {
	diff := FixHour(c2 - c1)
	return diff
}

// add a leading 0 if necessary
func twoDigitsFormat(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	} else {
		return fmt.Sprintf("%d", num)
	}
}

// calculate julian date from a calendar date
func JulianDate(year int, month int, day int) float64 {
	if month <= 2 {
		year -= 1
		month += 12
	}
	yearf, monthf, dayf := float64(year), float64(month), float64(day)

	A := math.Floor(yearf / 100.0)
	B := 2 - A + math.Floor(A/4)

	JD := math.Floor(365.25*(yearf+4716.0)) + math.Floor(30.6001*(monthf+1.0)) + dayf + B - 1524.5
	return JD
}

// degree sin
func dsin(d float64) float64 {
	return math.Sin(DegreeToRadian(d))
}

// degree cos
func dcos(d float64) float64 {
	return math.Cos(DegreeToRadian(d))
}

// degree tan
func dtan(d float64) float64 {
	return math.Tan(DegreeToRadian(d))
}

// degree arcsin
func darcsin(x float64) float64 {
	return RadianToDegree(math.Asin(x))
}

// degree arccos
func darccos(x float64) float64 {
	return RadianToDegree(math.Acos(x))
}

// degree arctan
func darctan(x float64) float64 {
	return RadianToDegree(math.Atan(x))
}

// degree arctan2
func darctan2(y float64, x float64) float64 {
	return RadianToDegree(math.Atan2(y, x))
}

// degree arccot
func darccot(x float64) float64 {
	return RadianToDegree(math.Atan(1 / x))
}

// Radian to Degree
func RadianToDegree(radian float64) float64 {

	return (radian * 180.0) / math.Pi
}

// degree to radian
func DegreeToRadian(degree float64) float64 {
	return (degree * math.Pi) / 180.0
}

func FixAngle(angle float64) float64 {

	angle = angle - 360.0*(math.Floor(angle/360.0))

	if angle < 0 {
		angle = angle + 360.0
	}
	return angle
}
func FixHour(hour float64) float64 {

	hour = hour - 24.0*(math.Floor(hour/24.0))
	if hour < 0 {
		hour += 24.0
	}
	return hour
}
func setTimeFormat(tf int) {
	timeFormat = tf
}
func setCalcMethod(methodId int) {
	calcMethod = methodId
}
func setAsrJuristic(j int) {
	asrJuristic = j
}
func setAdjustHighLats(h int) {
	adjustHighLats = h
}

/*func main() {

	latitude := -37.823689
	longitude := 145.121597
	timezone := 10
	// Test Prayer times here
	//PrayTime prayers = new PrayTime();

	setTimeFormat(TIME_12)
	setCalcMethod(JAFARI)
	setAsrJuristic(SHAFII)
	setAdjustHighLats(ANGLE_BASED)
	offsets := []int{0, 0, 0, 0, 0, 0, 0} // {Fajr,Sunrise,Dhuhr,Asr,Sunset,Maghrib,Isha}
	//prayers.tune(offsets)

	prayerTimes := getPrayerTimes("cal", latitude, longitude, timezone)
	prayerNames := getTimeNames()

	for i := 0; i < prayerTimes.size(); i++ {
		fmt.Println(prayerNames[i], " - ", prayerTimes[i])
	}
	///	fmt.Println(mapOfSlices)
	//fmt.Println(methodParams)
	//fmt.Println(methodParams2)
}
*/
