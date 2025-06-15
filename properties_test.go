package properties

import (
	"bufio"
	"strings"
	"testing"
)

const (
	KEY   = "key"
	VALUE = "value"
	REPR  = KEY + "=" + VALUE
)

func setUpTestInstance() *Properties {
	return New()
}

func assertSetAndGetBackSame(t *testing.T, key, value string) {
	prop := setUpTestInstance()
	prop.Set(key, value)
	if got, present := prop.Get(key); !present || got != value {
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
	repr := stringWriter.String()
	return repr[:len(repr)-1] /* Trim trailing newline */
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
	if got, present := prop.Get(KEY); !present || got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresLeadingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, " \t"+REPR)
	if got, present := prop.Get(KEY); !present || got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresTrailingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, KEY+"="+VALUE+" \t")
	if got, present := prop.Get(KEY); !present || got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadIgnoresWhitespaceAroundSeparator(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, KEY+" = "+VALUE)
	if got, present := prop.Get(KEY); !present || got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadHandlesEscapedSeparatorInKey(t *testing.T) {
	prop := setUpTestInstance()
	key := `key with\=separator`
	loadFromString(t, prop, key+"="+VALUE)
	if got, present := prop.Get("key with=separator"); !present || got != VALUE {
		t.Fatal("Expected: " + VALUE + "; got: " + got)
	}
}

func TestPropertiesLoadHandlesWrappedLines(t *testing.T) {
	prop := setUpTestInstance()
	value := "value broken and indented"
	loadFromString(t, prop,
		KEY+`=value broken \
		      and indented`)
	if got, present := prop.Get(KEY); !present || got != value {
		t.Fatal("Expected: " + value + "; got: " + got)
	}
}

func TestPropertiesLoadFailsOnWrappedLineWoCont(t *testing.T) {
	prop := setUpTestInstance()
	e := prop.Load(bufio.NewReader(strings.NewReader(KEY + `=value broken\`)))
	if e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func TestPropertiesLoadIgnoresComments(t *testing.T) {
	prop := setUpTestInstance()
	key := "# " + KEY
	loadFromString(t, prop, key+"="+VALUE)
	if _, present := prop.Get(key); present {
		t.Fatal("Expected: absent; got: present")
	}
}

func TestPropertiesLoadHasNoInlineComments(t *testing.T) {
	prop := setUpTestInstance()
	value := VALUE + " # not a comment"
	loadFromString(t, prop, KEY+"="+value)
	if got, present := prop.Get(KEY); !present || got != value {
		t.Fatal("Expected: " + value + "; got: " + got)
	}
}

func TestPropertiesWriteFollowsReprFormat(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	if stored := storeToString(t, prop); stored != REPR {
		t.Fatal("Expected: " + REPR + "; got: " + stored)
	}
}

func TestPropertiesStoreEscapesSeparatorInKey(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set("key with=separator", VALUE)
	expected := `key with\=separator=` + VALUE
	if stored := storeToString(t, prop); stored != expected {
		t.Fatal("Expected: " + expected + "; got: " + stored)
	}
}

func TestRoundTripStoreThenLoad(t *testing.T) {
	prop := setUpTestInstance()
	key := "key:with=special chars\tin#it\\"
	value := "value:with=special chars\tas#well"
	prop.Set(key, value)
	repr := storeToString(t, prop)
	prop2 := setUpTestInstance()
	loadFromString(t, prop2, repr)
	if got, present := prop2.Get(key); !present || got != value {
		t.Fatal("Expected: " + value + ", got: " + got)
	}
}

func TestRoundTripLoadThenStore(t *testing.T) {
	prop := setUpTestInstance()
	repr := "key:with\\=special chars\tin#it=value:with=special chars\tas#well"
	loadFromString(t, prop, repr)
	if stored := storeToString(t, prop); stored != repr {
		t.Fatal("Expected: " + repr + ", got: " + stored)
	}
}
