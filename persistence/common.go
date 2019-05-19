package persistence

import (
	"time"
)

//updateString - Use in an update opertion to figure out if a string value should be removed
func updateString(oldValue, newValue *string) *string {
	if newValue != nil {
		if *newValue == "" { // update is explicitly removing the field
			return nil
		}
		return newValue // update is explicitly setting the field to a new value
	}
	return oldValue // update did not explicitly set this field
}

//updateTime - Use in an update opertion to figure out if a time value should be removed
func updateTime(oldValue, newValue *time.Time) *time.Time {
	if newValue != nil {
		if time.Time.IsZero(*newValue) { // update is explicitly removing the field
			return nil
		}
		return newValue // update is explicitly setting the field to a new value
	}
	return oldValue // update did not explicitly set this field
}
