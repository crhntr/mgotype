package mongotype_test

import (
	"encoding/json"
	"testing"

	mongotype "github.com/crhntr/mongotype"
	"github.com/globalsign/mgo/bson"
)

func TestUpdate(t *testing.T) {
	t.Run("when CurrentDate is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.CurrentDate("last_updated"); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when CurrentDateAsTimestamp is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.CurrentDateAsTimestamp("update_timestamp"); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when CurrentDateAsDate is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.CurrentDateAsDate("last_updated"); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when IncrementInt is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.IncrementInt("update_count", 1); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when IncrementFloat64 is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.IncrementFloat64("price", 0.1); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Multiply is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Multiply("discount", 0.999); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Minimum is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Minimum("on_order", 100); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Maximum is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Maximum("sellable", 1000); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Rename is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Rename("sellable", "for_sale"); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Set is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Set("description", "this is a product"); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when SetOnInsert is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.SetOnInsert("sold", 0); err != nil {
			t.Error("it should not return an error")
		}
	})
	t.Run("when Unset is called with valid args", func(t *testing.T) {
		up := mongotype.NewUpdate()
		if err := up.Unset("invisable"); err != nil {
			t.Error("it should not return an error")
		}
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

func TestUpdate_MarshalJSON(t *testing.T) {
	update := mongotype.NewUpdate()
	update.Set("foo", "bar")
	exp := `{"$set":{"foo":"bar"}}`
	buf, err := json.Marshal(update)
	if err != nil {
		t.Fail()
	}
	got := string(buf)
	if got != exp {
		t.Errorf("%#v got %q expected %q", update, got, exp)
	}

	t.Run("when unmarshling bad json", func(t *testing.T) {
		var got mongotype.Update
		if err := got.UnmarshalJSON([]byte(`{`)); err == nil {
			t.Error("it should return an error")
		}
	})
}

func TestUpdate_MarshalingBSON(t *testing.T) {
	t.Run("when marshaling and unmarshaling the update shoud not change", func(t *testing.T) {
		update := mongotype.NewUpdate()
		update.Set("foo", "bar")
		update.IncrementInt("count", 1)
		buf, err := bson.Marshal(update)
		if err != nil {
			t.Error("it should not return an error")
		}
		var got mongotype.Update
		if err := bson.Unmarshal(buf, &got); err != nil {
			t.Error(err)
		}
		if update.Match(got) != nil {
			t.Errorf("expected \n%#v \nto equal \n%#v", update, got)
		}
	})
	t.Run("when unmarshling a non-update document", func(t *testing.T) {
		badUpdate := map[string]int{"$set": 5}
		buf, err := bson.Marshal(badUpdate)
		if err != nil {
			t.Fatal()
		}
		var got mongotype.Update
		if err := bson.Unmarshal(buf, &got); err == nil || err.Error() != "expected json object for key \"$set\"" {
			t.Error("it should return the correct error")
		}
	})
	t.Run("when unmarshling bad bson", func(t *testing.T) {
		var got mongotype.Update
		if err := got.SetBSON(bson.Raw{
			Kind: bson.ElementDocument,
			Data: []byte("not bson"),
		}); err == nil {
			t.Error("it should return an error")
		}
	})
}
