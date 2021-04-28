package main

import (
	"fmt"

	"github.com/CapregSoft/go-prayer-time/constants"
	"github.com/CapregSoft/go-prayer-time/prayers"
)

func main() {
	latitude := 33.57368163412395

	longitude := 73.17308661244054
	timezone := 5
	// Test Prayer times here
	//PrayTime prayers = new PrayTime();
	pray := &prayers.Prayer{}
	pray.Init()
	pray.TimeFormat = constants.TIME_12
	pray.CalcMethod = constants.JAFARI
	pray.AsrJuristic = constants.SHAFII
	pray.AdjustHighLats = constants.ANGLE_BASED

	//offsets = []int{0, 0, 0, 0, 0, 0, 0} // {Fajr,Sunrise,Dhuhr,Asr,Sunset,Maghrib,Isha}
	//prayers.tune(offsets)

	prayerTimes := pray.GetPrayerTimes(2021, 4, 29, latitude, longitude, timezone)
	prayerNames := pray.TimeName

	for i := 0; i < len(prayerTimes); i++ {
		fmt.Println(prayerNames[i], " - ", prayerTimes[i])
	}
}
