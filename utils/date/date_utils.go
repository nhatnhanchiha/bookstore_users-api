package date

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDbLayOut   = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDbFormat() string {
	return GetNow().Format(apiDbLayOut)
}
