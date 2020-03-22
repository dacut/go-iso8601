# Package iso8601
    `import "github.com/dacut/go-iso8601"`

Proper ISO 8601 (**not RFC 3339**) time parsing and formatting for Golang.

Using Go's `time.RFC3339` format string to parse timestamps will oftain fail
on valid ISO 8601 (and certain non-Go-produced RFC 3339) timestamps, e.g.:

* `20171031T235959Z` -- compact ISO 8601 representation
* `2017-10-31 23:59:59Z` -- RFC 3339 alternate separator
* `2017-10-31t23:59:59Z` -- ISO 8601 with lowercase 't' separator
* `2017-10-31T23:59:59.123Z` -- RFC 3339 with milliseconds

This package handles these formats. It properly rejects certain invalid formats
that might be accepted by naive regular expressions:

* `2017-1031T235959Z` -- separators must be all or none for a date
* `20171031T23:5959Z` -- same with time

## Index
ParseISO8601Timestamp converts an ISO 8601 timestamp into a time.Time result.
Compared to time.Parse(time.RFC3339) and time.Parse(time.RFC3339Nano), this
accepts the full range of ISO 8601 formats and the RFC 3339 variants.
```go
func ParseISO8601Timestamp(s string) (time.Time, error)
```
