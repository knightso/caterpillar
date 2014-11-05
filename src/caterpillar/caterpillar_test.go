package caterpillar

import (
	"testing"
	"appengine/aetest"
	"reflect"
)

//func parseWormhole(s string) (*Wormhole, string) {
func TestParseWormhole1(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	wh, rpl := parseWormhole(`[[hoge]]`)
	if wh == nil {
		t.Errorf("wh is nil")
	}

	if !reflect.DeepEqual(*wh, Wormhole {
		"hoge",
		"",
		false,
		AREA,
	}) {
		t.Errorf("illegal wh: %#v\n", wh)
	}

	if rpl != `{{index .Areas "hoge"}}` {
		t.Errorf("illegal rpl: %s\n", rpl)
	}
}

func TestParseWormhole2(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	wh, rpl := parseWormhole(`[[.Moke  もけ]]`)
	if wh == nil {
		t.Errorf("wh is nil")
	}

	if !reflect.DeepEqual(*wh, Wormhole {
		"Moke",
		"もけ",
		true,
		PROPERTY,
	}) {
		t.Errorf("illegal wh: %#v\n", wh)
	}

	if rpl != `{{index .Properties "Moke"}}` {
		t.Errorf("illegal rpl: %s\n", rpl)
	}
}

