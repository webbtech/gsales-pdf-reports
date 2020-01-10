package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DBHandler interface
type DBHandler interface {
	Close()
	GetDay(time.Time, primitive.ObjectID) (bson.M, error)
	GetEmployee(primitive.ObjectID) (*Employee, error)
	GetJournals(string, primitive.ObjectID) ([]*Journal, error)
	GetShift(string, primitive.ObjectID) (*Sales, error)
	GetStation(primitive.ObjectID) (*Station, error)
}

// Record interface
type Record interface {
	GetRecord()
}

// ===================== Helper Functions ====================================================== //

// SetFloat function
func SetFloat(num interface{}) float64 {

	var ret float64
	switch v := num.(type) {
	case *float64:
		// need to check for nil here to deal with null db values
		if v == nil {
			ret = 0.00
		} else {
			ret = *v
		}
	case float64:
		ret = v
	default:
		ret = 0.00
	}

	return ret
}

// SetString function
func SetString(s interface{}) string {
	var ret string
	switch v := s.(type) {
	case *string:
		// need to check for nil here to deal with null db values
		if v == nil {
			ret = ""
		} else {
			ret = *v
		}
	case string:
		ret = v
	default:
		ret = ""
	}

	return ret
}
