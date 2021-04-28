package constants

// Calculation Methods

const (
	JAFARI  int = 0
	KARACHI int = 1 // University of Islamic Sciences, Karachi
	ISNA    int = 2 // Islamic Society of North America (ISNA)
	MWL     int = 3 // Muslim World League (MWL)
	MAKKAH  int = 4 // Umm al-Qura, Makkah
	EGYPT   int = 5 // Egyptian General Authority of Survey
	CUSTOM  int = 6 // Custom Setting
	TEHRAN  int = 7 // Institute of Geophysics, University of Tehran
)

// Juristic Methods
const (
	SHAFII int = 0 // Shafii (standard)
	HANAFI int = 1 // Hanafi
)

// Adjusting Methods for Higher Latitudes
const (
	NONE        int = 0 // No adjustment
	MID_NIGHT   int = 1 // middle of night
	ONE_SEVENTH int = 2 // 1/7th of night
	ANGLE_BASED int = 3 // angle/60th of night
)

// Time Formats
const (
	TIME_24    int = 0 // 24-hour format
	TIME_12    int = 1 // 12-hour format
	TIME_12_NS int = 2 // 12-hour format with no suffix
	FLOATING   int = 3 // floating point number
)

const (
	INVALID_TIME string = "----"
)
