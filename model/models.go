package model

type PrayerData struct {
	Fajar   string `json:"fajar"`
	Sunrise string `json:"sunrise"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Sunset  string `json:"sunset"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}
