package main

import (
	"testing"
	"time"

	"github.com/CapregSoft/go-prayer-time/constants"
	"github.com/CapregSoft/go-prayer-time/prayers"
)

//Hanfi juristic
//Based on: University of Islamic Sciences, Karachi · Change GMT+05:00 · Times may vary
//TEST DATA from google search
var data = [][]string{
	//{"03:50", "05:21", "12:05", "16:51", "18:49", "20:20"},
	{"03:49", "05:20", "12:05", "16:52", "18:50", "20:21"},
	{"03:47", "05:19", "12:04", "16:52", "18:50", "20:22"},
	{"03:46", "05:18", "12:04", "16:53", "18:51", "20:23"},
	{"03:45", "05:17", "12:04", "16:53", "18:52", "20:24"},
	{"03:44", "05:16", "12:04", "16:53", "18:53", "20:25"},
	{"03:42", "05:15", "12:04", "16:54", "18:53", "20:26"},
	{"03:41", "05:14", "12:04", "16:54", "18:54", "20:28"},
	{"03:40", "05:13", "12:04", "16:54", "18:55", "20:29"},
	{"03:39", "05:12", "12:04", "16:55", "18:56", "20:30"},
	{"03:37", "05:12", "12:04", "16:55", "18:56", "20:31"},
	{"03:36", "05:11", "12:04", "16:56", "18:57", "20:32"},
	{"03:35", "05:10", "12:04", "16:56", "18:58", "20:33"},
	{"03:34", "05:09", "12:04", "16:56", "18:59", "20:34"},
	{"03:33", "05:08", "12:04", "16:57", "18:59", "20:35"},
	{"03:32", "05:08", "12:04", "16:57", "19:00", "20:36"},
	{"03:31", "05:07", "12:04", "16:57", "19:01", "20:37"},
	{"03:30", "05:06", "12:04", "16:58", "19:02", "20:38"},
	{"03:29", "05:05", "12:04", "16:58", "19:02", "20:39"},
	{"03:28", "05:05", "12:04", "16:58", "19:03", "20:40"},
	{"03:27", "05:04", "12:04", "16:59", "19:04", "20:41"},
	{"03:26", "05:04", "12:04", "16:59", "19:04", "20:42"},
	{"03:25", "05:03", "12:04", "17:00", "19:05", "20:43"},
	{"03:24", "05:02", "12:04", "17:00", "19:06", "20:44"},
	{"03:23", "05:02", "12:04", "17:00", "19:07", "20:45"},
	{"03:22", "05:01", "12:04", "17:01", "19:07", "20:46"},
	{"03:22", "05:01", "12:04", "17:01", "19:08", "20:47"},
	{"03:21", "05:00", "12:04", "17:01", "19:09", "20:48"},
	{"03:20", "05:00", "12:04", "17:02", "19:09", "20:49"},
	{"03:20", "05:00", "12:05", "17:02", "19:10", "20:50"},
}

func TestPrac(t *testing.T) {

	myDate := time.Now()

	latitude := 33.57368163412395

	longitude := 73.17308661244054
	timezone := 5

	pray := &prayers.Prayer{}
	pray.Init()
	pray.TimeFormat = constants.TIME_24
	pray.CalcMethod = constants.KARACHI
	pray.AsrJuristic = constants.HANAFI
	pray.AdjustHighLats = constants.ANGLE_BASED

	//offsets = []int{0, 0, 0, 0, 0, 0, 0} // {Fajr,Sunrise,Dhuhr,Asr,Sunset,Maghrib,Isha}
	//prayers.tune(offsets)

	for i := 0; i < 29; i++ {
		year, month, day := myDate.Date()
		prayerTimes := pray.GetPrayerTimes(year, int(month), (day), latitude, longitude, timezone)

		//newprayerTimes remove sun set from list
		newprayerTimes := func(prayerTimes []string) []string {
			temp := make([]string, 6)
			idx := 0
			for j := 0; j < len(prayerTimes); j++ {
				if j == 4 {

				} else {
					temp[idx] = prayerTimes[j]
					idx++
				}
			}
			return temp
		}(prayerTimes)

		for j := 0; j < len(newprayerTimes); j++ {
			if newprayerTimes[j] != data[i][j] {
				t.Errorf("got %s, wanted %s -- %d%d --- %d-%d-%d", newprayerTimes[j], data[i][j], i, j, year, month, day)
			}
		}
		myDate = myDate.Add(time.Hour * 24)
	}
}
