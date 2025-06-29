package properties

import (
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
		t.Fatalf("For key %s: expected value %q, got %q", key, value, got)
	}
}

func assertGetExpected(t *testing.T, prop *Properties, key, expected string) {
	if got, present := prop.Get(key); !present || got != expected {
		t.Fatalf("Expected: %q; got %q", expected, got)
	}
}

func assertGetAbsent(t *testing.T, prop *Properties, key string) {
	if _, present := prop.Get(key); present {
		t.Fatal("Expected: absent; got: present")
	}
}

func assertLoadReturnsError(t *testing.T, prop *Properties, repr string) {
	e := prop.Load(strings.NewReader(repr))
	if e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func loadFromString(t *testing.T, prop *Properties, data string) {
	e := prop.Load(strings.NewReader(data))
	if e != nil {
		t.Fatal(e)
	}
}

func storeToString(t *testing.T, prop *Properties) string {
	stringWriter := &strings.Builder{}
	e := prop.Store(stringWriter)
	if e != nil {
		t.Fatal(e)
	}
	repr := stringWriter.String()
	if len(repr) == 0 {
		return ""
	}
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
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadIgnoresLeadingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, " \t"+REPR)
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadIgnoresTrailingWhitespace(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, KEY+"="+VALUE+" \t")
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadIgnoresWhitespaceAroundSeparator(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, KEY+" = "+VALUE)
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadIgnoresEmptyLines(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, "\n\n"+REPR+"\n\n")
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadIgnoresBlankLines(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, "    \n\t  \n"+REPR+"\n\t \n  \t")
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadHandlesEscapedSeparatorInKey(t *testing.T) {
	prop := setUpTestInstance()
	key := `key with\=separator`
	loadFromString(t, prop, key+"="+VALUE)
	assertGetExpected(t, prop, "key with=separator", VALUE)
}

func TestPropertiesLoadHandlesWrappedLines(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop,
		KEY+`=value broken \
		      and indented`)
	assertGetExpected(t, prop, KEY, "value broken and indented")
}

func TestPropertiesLoadFailsOnWrappedLineWoCont(t *testing.T) {
	prop := setUpTestInstance()
	e := prop.Load(strings.NewReader(KEY + `=value broken\`))
	if e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func TestPropertiesLoadIgnoresComments(t *testing.T) {
	prop := setUpTestInstance()
	key := "# " + KEY
	loadFromString(t, prop, key+"="+VALUE)
	assertGetAbsent(t, prop, key)
}

func TestPropertiesLoadIgnoresIndentedComments(t *testing.T) {
	prop := setUpTestInstance()
	key := "# " + KEY
	loadFromString(t, prop, " \t "+key+"="+VALUE)
	assertGetAbsent(t, prop, key)
}

func TestPropertiesLoadHasNoInlineComments(t *testing.T) {
	prop := setUpTestInstance()
	value := VALUE + " # not a comment"
	loadFromString(t, prop, KEY+"="+value)
	assertGetExpected(t, prop, KEY, value)
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInKey(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, "illegal\\ escape-sequence="+VALUE)
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, KEY+"=illegal\\=escape-sequence")
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
	assertGetExpected(t, prop, key, value)
}

func TestRoundTripLoadThenStore(t *testing.T) {
	prop := setUpTestInstance()
	repr := "key:with\\=special chars\tin#it=value:with=special chars\tas#well"
	loadFromString(t, prop, repr)
	if stored := storeToString(t, prop); stored != repr {
		t.Fatal("Expected: " + repr + ", got: " + stored)
	}
}
