package tools

import "time"

func GetNowToTheMinute() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), 0, 0, now.Location(),
	)
}

func RoundDownToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GenerateTimeSlots(start time.Time, end time.Time, duration time.Duration) []int64 {
	var slots []int64
	s := start.Unix()
	e := end.Unix()
	for epoch := s; epoch <= e; epoch += int64(duration.Seconds()) {
		slots = append(slots, epoch)
	}
	return slots
}
