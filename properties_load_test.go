package properties_test

import (
	"fmt"
	"strings"
	"testing"

	properties "github.com/Moonstroke/propergol"
)

func assertLoadReturnsError(t *testing.T, prop *properties.Properties, repr string) error {
	t.Helper()
	var e error
	if e = prop.Load(strings.NewReader(repr)); e == nil {
		t.Fatal("Expected failure, but no error was raised")
	}
	return e
}

func assertLoadReturnsErrorContaining(t *testing.T, prop *properties.Properties, repr, excerpt string) {
	t.Helper()
	msg := assertLoadReturnsError(t, prop, repr).Error()
	if !strings.Contains(msg, excerpt) {
		t.Fatalf("%q not found in %q", excerpt, msg)
	}
}

func assertLoadReturnsErrorWithLineNum(t *testing.T, prop *properties.Properties, repr string, line uint) {
	t.Helper()
	msg := assertLoadReturnsError(t, prop, repr).Error()
	if !strings.Contains(msg, "line "+fmt.Sprint(line)) {
		t.Fatalf("Line number %d not found in error message %q", line, msg)
	}
}

/* PartialFailingReader is an implementation of io.Reader that fails after successfully reading the specified text.
 * It is basically a strings.Reader that fails with TEST_ERROR instead of io.EOF. */
type partialFailingReader struct {
	text string
}

func (pfr partialFailingReader) Read(b []byte) (int, error) {
	if len(b) >= len(pfr.text) {
		/* The whole text (or what's left of it) fits into the buffer */
		return copy(b, pfr.text), TEST_ERROR
	}
	/* Blit what fits and retain what did not for the next Read */
	n := copy(b, pfr.text)
	pfr.text = pfr.text[n+1:]
	return n, nil
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

func TestPropertiesLoadHandlesEscapedDQInKey(t *testing.T) {
	prop := setUpTestInstance()
	rawKey := `key with\"escaped double quotes`
	processedKey := "key with\"escaped double quotes"
	loadFromString(t, prop, rawKey+"="+VALUE)
	assertGetExpected(t, prop, processedKey, VALUE)
}

func TestPropertiesLoadHandlesEscapedDQInValue(t *testing.T) {
	prop := setUpTestInstance()
	rawValue := `value with\"escaped double quotes`
	processedValue := "value with\"escaped double quotes"
	loadFromString(t, prop, KEY+"="+rawValue)
	assertGetExpected(t, prop, KEY, processedValue)
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInKey(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, "illegal\\ escape-sequence="+VALUE)
}

func TestPropertiesLoadForbidsLegalEscapeSequenceWithCRInValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorContaining(t, prop, KEY+"=legal\\\rnescape-sequence", "\\ + CR")
}

func TestPropertiesLoadForbidsLegalEscapeSequenceWithCRInKey(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorContaining(t, prop, "legal\\\rnescape-sequence="+VALUE, "\\ + CR")
}

func TestPropertiesLoadForbidsIllegalEscapeSequenceWithCRInValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorContaining(t, prop, KEY+"=illegal\\\r escape-sequence", "\\ + CR")
}

func TestPropertiesLoadForbidsIllegalEscapeSequenceWithCRInKey(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorContaining(t, prop, "illegal\\\r escape-sequence="+VALUE, "\\ + CR")
}

