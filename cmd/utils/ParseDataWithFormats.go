package utils

import "time"

func ParseDateWithFormats(dateStr string, formats []string) (time.Time, error) {
	var err error
	var date time.Time

	for _, format := range formats {
		date, err = time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}
	return time.Time{}, err
}
