// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import "time"

// TimePeriod is a simple object that allows you to easily use seconds, minutes, hours, and days to define the length
// of a time period. It exists in this package to make it easy to unmarshall JSON configuration that includes such
// a duration.
type TimePeriod struct {
	Seconds int64
	Minutes int64
	Hours   int64
	Days    int64
}

// Duration returns a time.Duration object with a value equal to the sum of the Seconds, Minutes, Hours, and Days values
// of the TimePeriod object.
func (t TimePeriod) Duration() time.Duration {
	return time.Duration(t.Seconds)*time.Second +
		time.Duration(t.Minutes)*time.Minute +
		time.Duration(t.Hours)*time.Hour +
		time.Duration(t.Days*24)*time.Hour
}

// InSeconds returns the number of seconds that the TimePeriod represents. Many use cases in Splunk need a time period
// to be represented in seconds, so this is a convenience method to obtain that.
func (t TimePeriod) InSeconds() int64 {
	return int64(t.Duration().Seconds())
}

// InMinutes returns the number of minutes that the TimePeriod represents, truncated to an integer.
func (t TimePeriod) InMinutes() int64 {
	return int64(t.Duration().Minutes())
}

// InHours returns the number of hours that the TimePeriod represents, truncated to an integer.
func (t TimePeriod) InHours() int64 {
	return int64(t.Duration().Hours())
}

// InDays returns the number of days that the TimePeriod represents, truncated to an integer.
func (t TimePeriod) InDays() int64 {
	return int64(t.Duration().Hours() / 24)
}
