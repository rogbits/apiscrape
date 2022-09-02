package tools

import (
	"testing"
	"time"
)

func TestGenerateTimeSlots(t *testing.T) {
	now := GetNowToTheMinute()
	then := now.Add(-1 * time.Hour)
	slots := GenerateTimeSlots(then, now, time.Minute)
	if len(slots) != 61 {
		t.Fatal()
	}

	then = now.Add(-30 * 24 * time.Hour)
	slots = GenerateTimeSlots(then, now, time.Minute)
	if len(slots) != 43201 {
		t.Fatal()
	}
}
