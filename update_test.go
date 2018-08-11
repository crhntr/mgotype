package mongotype_test

import (
	"encoding/json"
	"testing"

	mongotype "github.com/crhntr/mongotype"
)

func TestUpdate(t *testing.T) {
	t.Run("when CurrentDate is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.CurrentDate("last_updated")
	})
	t.Run("when CurrentDateAsTimestamp is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.CurrentDateAsTimestamp("update_timestamp")
	})
	t.Run("when CurrentDateAsDate is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.CurrentDateAsDate("last_updated")
	})
	t.Run("when IncrementInt is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.IncrementInt("update_count", 1)
	})
	t.Run("when IncrementFloat64 is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.IncrementFloat64("price", 0.1)
	})
	t.Run("when Multiply is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Multiply("discount", 0.999)
	})
	t.Run("when Minimum is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Minimum("on_order", 100)
	})
	t.Run("when Maximum is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Maximum("sellable", 1000)
	})
	t.Run("when Rename is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Rename("sellable", "for_sale")
	})
	t.Run("when Set is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Set("description", "this is a product")
	})
	t.Run("when SetOnInsert is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.SetOnInsert("sold", 0)
	})
	t.Run("when Unset is called", func(t *testing.T) {
		up := mongotype.NewUpdate()
		up.Unset("invisable")
	})
}

func TestUpdate_Matches(t *testing.T) {
	t.Run("when two updates are equal", func(t *testing.T) {
		update1 := mongotype.NewUpdate()
		update1.Set("foo", "bar")
		update1.CurrentDate("today")

		update2 := mongotype.NewUpdate()
		update2.CurrentDate("today")
		update2.Set("foo", "bar")

		if err := update1.Match(update2); err != nil {
			t.Errorf("they should match: %s", err)
		}
	})

	t.Run("when two updated are not equal", func(t *testing.T) {
		update1 := mongotype.NewUpdate()
		update1.Set("foo", "bar")
		update1.CurrentDate("today")

		update2 := mongotype.NewUpdate()
		update2.CurrentDate("today")
		update2.Set("bar", "foo")

		if err := update1.Match(update2); err == nil {
			t.Errorf("they should not match")
		}
	})
}

func TestUpdate_String(t *testing.T) {
	update := mongotype.NewUpdate()
	update.Set("foo", "bar")
	exp := `{"$set":{"foo":"bar"}}`
	if str := update.String(); str != exp {
		t.Errorf("%#v got %q expected %q", update, str, exp)
	}
}

func TestUpdate_UnmarshalJSON(t *testing.T) {
	t.Run("when it is passed valid json", func(t *testing.T) {
		var update mongotype.Update

		if err := json.Unmarshal([]byte(`{"$set":{"foo":"bar"}}`), &update); err != nil {
			t.Errorf("expected no err: %s", err)
		}
		existing := mongotype.NewUpdate()
		existing.Set("foo", "bar")

		if err := update.Match(existing); err != nil {
			t.Errorf("expected %s to match %s", update, existing)
		}
	})
	t.Run("when it is passed invalid json", func(t *testing.T) {
		var data struct {
			Update mongotype.Update `json:"update"`
		}
		if err := json.Unmarshal([]byte(`{"update": {"$set":{"foo": ["something", wrong]}}, "bar": 5}`), &data); err == nil {
			t.Errorf("it should return an error")
		}
	})
	t.Run("when it is passed an unknown update key", func(t *testing.T) {
		var data struct {
			Update mongotype.Update `json:"update"`
		}
		if err := json.Unmarshal([]byte(`{"update": {"$do-a-thing":{"foo": ["something"]}}, "bar": 5}`), &data); err == nil {
			t.Errorf("it should return an error")
		}
	})
}
