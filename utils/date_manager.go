package utils

import "time"

func TimeStringToInt64(t string) (int64, error) {
	timeInt64, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, err
	}
	return timeInt64.Unix(), nil
}