func TestPropertiesLoadForbidsIllegalEscapeSequencesInValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsError(t, prop, KEY+"=illegal\\ escape-sequence")
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine1NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY, 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine1WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine2NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine2WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine3AfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine3AfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "key with\\\ncontinuation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "key with\\\ncontinuation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorAfterCmtContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey with\\\ncontinuation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorAfterCmtContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey with\\\ncontinuation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey with\\\ncontinuation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey with\\\ncontinuation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey with\\\ncontinuation line", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_NoSeparatorLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey with\\\ncontinuation line\n", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine1NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "=value", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine1WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "=value\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n=value2", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n=value2\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine2NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n=value2", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine2WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n=value2\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine3AfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\n=value2", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine3AfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\n=value2\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "=value with\\\ncontinuation line", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "=value with\\\ncontinuation line\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyAfterCmtContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n=value with\\\ncontinuation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyAfterCmtContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n=value with\\\ncontinuation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n=value with\\\ncontinuation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n=value with\\\ncontinuation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\n=value with\\\ncontinuation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_EmptyKeyLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\n=value with\\\ncontinuation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine1NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "ke\\y="+VALUE, 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine1WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "ke\\y="+VALUE+"\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nke\\y2="+VALUE, 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nke\\y2="+VALUE+"\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine2NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nke\\y2="+VALUE, 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine2WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nke\\y2="+VALUE+"\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine3AfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nke\\y2="+VALUE, 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine3AfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nke\\y2="+VALUE+"\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "ke\\y with\\\ncontinuation line="+VALUE, 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "ke\\y with\\\ncontinuation line="+VALUE+"\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyAfterCmtContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nke\\y with\\\ncontinuation line="+VALUE, 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyAfterCmtContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nke\\y with\\\ncontinuation line="+VALUE+"\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nke\\y with\\\ncontinuation line="+VALUE, 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nke\\y with\\\ncontinuation line="+VALUE+"\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nke\\y with\\\ncontinuation line="+VALUE, 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nke\\y with\\\ncontinuation line="+VALUE+"\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "key with\\\ncontin\\uation line="+VALUE, 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "key with\\\ncontin\\uation line="+VALUE+"\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey with\\\ncontin\\uation line="+VALUE, 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContAfterCmtContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey with\\\ncontin\\uation line"+VALUE+"\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey with\\\ncontin\\uation line="+VALUE, 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey with\\\ncontin\\uation line"+VALUE+"\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey with\\\ncontin\\uation line="+VALUE, 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInKeyContLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey with\\\ncontin\\uation line"+VALUE+"\n", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine1NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=val\\ue", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine1WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=val\\ue\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2=val\\ue", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInVaAfterWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2=val\\ue\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine2NoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2=val\\ue", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine2WithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2=val\\ue\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine3AfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2=val\\ue", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine3AfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2=val\\ue\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=val\\ue with\\\ncontinuation line", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=val\\ue with\\\ncontinuation line\n", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValAfterCmtContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n"+KEY+"=val\\ue with\\\ncontinuation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValAfterCmtContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\n"+KEY+"=val\\ue with\\\ncontinuation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey2=val\\ue with\\\ncontinuation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey2=val\\ue with\\\ncontinuation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey2=val\\ue with\\\ncontinuation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey2=val\\ue with\\\ncontinuation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine1ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=value with\\\ncontin\\uation line", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine1ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=value with\\\ncontin\\uation line\n", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine2ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey2=value with\\\ncontin\\uation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine2ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey2=value with\\\ncontin\\uation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine2ContNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey2=value with\\\ncontin\\uation line", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine2ContWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey2=value with\\\ncontin\\uation line\n", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine3ContAfterCmtNoLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey2=value with\\\ncontin\\uation line", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_InvalidEscSeqInValContLine3ContAfterCmtWithLineBreak(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey2=value with\\\ncontin\\uation line\n", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine1(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=value\\", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine1NoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"\\", 1)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscAfterCmt(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2=value\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscAfterCmtNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nKEY2\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2=value\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2NoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nKEY2\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine3AfterCmt(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2=value\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine3AfterCmtNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nKEY2\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine1Cont(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, KEY+"=value with\\\ncontinuation line\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine1ContNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "key with\\\ncontinuation line\\", 2)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2ContAfterCmt(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey=value with\\\ncontinuation line\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2ContAfterCmtNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, "# comment\nkey with\\\ncontinuation line\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2Cont(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey=value with\\\ncontinuation line\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine2ContNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\nkey with\\\ncontinuation line\\", 3)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine3ContAfterCmt(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey=value with\\\ncontinuation line\\", 4)
}

func TestPropertiesLoadDisplaysLineNumberInErrorMsg_TrailingEscLine3ContAfterNoValue(t *testing.T) {
	prop := setUpTestInstance()
	assertLoadReturnsErrorWithLineNum(t, prop, REPR+"\n# comment\nkey with\\\ncontinuation line\\", 4)
}

