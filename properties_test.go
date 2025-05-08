package properties

import (
	"testing"
)

const (
	KEY   = "key"
	VALUE = "value"
)

func TestPropertiesGetReturnsValuePassedToSet(t *testing.T) {
	prop := Properties{}
	prop.Set(KEY, VALUE)
	if prop.Get(KEY) != VALUE {
		t.Fail()
	}
}
