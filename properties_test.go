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

func assertSetAndGetBackSame(t *testing.T, key, value string) {
	prop := setUpTestInstance()
	prop.Set(key, value)
	if got := prop.Get(key); got != value {
		t.Fatal("For key " + key + `: expected value "` + value + `", got "` + got + `"`)
	}
}

func loadFromString(t *testing.T, prop *Properties, data string) {
	e := prop.Load(bufio.NewReader(strings.NewReader(data)))
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
	assertSetAndGetBackSame(t, KEY, VALUE)
}

func TestPropertiesAcceptKeysWithSpaces(t *testing.T) {
	assertSetAndGetBackSame(t, "a key with spaces", "whatever")
}

func TestPropertiesAcceptValuesWithSpaces(t *testing.T) {
	assertSetAndGetBackSame(t, "whatever", "a value with spaces")
}

func TestPropertiesAcceptValuesWithColons(t *testing.T) {
	assertSetAndGetBackSame(t, "whatever", "a:value:with:colons")
}

func TestPropertiesAcceptValuesWithSeparators(t *testing.T) {
	assertSetAndGetBackSame(t, "whatever", "a=value=with=separators")
}

func TestPropertiesLoadParsesRepresentation(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, REPR)
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresLeadingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, " \t"+REPR)
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresTrailingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, REPR+" \t")
	if got := prop.Get(KEY); got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresWhitespaceAroundSeparator(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, KEY+" = "+VALUE)
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
