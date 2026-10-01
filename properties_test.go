package properties

import (
	"errors"
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
	t.Helper()
	prop := setUpTestInstance()
	prop.Set(key, value)
	if got, present := prop.Get(key); !present || got != value {
		t.Fatalf("For key %s: expected value %q, got %q", key, value, got)
	}
}

func assertGetExpected(t *testing.T, prop *Properties, key, expected string) {
	t.Helper()
	got, present := prop.Get(key)
	if !present {
		t.Fatalf("Expected: %q; got absent", expected)
	} else if got != expected {
		t.Fatalf("Expected: %q; got %q", expected, got)
	}
}

func assertGetAbsent(t *testing.T, prop *Properties, key string) {
	t.Helper()
	if _, present := prop.Get(key); present {
		t.Fatal("Expected: absent; got: present")
	}
}

func assertLoadReturnsError(t *testing.T, prop *Properties, repr string) {
	t.Helper()
	if e := prop.Load(strings.NewReader(repr)); e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func loadFromString(t *testing.T, prop *Properties, data string) {
	t.Helper()
	if e := prop.Load(strings.NewReader(data)); e != nil {
		t.Fatal(e)
	}
}

func storeToString(t *testing.T, prop *Properties) string {
	t.Helper()
	stringWriter := &strings.Builder{}
	if e := prop.Store(stringWriter); e != nil {
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
	rawKey := `key with\=separator`
	processedKey := `key with=separator`
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadAcceptsEscapedSeparatorInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\=separator`
	processedValue := "value with=separator"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
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
	if e := prop.Load(strings.NewReader(KEY + `=value broken\`)); e == nil {
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

func TestPropertiesLoadHandlesEscapedBackslashInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\\escaped backslash`
	processedKey := "key with\\escaped backslash"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedBackslashInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\\escaped backslash`
	processedValue := "value with\\escaped backslash"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedLFInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\nescaped LF`
	processedKey := "key with\nescaped LF"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedLFInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\nescaped LF`
	processedValue := "value with\nescaped LF"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedCRInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\rescaped CR`
	processedKey := "key with\rescaped CR"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedCRInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\rescaped CR`
	processedValue := "value with\rescaped CR"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedTabInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\tescaped Tab`
	processedKey := "key with\tescaped Tab"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedTabInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\tescaped Tab`
	processedValue := "value with\tescaped Tab"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedVTabInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\vescaped VTab`
	processedKey := "key with\vescaped VTab"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedVTabInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\vescaped VTab`
	processedValue := "value with\vescaped VTab"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedFFInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\fescaped FF`
	processedKey := "key with\fescaped FF"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedFFInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\fescaped FF`
	processedValue := "value with\fescaped FF"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedNulInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\fescaped FF`
	processedKey := "key with\fescaped FF"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedNulInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\0escaped NUL`
	processedValue := "value with\000escaped NUL"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedBelInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\aescaped BEL`
	processedKey := "key with\aescaped BEL"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedBelInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\aescaped BEL`
	processedValue := "value with\aescaped BEL"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedBSInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\bescaped BS`
	processedKey := "key with\bescaped BS"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedBSInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\bescaped BS`
	processedValue := "value with\bescaped BS"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadHandlesEscapedEscInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\eescaped Esc`
	processedKey := "key with\x1bescaped Esc"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedEscInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\eescaped Esc`
	processedValue := "value with\x1bescaped Esc"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInKey(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, "illegal\\ escape-sequence="+VALUE)
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, KEY+"=illegal\\ escape-sequence")
}

func TestPropertiesStoreFollowsReprFormat(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	if stored := storeToString(t, prop); stored != REPR {
		t.Fatalf("Expected: %q; got: %q", REPR, stored)
	}
}

func TestPropertiesStoreEscapesSeparatorInKey(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set("key with=separator", VALUE)
	expected := `key with\=separator=` + VALUE
	if stored := storeToString(t, prop); stored != expected {
		t.Fatalf("Expected: %q; got: %q", REPR, stored)
	}
}

func TestPropertiesStoreCannotStoreKeyPrefixedWithHashSign(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set("# key", VALUE)
	e := prop.Store(&strings.Builder{})
	if e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func TestRoundTripStoreThenLoad(t *testing.T) {
	prop := setUpTestInstance()
	key := "\bkey:with=special\fchars\tin#it\\\000"
	value := "\avalue\nwith=special\rchars\vas#well\x1b"
	prop.Set(key, value)
	repr := storeToString(t, prop)
	prop2 := setUpTestInstance()
	loadFromString(t, prop2, repr)
	assertGetExpected(t, prop, key, value)
}

func TestRoundTripLoadThenStore(t *testing.T) {
	prop := setUpTestInstance()
	repr := `\bkey:with\=special\fchars` + "\t" + `in#it\\\0=\avalue\nwith=special\rchars\vas#well\e`
	loadFromString(t, prop, repr)
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
	}
}

func TestRoundTripLoadThenStoreDoesNotPreserveAllEscSeqs(t *testing.T) {
	prop := setUpTestInstance()
	rawRepr := `key with a discaradable\tescape sequence=value with\=discardable\tescape sequences`
	processedRepr := "key with a discaradable\tescape sequence=value with=discardable\tescape sequences"
	loadFromString(t, prop, rawRepr)
	if stored := storeToString(t, prop); stored != processedRepr {
		t.Fatalf("Expected: %q; got %q", processedRepr, stored)
	}
}

func TestPropertiesLoadHandlesCRLFLineEnding(t *testing.T) {
	prop := setUpTestInstance()
	repr := "key1 = value1\r\nkey2 = value2\r\n"
	loadFromString(t, prop, repr)
	assertGetExpected(t, prop, "key1", "value1")
	assertGetExpected(t, prop, "key2", "value2")
}

func TestPropertiesLoadHandlesCRLFInLineContinuation(t *testing.T) {
	prop := setUpTestInstance()
	repr := "key1 = value1 \\\r\n       continued"
	loadFromString(t, prop, repr)
	assertGetExpected(t, prop, "key1", "value1 continued")
}

func TestPropertiesLoadPreservesTrailingCR(t *testing.T) {
	prop := setUpTestInstance()
	repr := "key1 = value1\r"
	loadFromString(t, prop, repr)
	assertGetExpected(t, prop, "key1", "value1\r")
}

var TEST_ERROR error = errors.New("test error")

// An implementation of io.Reader and io.Writer that always fail.
type failingReaderWriter struct{}

func (_ failingReaderWriter) Read(b []byte) (int, error) {
	return 0, TEST_ERROR
}

func (_ failingReaderWriter) Write(p []byte) (n int, err error) {
	return 0, TEST_ERROR
}

type partialFailingReader struct {
	text string
}

func (pfr partialFailingReader) Read(b []byte) (int, error) {
	// Trash implementation (discards the part of pfr.text that doesn't go into b)
	// but good enough for the present case
	return copy(b, pfr.text), TEST_ERROR
}

func TestPropertiesLoadHandlesReadError(t *testing.T) {
	prop := setUpTestInstance()
	err := prop.Load(failingReaderWriter{})
	if err != TEST_ERROR {
		t.Fatalf("Expected error %v, got %v", TEST_ERROR, err)
	}
}

func TestEmptyPropertiesStoreRaisesNoError(t *testing.T) {
	prop := setUpTestInstance()
	err := prop.Store(failingReaderWriter{})
	if err != nil {
		t.Fatalf("Expected error %v, got %v", nil, err)
	}
}

func TestPropertiesStoreHandlesWriteError(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, VALUE)
	err := prop.Store(failingReaderWriter{})
	if err != TEST_ERROR {
		t.Fatalf("Expected error %v, got %v", TEST_ERROR, err)
	}
}
