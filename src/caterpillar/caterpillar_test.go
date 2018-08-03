package caterpillar

import (
	"reflect"
	"testing"

	"google.golang.org/appengine/aetest"
)

//func parseWormhole(s string) (*Wormhole, string) {
func TestParseWormhole1(t *testing.T) {
	_, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	wh, rpl := parseWormhole(`[[hoge]]`)
	if wh == nil {
		t.Errorf("wh is nil")
	}

	if !reflect.DeepEqual(*wh, Wormhole{
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
	_, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	wh, rpl := parseWormhole(`[[.Moke  もけ]]`)
	if wh == nil {
		t.Errorf("wh is nil")
	}

	if !reflect.DeepEqual(*wh, Wormhole{
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
