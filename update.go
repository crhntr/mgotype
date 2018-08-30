package mgotype

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/globalsign/mgo/bson"
)

type Update struct {
	doc map[UpdateOperator]map[string]interface{}
}

func NewUpdate() Update {
	var update Update
	update.doc = make(map[UpdateOperator]map[string]interface{})
	return update
}

type UpdateOperator string

const (
	// Field Update Operators
	UpdateOperatorCurrentDate = UpdateOperator("$currentDate")
	UpdateOperatorInc         = UpdateOperator("$inc")
	UpdateOperatorMul         = UpdateOperator("$mul")
	UpdateOperatorMin         = UpdateOperator("$min")
	UpdateOperatorMax         = UpdateOperator("$max")
	UpdateOperatorRename      = UpdateOperator("$rename")
	UpdateOperatorSet         = UpdateOperator("$set")
	UpdateOperatorSetOnInsert = UpdateOperator("$setOnInsert")
	UpdateOperatorUnset       = UpdateOperator("$unset")
)

func IsUpdateOperator(str string) bool {
	for _, op := range []UpdateOperator{
		UpdateOperatorCurrentDate,
		UpdateOperatorInc,
		UpdateOperatorMul,
		UpdateOperatorMin,
		UpdateOperatorMax,
		UpdateOperatorRename,
		UpdateOperatorSet,
		UpdateOperatorSetOnInsert,
		UpdateOperatorUnset,
	} {
		if string(op) == str {
			return true
		}
	}
	return false
}

func ensureUpOp(doc map[UpdateOperator]map[string]interface{}, field UpdateOperator) {
	if _, ok := doc[field]; !ok {
		doc[field] = make(map[string]interface{})
	}
}

func (update Update) CurrentDate(key string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][key] = true
	return nil
}

func (update Update) CurrentDateAsTimestamp(key string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][key] = map[string]string{"$type": "timestamp"}
	return nil
}

func (update Update) CurrentDateAsDate(key string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][key] = map[string]string{"$type": "date"}
	return nil
}

func (update Update) IncrementInt(key string, amount int) error {
	ensureUpOp(update.doc, UpdateOperatorInc)
	update.doc[UpdateOperatorInc][key] = amount
	return nil
}

func (update Update) IncrementFloat64(key string, amount float64) error {
	ensureUpOp(update.doc, UpdateOperatorInc)
	update.doc[UpdateOperatorInc][key] = amount
	return nil
}

func (update Update) Multiply(key string, num float64) error {
	ensureUpOp(update.doc, UpdateOperatorMul)
	update.doc[UpdateOperatorMul][key] = num
	return nil
}

func (update Update) Minimum(key string, min float64) error {
	ensureUpOp(update.doc, UpdateOperatorMin)
	update.doc[UpdateOperatorMin][key] = min
	return nil
}

func (update Update) Maximum(key string, max float64) error {
	ensureUpOp(update.doc, UpdateOperatorMax)
	update.doc[UpdateOperatorMax][key] = max
	return nil
}

func (update Update) Rename(oldKey, newKey string) error {
	ensureUpOp(update.doc, UpdateOperatorRename)
	update.doc[UpdateOperatorRename][oldKey] = newKey
	return nil
}

func (update Update) Set(key string, value interface{}) error {
	ensureUpOp(update.doc, UpdateOperatorSet)
	update.doc[UpdateOperatorSet][key] = value
	return nil
}

func (update Update) SetOnInsert(key string, value interface{}) error {
	ensureUpOp(update.doc, UpdateOperatorSetOnInsert)
	update.doc[UpdateOperatorSetOnInsert][key] = value
	return nil
}

func (update Update) Unset(key string) error {
	ensureUpOp(update.doc, UpdateOperatorUnset)
	update.doc[UpdateOperatorUnset][key] = ""
	return nil
}

func (update *Update) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &update.doc); err != nil {
		return err
	}
	for opkey, _ := range update.doc {
		if !IsUpdateOperator(string(opkey)) {
			return fmt.Errorf("[ %q ] is not a valid update operator", opkey)
		}
	}
	return nil
}

func (update Update) MarshalJSON() ([]byte, error) {
	return json.Marshal(update.doc)
}

func (update Update) GetBSON() (interface{}, error) {
	return update.doc, nil
}

func (update *Update) SetBSON(raw bson.Raw) error {
	doc := make(map[string]interface{})
	if err := raw.Unmarshal(&doc); err != nil {
		return err // untested
	}
	if update.doc == nil {
		update.doc = make(map[UpdateOperator]map[string]interface{})
	}
	for key, val := range doc {
		part, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected json object for key %q", key)
		}
		update.doc[UpdateOperator(key)] = part
	}
	return nil
}

type EachModifier struct {
	Items []interface{} `json:"$each" bson:"$each"`
}

func Each(items ...interface{}) EachModifier {
	return EachModifier{items}
}

func (update Update) Match(other Update) error {
	if !reflect.DeepEqual(update.doc, other.doc) {
		return errors.New("failed deep equal")
	}
	return nil
}

func (update Update) String() string {
	buf, _ := json.Marshal(update.doc)
	return string(buf)
}
