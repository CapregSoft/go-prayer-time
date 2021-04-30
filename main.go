package main

import (
	"fmt"
	"time"

	"github.com/CapregSoft/go-prayer-time/constants"
	"github.com/CapregSoft/go-prayer-time/prayers"
)

func main() {

	myDate := time.Now()
	//year, month, day := myDate.Date()
	latitude := 33.57368163412395

	longitude := 73.17308661244054
	timezone := 5
	// Test Prayer times here
	//PrayTime prayers = new PrayTime();
	pray := &prayers.Prayer{}
	pray.Init()
	pray.TimeFormat = constants.TIME_12_NS
	pray.CalcMethod = constants.KARACHI
	pray.AsrJuristic = constants.HANAFI
	pray.AdjustHighLats = constants.ANGLE_BASED

	for i := 0; i < 1; i++ {
		myDate = myDate.Add(time.Hour * 24)
		//offsets = []int{0, 0, 0, 0, 0, 0, 0} // {Fajr,Sunrise,Dhuhr,Asr,Sunset,Maghrib,Isha}
		//prayers.tune(offsets)
		//2021-05-24
		prayerTimes := pray.GetPrayerTimes(2021, 5, 24, latitude, longitude, timezone)

		//prayerTimes := pray.GetPrayerTimes(year, int(month), (day), latitude, longitude, timezone)
		prayerNames := pray.TimeName

		for i := 0; i < len(prayerTimes); i++ {
			fmt.Println(prayerNames[i], " - ", prayerTimes[i])
		}
	}

	//prayerTimes := pray.GetPrayerTimes(2021, 4, 29, latitude, longitude, timezone)

}
