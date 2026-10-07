package properties_test

import (
	"errors"
	"strings"
	"testing"

	properties "github.com/Moonstroke/propergol"
)

const (
	KEY   = "key"
	VALUE = "value"
	REPR  = KEY + "=" + VALUE
)

func setUpTestInstance() *properties.Properties {
	return properties.New()
}

func assertSetAndGetBackSame(t *testing.T, key, value string) {
	t.Helper()
	prop := setUpTestInstance()
	prop.Set(key, value)
	if got, present := prop.Get(key); !present || got != value {
		t.Fatalf("For key %s: expected value %q, got %q", key, value, got)
	}
}

func assertGetExpected(t *testing.T, prop *properties.Properties, key, expected string) {
	t.Helper()
	got, present := prop.Get(key)
	if !present {
		t.Fatalf("Expected: %q; got absent", expected)
	} else if got != expected {
		t.Fatalf("Expected: %q; got %q", expected, got)
	}
}

func assertGetAbsent(t *testing.T, prop *properties.Properties, key string) {
	t.Helper()
	if _, present := prop.Get(key); present {
		t.Fatal("Expected: absent; got: present")
	}
}

func assertLoadReturnsError(t *testing.T, prop *properties.Properties, repr string) {
	t.Helper()
	if e := prop.Load(strings.NewReader(repr)); e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
}

func loadFromString(t *testing.T, prop *properties.Properties, data string) {
	t.Helper()
	if e := prop.Load(strings.NewReader(data)); e != nil {
		t.Fatal(e)
	}
}

func storeToString(t *testing.T, prop *properties.Properties) string {
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

var TEST_ERROR error = errors.New("test error")

/* FailingReaderWriter is an implementation of io.Reader and io.Writer that always fail. */
type failingReaderWriter struct{}

func (_ failingReaderWriter) Read(b []byte) (int, error) {
	return 0, TEST_ERROR
}

func (_ failingReaderWriter) Write(p []byte) (n int, err error) {
	return 0, TEST_ERROR
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
