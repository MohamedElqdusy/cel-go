// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package overloads defines the internal overload identifiers for function and
// operator overloads.
package overloads

const (
	// Boolean logic overloads
	Conditional            = "conditional"
	LogicalAnd             = "logical_and"
	LogicalOr              = "logical_or"
	LogicalNot             = "logical_not"
	Equals                 = "equals"
	NotEquals              = "not_equals"
	LessBool               = "less_bool"
	LessInt64              = "less_int64"
	LessUint64             = "less_uint64"
	LessDouble             = "less_double"
	LessString             = "less_string"
	LessBytes              = "less_bytes"
	LessTimestamp          = "less_timestamp"
	LessDuration           = "less_duration"
	LessEqualsBool         = "less_equals_bool"
	LessEqualsInt64        = "less_equals_int64"
	LessEqualsUint64       = "less_equals_uint64"
	LessEqualsDouble       = "less_equals_double"
	LessEqualsString       = "less_equals_string"
	LessEqualsBytes        = "less_equals_bytes"
	LessEqualsTimestamp    = "less_equals_timestamp"
	LessEqualsDuration     = "less_equals_duration"
	GreaterBool            = "greater_bool"
	GreaterInt64           = "greater_int64"
	GreaterUint64          = "greater_uint64"
	GreaterDouble          = "greater_double"
	GreaterString          = "greater_string"
	GreaterBytes           = "greater_bytes"
	GreaterTimestamp       = "greater_timestamp"
	GreaterDuration        = "greater_duration"
	GreaterEqualsBool      = "greater_equals_bool"
	GreaterEqualsInt64     = "greater_equals_int64"
	GreaterEqualsUint64    = "greater_equals_uint64"
	GreaterEqualsDouble    = "greater_equals_double"
	GreaterEqualsString    = "greater_equals_string"
	GreaterEqualsBytes     = "greater_equals_bytes"
	GreaterEqualsTimestamp = "greater_equals_timestamp"
	GreaterEqualsDuration  = "greater_equals_duration"

	// Math overloads
	AddInt64                   = "add_int64"
	AddUint64                  = "add_uint64"
	AddDouble                  = "add_double"
	AddString                  = "add_string"
	AddBytes                   = "add_bytes"
	AddList                    = "add_list"
	AddTimestampDuration       = "add_timestamp_duration"
	AddDurationTimestamp       = "add_duration_timestamp"
	AddDurationDuration        = "add_duration_duration"
	SubtractInt64              = "subtract_int64"
	SubtractUint64             = "subtract_uint64"
	SubtractDouble             = "subtract_double"
	SubtractTimestampTimestamp = "subtract_timestamp_timestamp"
	SubtractTimestampDuration  = "subtract_timestamp_duration"
	SubtractDurationDuration   = "subtract_duration_duration"
	MultiplyInt64              = "multiply_int64"
	MultiplyUint64             = "multiply_uint64"
	MultiplyDouble             = "multiply_double"
	DivideInt64                = "divide_int64"
	DivideUint64               = "divide_uint64"
	DivideDouble               = "divide_double"
	ModuloInt64                = "modulo_int64"
	ModuloUint64               = "modulo_uint64"
	NegateInt64                = "negate_int64"
	NegateDouble               = "negate_double"

	// Index overloads
	IndexList    = "index_list"
	IndexMap     = "index_map"
	IndexMessage = "index_message" // TODO: introduce concept of types.Message

	// In operators
	DeprecatedIn = "in"
	InList       = "in_list"
	InMap        = "in_map"
	InMessage    = "in_message" // TODO: introduce concept of types.Message

	// Size overloads
	Size           = "size"
	SizeString     = "size_string"
	SizeBytes      = "size_bytes"
	SizeList       = "size_list"
	SizeMap        = "size_map"
	SizeStringInst = "string_size"
	SizeBytesInst  = "bytes_size"
	SizeListInst   = "list_size"
	SizeMapInst    = "map_size"

	// Matches function
	Matches     = "matches"
	MatchString = "matches_string"

	// Time-based functions
	TimeGetFullYear     = "getFullYear"
	TimeGetMonth        = "getMonth"
	TimeGetDayOfYear    = "getDayOfYear"
	TimeGetDate         = "getDate"
	TimeGetDayOfMonth   = "getDayOfMonth"
	TimeGetDayOfWeek    = "getDayOfWeek"
	TimeGetHours        = "getHours"
	TimeGetMinutes      = "getMinutes"
	TimeGetSeconds      = "getSeconds"
	TimeGetMilliseconds = "getMilliseconds"

	// Timestamp overloads for time functions without timezones.
	TimestampToYear                = "timestamp_to_year"
	TimestampToMonth               = "timestamp_to_month"
	TimestampToDayOfYear           = "timestamp_to_day_of_year"
	TimestampToDayOfMonthZeroBased = "timestamp_to_day_of_month"
	TimestampToDayOfMonthOneBased  = "timestamp_to_day_of_month_1_based"
	TimestampToDayOfWeek           = "timestamp_to_day_of_week"
	TimestampToHours               = "timestamp_to_hours"
	TimestampToMinutes             = "timestamp_to_minutes"
	TimestampToSeconds             = "timestamp_to_seconds"
	TimestampToMilliseconds        = "timestamp_to_milliseconds"

	// Timestamp overloads for time functions with timezones.
	TimestampToYearWithTz                = "timestamp_to_year_with_tz"
	TimestampToMonthWithTz               = "timestamp_to_month_with_tz"
	TimestampToDayOfYearWithTz           = "timestamp_to_day_of_year_with_tz"
	TimestampToDayOfMonthZeroBasedWithTz = "timestamp_to_day_of_month_with_tz"
	TimestampToDayOfMonthOneBasedWithTz  = "timestamp_to_day_of_month_1_based_with_tz"
	TimestampToDayOfWeekWithTz           = "timestamp_to_day_of_week_with_tz"
	TimestampToHoursWithTz               = "timestamp_to_hours_with_tz"
	TimestampToMinutesWithTz             = "timestamp_to_minutes_with_tz"
	TimestampToSecondsWithTz             = "timestamp_to_seconds_tz"
	TimestampToMillisecondsWithTz        = "timestamp_to_milliseconds_with_tz"

	// Duration overloads for time functions.
	DurationToHours        = "duration_to_hours"
	DurationToMinutes      = "duration_to_minutes"
	DurationToSeconds      = "duration_to_seconds"
	DurationToMilliseconds = "duration_to_milliseconds"

	// Type conversion methods and overloads
	TypeConvertInt       = "int"
	TypeConvertUint      = "uint"
	TypeConvertDouble    = "double"
	TypeConvertBool      = "bool"
	TypeConvertString    = "string"
	TypeConvertBytes     = "bytes"
	TypeConvertTimestamp = "timestamp"
	TypeConvertDuration  = "duration"
	TypeConvertType      = "type"
	TypeConvertDyn       = "dyn"

	// Int conversion functions.
	IntToInt       = "int64_to_int64"
	UintToInt      = "uint64_to_int64"
	DoubleToInt    = "double_to_int64"
	StringToInt    = "string_to_int64"
	TimestampToInt = "timestamp_to_int64"
	DurationToInt  = "duration_to_int64"

	// Uint conversion functions.
	UintToUint   = "uint64_to_uint64"
	IntToUint    = "int64_to_uint64"
	DoubleToUint = "double_to_uint64"
	StringToUint = "string_to_uint64"

	// Double conversion functions.
	DoubleToDouble = "double_to_double"
	IntToDouble    = "int64_to_double"
	UintToDouble   = "uint64_to_double"
	StringToDouble = "string_to_double"

	// Bool conversion functions.
	BoolToBool   = "bool_to_bool"
	StringToBool = "string_to_bool"

	// Bytes conversion functions.
	BytesToBytes  = "bytes_to_bytes"
	StringToBytes = "string_to_bytes"

	// String conversion functions.
	StringToString    = "string_to_string"
	BoolToString      = "bool_to_string"
	IntToString       = "int64_to_string"
	UintToString      = "uint64_to_string"
	DoubleToString    = "double_to_string"
	BytesToString     = "bytes_to_string"
	TimestampToString = "timestamp_to_string"
	DurationToString  = "duration_to_string"

	// Timestamp conversion functions
	TimestampToTimestamp = "timestamp_to_timestamp"
	StringToTimestamp    = "string_to_timestamp"
	IntToTimestamp       = "int64_to_timestamp"

	// Convert duration from string
	DurationToDuration = "duration_to_duration"
	StringToDuration   = "string_to_duration"
	IntToDuration      = "int64_to_duration"

	// Convert to dyn
	ToDyn = "to_dyn"

	// Comprehensions helper methods, not directly accessible via a developer.
	Iterator = "@iterator"
	HasNext  = "@hasNext"
	Next     = "@next"
)
