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

package common

// Location interface to represent a location within Source.
type Location interface {
	Line() int   // 1-based line number within source.
	Column() int // 0-based column number within source.
}

// SourceLocation helper type to manually construct a location.
type SourceLocation struct {
	line   int
	column int
}

var (
	// Ensure the SourceLocation implements the Location interface.
	_          Location = &SourceLocation{}
	NoLocation          = &SourceLocation{-1, -1}
)

// Create a new location.
func NewLocation(line, column int) Location {
	return &SourceLocation{
		line:   line,
		column: column}
}

func (l *SourceLocation) Line() int {
	return l.line
}

func (l *SourceLocation) Column() int {
	return l.column
}
