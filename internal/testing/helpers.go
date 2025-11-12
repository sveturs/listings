package testing

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// Utility Helpers
// =============================================================================

// MustNewStruct creates a new structpb.Struct or panics.
// Use this in test setup code where you want to fail fast on invalid data.
func MustNewStruct(data map[string]interface{}) *structpb.Struct {
	s, err := structpb.NewStruct(data)
	if err != nil {
		panic(err)
	}
	return s
}

// MustNewValue creates a new structpb.Value or panics.
func MustNewValue(data interface{}) *structpb.Value {
	v, err := structpb.NewValue(data)
	if err != nil {
		panic(err)
	}
	return v
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Int32Ptr returns a pointer to an int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// TimestampNow returns the current timestamp
func TimestampNow() *timestamppb.Timestamp {
	return timestamppb.Now()
}

// TimestampFromTime converts a time.Time to a protobuf Timestamp
func TimestampFromTime(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// TimeNowString returns current time as RFC3339 string (for proto string timestamps)
func TimeNowString() string {
	return time.Now().Format(time.RFC3339)
}

// TimeToString converts time.Time to RFC3339 string
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// TimeYesterday returns yesterday's time
func TimeYesterday() time.Time {
	return time.Now().Add(-24 * time.Hour)
}

// TimeTomorrow returns tomorrow's time
func TimeTomorrow() time.Time {
	return time.Now().Add(24 * time.Hour)
}
