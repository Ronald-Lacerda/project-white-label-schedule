package auth

import "testing"

func TestDefaultBusinessHourSeeds_StartClosedForAllDays(t *testing.T) {
	hours := defaultBusinessHourSeeds()

	if len(hours) != 7 {
		t.Fatalf("expected 7 business hour seeds, got %d", len(hours))
	}

	for day, hour := range hours {
		if hour.DayOfWeek != day {
			t.Fatalf("expected day_of_week %d, got %d", day, hour.DayOfWeek)
		}
		if !hour.IsClosed {
			t.Fatalf("expected day %d to start closed", day)
		}
		if hour.OpenTime != "08:00:00" {
			t.Fatalf("expected open_time 08:00:00 for day %d, got %s", day, hour.OpenTime)
		}
		if hour.CloseTime != "18:00:00" {
			t.Fatalf("expected close_time 18:00:00 for day %d, got %s", day, hour.CloseTime)
		}
	}
}
