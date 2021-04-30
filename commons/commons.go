package commons

import (
	"fmt"
	"math"
	"time"

	"github.com/CapregSoft/go-prayer-time/constants"
)

func FixAngle(angle float64) float64 {

	angle = angle - 360.0*(math.Floor(angle/360.0))

	if angle < 0 {
		angle = angle + 360.0
	}
	return angle
}

// range reduce hours to 0..23
func FixHour(hour float64) float64 {

	hour = hour - 24.0*(math.Floor(hour/24.0))
	if hour < 0 {
		hour += 24.0
	}
	return hour
}

// radian to degree
func RadianToDegree(radian float64) float64 {

	return (radian * 180.0) / math.Pi
}

// deree to radian
func DegreeToRadian(degree float64) float64 {
	return (degree * math.Pi) / 180.0
}

// degree sin
func Dsin(d float64) float64 {
	return math.Sin(DegreeToRadian(d))
}

// degree cos
func Dcos(d float64) float64 {
	return math.Cos(DegreeToRadian(d))
}

// degree tan
func Dtan(d float64) float64 {
	return math.Tan(DegreeToRadian(d))
}

// degree arcsin
func Darcsin(x float64) float64 {
	return RadianToDegree(math.Asin(x))
}

// degree arccos
func Darccos(x float64) float64 {
	return RadianToDegree(math.Acos(x))
}

// degree arctan
func Darctan(x float64) float64 {
	return RadianToDegree(math.Atan(x))
}

// degree arctan2
func Darctan2(y float64, x float64) float64 {
	return RadianToDegree(math.Atan2(y, x))
}

// degree arccot
func Darccot(x float64) float64 {

	return RadianToDegree(math.Atan2(1.0, x))
}

// ---------------------- Time-Zone Functions -----------------------
// compute local time-zone for a specific date

func GetTimeZone1() float64 {
	t := time.Now()
	_, offset := t.Zone()
	//fmt.Println(zone, offset)
	return float64(offset) / 3600.0

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

// ---------------------- Calculation Functions -----------------------
// References:
// http://www.ummah.net/astronomy/saltime
// http://aa.usno.navy.mil/faq/docs/SunApprox.html
// compute declination angle of sun and equation of time
func SunPosition(jd float64) []float64 {
	D := jd - 2451545.0
	g := FixAngle(357.529 + 0.98560028*D)
	q := FixAngle(280.459 + 0.98564736*D)
	L := FixAngle(q + (1.915 * Dsin(g)) + (0.020 * Dsin(2*g)))

	e := 23.439 - (0.00000036 * D)

	d := Darcsin(Dsin(e) * Dsin(L))

	RA := Darctan2((Dcos(e)*Dsin(L)), (Dcos(L))) / 15.0
	RA = FixHour(RA)
	EqT := q/15.0 - RA
	return []float64{d, EqT}
}

// compute equation of time
func EquationOfTime(jd float64) float64 {

	return SunPosition(jd)[1]
}

// compute declination angle of sun
func SunDeclination(jd float64) float64 {
	return SunPosition(jd)[0]
}

// ---------------------- Misc Functions -----------------------
// compute the difference between two times
func TimeDiff(time1 float64, time2 float64) float64 {
	return FixHour(time2 - time1)
}
func FloatToTime24(time float64) string {
	if time < 0 {
		return constants.INVALID_TIME
	}
	time = FixHour(time + 0.5/60) // add 0.5 minutes to round
	hours := math.Floor(time)
	minutes := math.Floor((time - hours) * 60)
	return fmt.Sprintf("%s%s%s", TwoDigitsFormat(int(hours)), ":", TwoDigitsFormat(int(minutes)))
}

// convert float hours to 12h format
func FloatToTime12(time float64, noSuffix bool) string {
	if time < 0 {
		return constants.INVALID_TIME
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
	temp := fmt.Sprintf("%s%s%s", TwoDigitsFormat(int(hours)), ":", TwoDigitsFormat(int(minutes)))
	if !noSuffix {
		temp = fmt.Sprintf("%s%s", temp, suffix)
	}
	return temp
}

// convert float hours to 12h format with no suffix
func FloatToTime12NS(time float64) string {
	return FloatToTime12(time, true)
}
func TwoDigitsFormat(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	} else {
		return fmt.Sprintf("%d", num)
	}
}