func TestPropertiesLoadStripsQuotesAroundQuotedKeyWLeadingHash(t *testing.T) {
	prop := setUpTestInstance()
	key := "# " + KEY
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetAbsent(t, prop, `"`+key+`"`)
}

func TestPropertiesLoadPreservesQuotedKeyWLeadingHash(t *testing.T) {
	prop := setUpTestInstance()
	key := "# " + KEY
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetExpected(t, prop, key, VALUE)
}

func TestPropertiesLoadStripsQuotesAroundQuotedWhitespaceOnlyKey(t *testing.T) {
	prop := setUpTestInstance()
	key := "   "
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetAbsent(t, prop, `"`+key+`"`)
}

func TestPropertiesLoadPreservesQuotedWhitespaceOnlyKey(t *testing.T) {
	prop := setUpTestInstance()
	key := "   "
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetExpected(t, prop, key, VALUE)
}

func TestPropertiesLoadPreservesQuotedWhitespaceOnlyValue(t *testing.T) {
	prop := setUpTestInstance()
	value := "   "
	loadFromString(t, prop, KEY+`="`+value+`"`)
	assertGetExpected(t, prop, KEY, value)
}

func TestPropertiesLoadStripsQuotesAroundQuotedKeyWSurroundingWS(t *testing.T) {
	prop := setUpTestInstance()
	key := " " + KEY + " "
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetAbsent(t, prop, `"`+key+`"`)
}

func TestPropertiesLoadPreservesQuotedKeyWSurroundingWS(t *testing.T) {
	prop := setUpTestInstance()
	key := " " + KEY + " "
	loadFromString(t, prop, `"`+key+`"=`+VALUE)
	assertGetExpected(t, prop, key, VALUE)
}

func TestPropertiesLoadPreservesValueWSurroundingWS(t *testing.T) {
	prop := setUpTestInstance()
	value := " " + VALUE + " "
	loadFromString(t, prop, KEY+`="`+value+`"`)
	assertGetExpected(t, prop, KEY, value)
}

func TestPropertiesLoadDiscardsWSOutOfQuotedKey(t *testing.T) {
	prop := setUpTestInstance()
	loadFromString(t, prop, "\t\" "+KEY+` " =`+VALUE)
	assertGetAbsent(t, prop, "\t "+KEY+"  ")
}

func TestPropertiesLoadDiscardsNonQuotedWSAroundKey(t *testing.T) {
	prop := setUpTestInstance()
	key := " " + KEY + " "
	loadFromString(t, prop, "\t\""+key+`" =`+VALUE)
	assertGetExpected(t, prop, key, VALUE)
}

func TestPropertiesLoadDiscardsNonQuotedWSAroundValue(t *testing.T) {
	prop := setUpTestInstance()
	value := " " + VALUE + " "
	loadFromString(t, prop, KEY+`= "`+value+`"  `)
	assertGetExpected(t, prop, KEY, value)
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

func TestPropertiesLoadHandlesEmptyInput(t *testing.T) {
	prop := setUpTestInstance()
	t.Helper()
	/* An empty reader is basically equivalent to a failingReaderWriter returning io.EOF */
	if err := prop.Load(strings.NewReader("")); err != nil {
		t.Fatal(err)
	}
}

func TestPropertiesLoadHandlesReadError(t *testing.T) {
	prop := setUpTestInstance()
	err := prop.Load(failingReaderWriter{})
	if err != TEST_ERROR {
		t.Fatalf("Expected error %v, got %v", TEST_ERROR, err)
	}
}

func TestPropertiesLoadHandlesErrorAfterPartialContent(t *testing.T) {
	prop := setUpTestInstance()
	err := prop.Load(partialFailingReader{REPR})
	if err != TEST_ERROR {
		t.Fatalf("Expected error %v, got %v", TEST_ERROR, err)
	}
	assertGetExpected(t, prop, KEY, VALUE)
}

func TestPropertiesLoadProcessesPartialContent(t *testing.T) {
	prop := setUpTestInstance()
	prop.Load(partialFailingReader{REPR})
	assertGetExpected(t, prop, KEY, VALUE)
}
