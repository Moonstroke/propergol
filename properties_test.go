package properties

import (
	"bufio"
	"strings"
	"testing"
)

const (
	KEY   = "key"
	VALUE = "value"
	REPR  = KEY + "=" + VALUE + "\n"
)

func setUpTestInstance() *Properties {
	return New()
}

func TestPropertiesGetReturnsValuePassedToSet(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadParsesRepresentation(t *testing.T) {
	prop := setUpTestInstance()
	e := prop.Load(bufio.NewReader(strings.NewReader(REPR)))
	if e != nil {
		t.Fatal(e)
	}
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesWriteFollowsReprFormat(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	writer := strings.Builder{}
	e := prop.Store(bufio.NewWriter(&writer))
	if e != nil {
		t.Fatal(e)
	}
	if stored := writer.String(); stored != REPR {
		t.Fatal("Expected: " + REPR + "; got: " + stored)
	}
}
