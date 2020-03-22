package iso8601

import (
	"testing"
)

func TestValidDates(t *testing.T) {
	// Regular Zulu times
	for _, s := range []string{
		"1900-12-31T00:10:20Z", "1900-12-31T001020Z", "19001231T00:10:20Z", "19001231T001020Z",
		"1900-12-31t00:10:20Z", "1900-12-31t001020Z", "19001231t00:10:20Z", "19001231t001020Z",
		"1900-12-31 00:10:20Z", "1900-12-31 001020Z", "19001231 00:10:20Z", "19001231 001020Z"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 1900 || ts.Month() != 12 || ts.Day() != 31 || ts.Hour() != 0 || ts.Minute() != 10 || ts.Second() != 20 {
			t.Errorf("Incorrect timestamp value for %#v: expected 1900-12-31T00:10:20Z, got %v", s, ts)
		}
	}

	// Date only
	for _, s := range []string{"1900-12-31", "19001231"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 1900 || ts.Month() != 12 || ts.Day() != 31 || ts.Hour() != 0 || ts.Minute() != 0 || ts.Second() != 0 {
			t.Errorf("Incorrect timestamp value for %#v: expected 1900-12-31T00:00:00Z, got %v", s, ts)
		}
	}

	// Fractional second handling
	for _, s := range []string{"2020-02-17T11:39:27.658731+00:00", "2020-02-17T11:39:27.658731Z"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 2020 || ts.Month() != 2 || ts.Day() != 17 || ts.Hour() != 11 || ts.Minute() != 39 || ts.Second() != 27 || ts.Nanosecond() != 658731000 {
			t.Errorf("Incorrect timestamp value for %#v: expected 2020-02-17T11:39:27.658731Z, got %v", s, ts)
		}
	}

	for _, s := range []string{"2020-02-17T11:39:27.6+00:00", "2020-02-17T11:39:27.6Z"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 2020 || ts.Month() != 2 || ts.Day() != 17 || ts.Hour() != 11 || ts.Minute() != 39 || ts.Second() != 27 || ts.Nanosecond() != 600000000 {
			t.Errorf("Incorrect timestamp value for %#v: expected 2020-02-17T11:39:27.6Z, got %v", s, ts)
		}
	}

	// Time zone conversion
	for _, s := range []string{"2020-02-17T11:39:27.658731-02:30"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 2020 || ts.Month() != 2 || ts.Day() != 17 || ts.Hour() != 11 || ts.Minute() != 39 || ts.Second() != 27 || ts.Nanosecond() != 658731000 {
			t.Errorf("Incorrect timestamp value for %#v: expected 2020-02-17T11:39:27.658731Z, got %v", s, ts)
		} else {
			tsUTC := ts.UTC()
			if tsUTC.Hour() != 14 || tsUTC.Minute() != 9 || tsUTC.Second() != 27 || tsUTC.Nanosecond() != 658731000 {
				t.Errorf("Incorrect conversion to UTC for %#v: expected 2020-02-17T14:09:27.658731Z, got %v", s, tsUTC)
			}
		}
	}

	// Leap seconds
	for _, s := range []string{"1900-12-31T00:10:60Z", "1900-12-31T001060Z", "19001231T00:10:60Z", "19001231T001060Z"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 1900 || ts.Month() != 12 || ts.Day() != 31 || ts.Hour() != 0 || !((ts.Minute() == 10 && ts.Second() == 60) || (ts.Minute() == 11 && ts.Second() == 0)) {
			t.Errorf("Incorrect timestamp value for %#v: expected 1900-12-31T00:10:60Z, got %v", s, ts)
		}
	}

	// Double-leap seconds
	for _, s := range []string{"1900-12-31T00:10:61Z", "1900-12-31T001061Z", "19001231T00:10:61Z", "19001231T001061Z"} {
		if ts, err := ParseISO8601Timestamp(s); err != nil {
			t.Errorf("Failed to parse timestamp: %#v %#v\n", s, err)
		} else if ts.Year() != 1900 || ts.Month() != 12 || ts.Day() != 31 || ts.Hour() != 0 || !((ts.Minute() == 10 && ts.Second() == 61) || (ts.Minute() == 11 && ts.Second() == 1)) {
			t.Errorf("Incorrect timestamp value for %#v: expected 1900-12-31T00:10:60Z, got %v", s, ts)
		}
	}

}

func TestInvalidDates(t *testing.T) {
	for _, s := range []string{"1900-1231T00:10:20Z", "1900-12-31T00:1020Z"} {
		if _, err := ParseISO8601Timestamp(s); err == nil {
			t.Errorf("Expected an error on timestamp: %#v\n", s)
		}
	}
}
