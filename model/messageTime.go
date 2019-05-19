package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/bsontype"
	"github.com/mongodb/mongo-go-driver/x/bsonx/bsoncore"
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

// MarshalJSON - marshals to JSON as time.Time
func (mt MessageTime) MarshalJSON() ([]byte, error) {
	tmpTime := time.Time(mt)
	return tmpTime.MarshalJSON()
}

// UnmarshalBSONValue - unmarshal BSON to time.Time
func (mt *MessageTime) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	if t == bsontype.DateTime {
		if tmpTime, _, ok := bsoncore.ReadTime(raw); ok {
			*mt = MessageTime(tmpTime)
			return nil
		}
	}

	return errors.New("unable to unmarshal bson MessageTime")
}

//MarshalBSONValue - marshals to BSON as time.Time
func (mt *MessageTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	tmpTime := time.Time{}
	if mt != nil {
		tmpTime = time.Time(*mt)
	}

	return bsontype.DateTime, bsoncore.AppendTime(nil, tmpTime), nil
}
