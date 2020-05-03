package iso8601

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ISO8601CompactFormat is a layout for time.Format that prints a time value
// in the most compact ISO8601 format available. This assumes the time value
// is in UTC, and returns the time zone as 'Z'.
const ISO8601CompactFormat = "20060102T150405Z"

var iso8601Variants [6]*regexp.Regexp

func init() {
	iso8601Variants[0] = regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])[Tt ]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|6[01])(?:[.,]([0-9]{1,9}))?(Z|[-+][01][0-9]:?(?:[0-5][0-9])?)$`)
	iso8601Variants[1] = regexp.MustCompile(`^([0-9]{4})(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])[Tt ]([01][0-9]|2[0-3])([0-5][0-9])([0-5][0-9]|6[01])(?:[.,]([0-9]{1,9}))?(Z|[-+][01][0-9]:?(?:[0-5][0-9])?)$`)
	iso8601Variants[2] = regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])[Tt ]([01][0-9]|2[0-3])([0-5][0-9])([0-5][0-9]|6[01])(?:[.,]([0-9]{1,9}))?(Z|[-+][01][0-9]:?(?:[0-5][0-9])?)$`)
	iso8601Variants[3] = regexp.MustCompile(`^([0-9]{4})(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])[Tt ]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|6[01])(?:[.,]([0-9]{1,9}))?(Z|[-+][01][0-9]:?(?:[0-5][0-9])?)$`)
	iso8601Variants[4] = regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
	iso8601Variants[5] = regexp.MustCompile(`^([0-9]{4})(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`)
}

func atoip(s string) int {
	if value, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprintf("Failed to convert %#v to int", s))
	} else {
		return int(value)
	}
}

// ParseISO8601Timestamp converts an ISO 8601 timestamp into a time.Time
// result. Compared to time.Parse(time.RFC3339) and
// time.Parse(time.RFC3339Nano), this accepts the full range of ISO 8601
// formats and the RFC 3339 variants.
func ParseISO8601Timestamp(s string) (time.Time, error) {
	for _, re := range iso8601Variants {
		if re.MatchString(s) {
			var year, day, hour, minute, second, nanosecs int
			var month time.Month
			var fracSecStr, tzStr string

			parts := re.FindStringSubmatch(s)
			year = atoip(parts[1])
			month = time.Month(atoip(parts[2]))
			day = atoip(parts[3])
			nanosecs = 0

			if len(parts) > 4 {
				hour = atoip(parts[4])
				minute = atoip(parts[5])
				second = atoip(parts[6])
				fracSecStr = parts[7]

				// fractional seconds don't need to be nanosecond resolution.
				// Pad the right with zeros to make it so.
				if fracSecStr != "" {
					for len(fracSecStr) < 9 {
						fracSecStr = fracSecStr + "0"
					}

					nanosecs = atoip(fracSecStr)
				}
			}

			var loc *time.Location

			if len(parts) > 4 {
				tzStr = parts[8]
			} else {
				tzStr = "Z"
			}

			if tzStr == "Z" {
				loc = time.UTC
			} else {
				signStr := tzStr[0]
				sign := 1
				if signStr == '-' {
					sign = -1
				}

				tzHour := atoip(tzStr[1:3])
				var tzMin int

				if tzStr[3] == ':' {
					tzMin = atoip(tzStr[4:6])
				} else {
					tzMin = atoip(tzStr[3:5])
				}

				loc = time.FixedZone(
					fmt.Sprintf("%c%02d:%02d", signStr, tzHour, tzMin),
					sign*(tzHour*3600+tzMin*60))
			}

			return time.Date(year, month, day, hour, minute, second, nanosecs, loc), nil
		}
	}

	return time.Time{}, &time.ParseError{
		Layout:     "ISO 8601",
		Value:      s,
		LayoutElem: "",
		ValueElem:  "",
		Message:    ": timestamp is not in ISO 8601 format",
	}
}

// Time extends time.Time to handle valid ISO 8601 (but invalid RFC 3336)
// formatted timestamps in UnmarshalText.
type Time struct {
	time.Time
}

// Date returns the Time corresponding to
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice. In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// Date returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// Date panics if loc is nil.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// Now returns the current local time.
func Now() Time {
	return Time{time.Now()}
}

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such
// value is 1<<63-1 (the largest int64 value).
func Unix(sec int64, nsec int64) Time {
	return Time{time.Unix(sec, nsec)}
}

// Add returns the time t+d.
func (t Time) Add(d time.Duration) Time {
	return Time{t.Time.Add(d)}
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (t Time) AddDate(years int, months int, days int) Time {
	return Time{t.Time.AddDate(years, months, days)}
}

// In returns a copy of t representing the same time instant, but
// with the copy's location information set to loc for display
// purposes.
//
// In panics if loc is nil.
func (t Time) In(loc *time.Location) Time {
	return Time{t.Time.In(loc)}
}

// Local returns t with the location set to local time.
func (t Time) Local() Time {
	return Time{t.Time.Local()}
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Round(d time.Duration) Time {
	return Time{t.Time.Round(d)}
}

// String returns the time formatted using the RFC3339Nano string:
//	"2006-01-02T15:04:05.999999999Z07:00"
func (t Time) String() string {
	return t.Format(time.RFC3339Nano)
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Truncate(d time.Duration) Time {
	return Time{t.Time.Truncate(d)}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in any ISO 8601 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	// Make sure the string is quoted properly.
	s := string(data)
	if len(s) < 2 || (!(s[0] == '"' && s[len(s)-1] == '"') && !(s[0] == '\'' && s[len(s)-1] == '\'')) {
		return &time.ParseError{
			Layout:     "ISO 8601",
			Value:      s,
			LayoutElem: "",
			ValueElem:  "",
			Message:    ": timestamp must be a JSON string literal",
		}
	}

	// Remove the quotation marks.
	s = s[1 : len(s)-1]

	var err error
	t.Time, err = ParseISO8601Timestamp(s)
	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in ISO 8601 format.
func (t *Time) UnmarshalText(data []byte) error {
	var err error
	t.Time, err = ParseISO8601Timestamp(string(data))
	return err
}

// UTC returns t with the location set to UTC.
func (t Time) UTC() Time {
	return Time{t.Time.UTC()}
}
