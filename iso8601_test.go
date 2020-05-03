package iso8601

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"
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
	for _, s := range []string{"2020-02-17T11:39:27.658731-02:30", "2020-02-17T11:39:27.658731-0230"} {
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

func TestParseJSON(t *testing.T) {
	var jsonStruct struct {
		Time Time `json:"time"`
	}

	for _, timeStamp := range []string{
		"1900-12-31T00:10:20Z", "1900-12-31T001020Z", "19001231T00:10:20Z", "19001231T001020Z",
		"1900-12-31t00:10:20Z", "1900-12-31t001020Z", "19001231t00:10:20Z", "19001231t001020Z",
		"1900-12-31 00:10:20Z", "1900-12-31 001020Z", "19001231 00:10:20Z", "19001231 001020Z",
		"2020-02-17T11:39:27.658731+00:00", "2020-02-17T11:39:27.658731Z", "2020-02-17T11:39:27.658731-02:30",
	} {
		jsonText := `{"time": "` + timeStamp + `"}`
		if err := json.Unmarshal([]byte(jsonText), &jsonStruct); err != nil {
			t.Errorf("Failed to unmarshall JSON: %s: %#v\n", jsonText, err)
		}
	}

	if err := json.Unmarshal([]byte(`{"time": null}`), &jsonStruct); err != nil {
		t.Errorf("Failed to unmarhsall null time")
	}

	err := json.Unmarshal([]byte(`{"time": 1}`), &jsonStruct)
	if err == nil {
		t.Errorf("Expected a time.ParseError")
	} else if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("Expected a time.ParseError")
	}
}

func TestParseText(t *testing.T) {
	var xmlAttrStruct struct {
		Time Time `xml:"time,attr"`
	}

	var xmlEmbedStruct struct {
		Time Time `xml:"time"`
	}

	for _, timeStamp := range []string{
		"1900-12-31T00:10:20Z", "1900-12-31T001020Z", "19001231T00:10:20Z", "19001231T001020Z",
		"1900-12-31t00:10:20Z", "1900-12-31t001020Z", "19001231t00:10:20Z", "19001231t001020Z",
		"1900-12-31 00:10:20Z", "1900-12-31 001020Z", "19001231 00:10:20Z", "19001231 001020Z",
		"2020-02-17T11:39:27.658731+00:00", "2020-02-17T11:39:27.658731Z", "2020-02-17T11:39:27.658731-02:30",
	} {
		xmlAttrText := `<value time="` + timeStamp + `" />`
		xmlEmbedText := `<value><time>` + timeStamp + `</time></value>`

		if err := xml.Unmarshal([]byte(xmlAttrText), &xmlAttrStruct); err != nil {
			t.Errorf("Failed to unmarshall XML: %s: %#v\n", xmlAttrText, err)
		}

		if err := xml.Unmarshal([]byte(xmlEmbedText), &xmlEmbedStruct); err != nil {
			t.Errorf("Failed to unmarshall XML: %s: %#v\n", xmlEmbedText, err)
		}
	}

}

func TestTimeCompilation(t *testing.T) {
	// This just makes sure the iso8601.Time type is returned from various
	// methods properly.
	oneDay, _ := time.ParseDuration("24h")
	oneMinute, _ := time.ParseDuration("1m")
	start := Date(2001, 1, 1, 3, 15, 45, 0, time.FixedZone("UTC", 0))
	tomorrow := start.Add(oneDay)
	zone333 := time.FixedZone("333", -(3*3600 + 33*60))

	if !tomorrow.After(start.Time) {
		t.Errorf("Expected tomorrow to be after today: tomorrow=%v start=%v oneDay=%v\n", tomorrow, start, oneDay)
	}

	dayAfterTomorrow := start.AddDate(0, 0, 2)
	if !dayAfterTomorrow.After(tomorrow.Time) {
		t.Errorf("Expected dayAfterTomorrow to be after tomorrow")
	}

	rounded := start.Round(oneMinute)
	truncated := start.Truncate(oneMinute)

	roundedDiff := rounded.Sub(start.Time)
	if roundedDiff != 15000000000 {
		t.Errorf("Expected rounded to be 15 seconds after start")
	}

	truncatedDiff := start.Sub(truncated.Time)
	if truncatedDiff != 45000000000 {
		t.Errorf("Expected truncated to be 45 seconds before start")
	}

	localStart := start.Local()
	start.Sub(localStart.Time)

	startUTC := localStart.UTC()
	startUTCDiff := startUTC.Sub(start.Time)
	if startUTCDiff != 0 {
		t.Errorf("Expected startUTC to equal start: startUTC=%v start=%v", startUTC, start)
	}

	start333 := start.In(zone333)
	diff333 := start.Sub(start333.Time)
	if diff333 != 0 {
		t.Errorf("Expected start and start333 to be the same: start=%v start333=%v", start, start333)
	}

	Now().Sub(time.Now())
	Unix(0, 0).Sub(time.Unix(0, 0))
}
