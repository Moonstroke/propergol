package properties_test

import "testing"

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

func TestPropertiesStoreEscapesDQuotesInKey(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(`key"with"embedded"quotes`, VALUE)
	expected := `key\"with\"embedded\"quotes=` + VALUE
	if stored := storeToString(t, prop); stored != expected {
		t.Fatalf("Expected: %q; got: %q", expected, stored)
	}
}

func TestPropertiesStoreEscapesDQuotesInValue(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, `value"with"embedded"quotes`)
	expected := KEY + `=value\"with\"embedded\"quotes`
	if stored := storeToString(t, prop); stored != expected {
		t.Fatalf("Expected: %q; got: %q", expected, stored)
	}
}

func TestPropertiesStoreStoresKeyPrefixedWithHashSignInDQ(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set("# "+KEY, VALUE)
	repr := `"# ` + KEY + `"=` + VALUE
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
	}
}

func TestPropertiesStoreQuotesWhitespaceOnlyKey(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set("   ", VALUE)
	repr := `"   "=` + VALUE
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
	}
}

func TestPropertiesStoreQuotesWhitespaceOnlyValue(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, "   ")
	repr := KEY + `="   "`
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
	}
}

func TestPropertiesStoreQuoteKeyWSurroundingWs(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(" "+KEY+" ", VALUE)
	repr := `" ` + KEY + ` "=` + VALUE
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
	}
}

func TestPropertiesStoreQuotesValueWSurroundingWs(t *testing.T) {
	prop := setUpTestInstance()
	prop.Set(KEY, " "+VALUE+" ")
	repr := KEY + `=" ` + VALUE + ` "`
	if stored := storeToString(t, prop); stored != repr {
		t.Fatalf("Expected: %q; got %q", repr, stored)
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
