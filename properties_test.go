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

func assertSetAndGetBack(t *testing.T, prop *Properties, key, value string) {
	prop.Set(key, value)
	if got := prop.Get(key); got != value {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func loadFromString(t *testing.T, prop *Properties) {
	e := prop.Load(bufio.NewReader(strings.NewReader(REPR)))
	if e != nil {
		t.Fatal(e)
	}
}

func storeToString(t *testing.T, prop *Properties) string {
	stringWriter := strings.Builder{}
	bufWriter := bufio.NewWriter(&stringWriter)
	e := prop.Store(bufWriter)
	if e != nil {
		t.Fatal(e)
	}
	/* Ensure that the text is passed down to the string writer */
	bufWriter.Flush()
	return stringWriter.String()
}

func TestPropertiesGetReturnsValuePassedToSet(t *testing.T) {
	prop := setUpTestInstance()
	assertSetAndGetBack(t, prop, KEY, VALUE)
}

func TestPropertiesLoadParsesRepresentation(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop)
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesWriteFollowsReprFormat(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	if stored := storeToString(t, prop); stored != REPR {
		t.Fatal("Expected: " + REPR + "; got: " + stored)
	}
}
