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

	// Array Update Operators

	UpdateOperatorAddToSet = UpdateOperator("$addToSet")
	UpdateOperatorPop      = UpdateOperator("$pop")
	UpdateOperatorPull     = UpdateOperator("$pull")
	UpdateOperatorPush     = UpdateOperator("$push")
	UpdateOperatorPullAll  = UpdateOperator("$pullAll")
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

func (update Update) CurrentDate(field string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][field] = true
	return nil
}

func (update Update) CurrentDateAsTimestamp(field string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][field] = map[string]string{"$type": "timestamp"}
	return nil
}

func (update Update) CurrentDateAsDate(field string) error {
	ensureUpOp(update.doc, UpdateOperatorCurrentDate)
	update.doc[UpdateOperatorCurrentDate][field] = map[string]string{"$type": "date"}
	return nil
}

func (update Update) IncrementInt(field string, amount int) error {
	ensureUpOp(update.doc, UpdateOperatorInc)
	update.doc[UpdateOperatorInc][field] = amount
	return nil
}

func (update Update) IncrementFloat64(field string, amount float64) error {
	ensureUpOp(update.doc, UpdateOperatorInc)
	update.doc[UpdateOperatorInc][field] = amount
	return nil
}

func (update Update) Multiply(field string, num float64) error {
	ensureUpOp(update.doc, UpdateOperatorMul)
	update.doc[UpdateOperatorMul][field] = num
	return nil
}

func (update Update) Minimum(field string, min float64) error {
	ensureUpOp(update.doc, UpdateOperatorMin)
	update.doc[UpdateOperatorMin][field] = min
	return nil
}

func (update Update) Maximum(field string, max float64) error {
	ensureUpOp(update.doc, UpdateOperatorMax)
	update.doc[UpdateOperatorMax][field] = max
	return nil
}

func (update Update) Rename(oldfield, newfield string) error {
	ensureUpOp(update.doc, UpdateOperatorRename)
	update.doc[UpdateOperatorRename][oldfield] = newfield
	return nil
}

func (update Update) Set(field string, value interface{}) error {
	ensureUpOp(update.doc, UpdateOperatorSet)
	update.doc[UpdateOperatorSet][field] = value
	return nil
}

func (update Update) SetOnInsert(field string, value interface{}) error {
	ensureUpOp(update.doc, UpdateOperatorSetOnInsert)
	update.doc[UpdateOperatorSetOnInsert][field] = value
	return nil
}

func (update Update) Unset(field string) error {
	ensureUpOp(update.doc, UpdateOperatorUnset)
	update.doc[UpdateOperatorUnset][field] = ""
	return nil
}

func (update Update) AddToSet(field string) error { return nil }
func (update Update) Pop(field string) error      { return nil }
func (update Update) Pull(field string) error     { return nil }
func (update Update) Push(field string) error     { return nil }
func (update Update) PullAll(field string) error  { return nil }

func (update *Update) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &update.doc); err != nil {
		return err
	}
	for opfield, _ := range update.doc {
		if !IsUpdateOperator(string(opfield)) {
			return fmt.Errorf("[ %q ] is not a valid update operator", opfield)
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
	for field, val := range doc {
		part, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected json object for field %q", field)
		}
		update.doc[UpdateOperator(field)] = part
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
