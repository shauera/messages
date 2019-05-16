package model

import (
	"encoding/json"
	"time"
)

// MessageTime - support "null" as zero time
// we want this to be an alias to time.Time so swagger will pick this as time when creating examples
type MessageTime time.Time

// UnmarshalJSON - used to accomodate passing null as "zero" time in JSON
func (mt *MessageTime) UnmarshalJSON(b []byte) error {
	// first try unmarshaling as string and see if we got null
	var messageTimeString string
	if err := json.Unmarshal(b, &messageTimeString); err == nil {
		if messageTimeString == "null" {
			*mt = MessageTime(time.Time{}) // zero time
			return nil
		}
	}

	// we did not get null as time so try unmarshaling as time
	var tryTime time.Time
	if err := json.Unmarshal(b, &tryTime); err != nil {
		return err
	}
	*mt = MessageTime(tryTime)
	return nil
}

// MarshalJSON - marshals as time.Time
func (mt MessageTime) MarshalJSON() ([]byte, error) {
	tmpTime := time.Time(mt)
	return tmpTime.MarshalJSON()
}
