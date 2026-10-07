// Package properties provides a structure that centralizes and manipulates application properties.
package properties

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// This structure represents a mapping of keys to values.
// It is intended to be used to centralize configuration data of an application.
// The property keys and values are represented as string objects.
type Properties struct {
	values sync.Map
}

// Create an empty instance of the Properties structure.
func New() *Properties {
	return &Properties{}
}

// Assign the given value to the property with the specified key.
// If no property with this key exists, it is added;
// otherwise, the value is replaced by the one given and the former value is discarded.
func (p *Properties) Set(key string, value string) {
	p.values.Store(key, value)
}

// Retrieve the value of the property with the specified key.
// If there is no property with this key, the empty string is returned.
func (p *Properties) Get(key string) (string, bool) {
	val, present := p.values.Load(key)
	if !present {
		return "", false
	}
	str, is_string := val.(string)
	if !is_string {
		return "", false
	}
	return str, true
}

type propDefError struct {
	lineNumber uint
	message    string
}

func (e propDefError) Error() string {
	return fmt.Sprintf("invalid property definition on line %d: %s", e.lineNumber, e.message)
}

func unescape(c byte) (byte, bool) {
	switch c {
	case '\\', '=', '"':
		return c, true
	case 'n':
		return '\n', true
	case 'r':
		return '\r', true
	case 't':
		return '\t', true
	case 'f':
		return '\f', true
	case 'v':
		return '\v', true
	case 'a':
		return '\a', true
	case 'b':
		return '\b', true
	case 'e':
		return 27, true
	case '0':
		return 0, true
	}
	return '?', false
}

// Holds data used while processing input
type loadState struct {
	lineNumber uint
	// Retains the key of the current definition (empty before the separator has been found)
	key string
	// Used to construct each property member in turn
	builder strings.Builder
	// Index of the last significant (i.e. not discardable whitespace) character in the above builder
	lastChar uint
	// Indicates whether the scanner is currently parsing an escape sequence
	escaped bool
	// Indicates whether the current property member (key or value) is being parsed
	// (i.e. if we are no longer scanning leading whitespace)
	inMember bool
	// Indicates whether we are parsing the key or value (i.e. the separator has been met)
	inKey bool
	// Indicates whether we are currently reading a comment line (to be skipped)
	skipLine bool
	// Indicats that the previous character was a Carriage Return
	wasCR bool
}

func processByte(c byte, p *Properties, state *loadState) error {
	// Does not fit in the state switch because c still needs to be processed after an optional stray CR is handled
	if state.wasCR && c != '\n' {
		state.builder.WriteByte('\r')
		state.wasCR = false
	}
	switch {
	case state.skipLine:
		if c == '\n' {
			state.skipLine = false
		}
	case state.escaped:
		if c == '\r' {
			state.wasCR = true
		} else {
			if c == '\n' {
				// Wrapped line
				state.lineNumber++
				state.inMember = false
				// Reset CRLF sequence flag
				state.wasCR = false
			} else {
				u, ok := unescape(c)
				if !ok {
					return propDefError{state.lineNumber, "illegal escape sequence \\" + string(c)}
				}
				state.builder.WriteByte(u)
				state.lastChar++
			}
			state.escaped = false
		}
	case c == '\\':
		state.escaped = true
		state.inMember = true
	case c == '\r':
		state.wasCR = true
	case c == '\n':
		// End of physical line (escaped line breaks already handled above)
		// not in a member => blank or empty line: no property to add.
		if state.inMember {
			if state.inKey {
				// No separator found: ill-formed definition
				return propDefError{state.lineNumber, "no separator"}
			}
			p.Set(state.key, state.builder.String()[:state.lastChar])
			state.builder.Reset()
			state.inKey = true
			state.inMember = false
			state.lastChar = 0
		}
		// Reset CRLF sequence flag
		state.wasCR = false
	case c == '=' && state.inKey:
		if !state.inMember {
			return propDefError{state.lineNumber, "empty key"}
		}
		// Actual separator met. Finalize the key and prepare to build the value
		state.key = state.builder.String()[:state.lastChar]
		state.builder.Reset()
		state.inKey = false
		state.inMember = false
		state.lastChar = 0
	case !state.inMember && state.inKey && c == '#':
		// (!state.inMember && state.inKey) <=> at the beginning of the line (index 0 or in indentation whitespace)
		state.skipLine = true
	case c == ' ' || c == '\t':
		// Only write significant whitespace (i.e. not leading indentation)
		if state.inMember {
			state.builder.WriteByte(c)
		}
	default:
		if c != '"' {
			state.builder.WriteByte(c)
		}
		state.inMember = true
		state.lastChar = uint(state.builder.Len())
	}
	return nil
}

func pullBytes(p *Properties, state *loadState, byteCh <-chan byte, errCh chan<- error) {
	defer close(errCh)
	for c := range byteCh {
		if err := processByte(c, p, state); err != nil {
			errCh <- err
			return
		}
	}
}

func pushBytes(reader io.Reader, buffer []byte, byteCh chan<- byte, errCh <-chan error) (error, bool) {
	defer close(byteCh)
	var err error
	var n int
	for err == nil {
		n, err = reader.Read(buffer)
		for _, c := range buffer[:n] {
			byteCh <- c
			select {
			case processErr := <-errCh:
				return processErr, true
			default:
				// No error, continue
			}
		}
	}
	if err != io.EOF {
		return err, false
	}
	return nil, false
}

// Parse properties in text form from the given reader.
func (p *Properties) Load(reader io.Reader) error {
	buffer := make([]byte, 1024)
	byteCh := make(chan byte, 1024)
	state := loadState{
		lineNumber: 1,
		inKey:      true,
	}
	errCh := make(chan error, 1)
	go pullBytes(p, &state, byteCh, errCh)
	var err error
	var exitNow bool
	if err, exitNow = pushBytes(reader, buffer, byteCh, errCh); exitNow {
		return err
	}
	if processErr := <-errCh; processErr != nil {
		return processErr
	}
	if state.escaped {
		return propDefError{state.lineNumber, "line wrapped without a continuation"}
	}
	// Process last line if no trailing EOL was found
	if state.inMember {
		if state.inKey {
			// No separator found: ill-formed definition
			return propDefError{state.lineNumber, "no separator"}
		}
		if state.wasCR {
			state.builder.WriteByte('\r')
			state.lastChar++
		}
		p.Set(state.key, state.builder.String()[:state.lastChar])
	}
	return err
}

var keyEscaper, valueEscaper *strings.Replacer
var replacerInit sync.Once

// Output the properties in text form to the given writer.
func (p *Properties) Store(writer io.Writer) error {
	replacerInit.Do(func() {
		oldnew := []string{
			"=", `\=`,
			`\`, `\\`,
			"\n", `\n`,
			"\r", `\r`,
			"\f", `\f`,
			"\v", `\v`,
			"\a", `\a`,
			"\b", `\b`,
			"\x1b", `\e`,
			"\000", `\0`,
			`"`, `\"`,
		}
		keyEscaper = strings.NewReplacer(oldnew...)
		valueEscaper = strings.NewReplacer(oldnew[2:]...) /* Skip escaping of = as it has no special meaning in the value */
	})
	var e error = nil
	p.values.Range(func(key, val any) bool {
		key_str, key_is_string := key.(string)
		if !key_is_string {
			e = fmt.Errorf("Expected string, got %[1]T: %[1]v", key)
			return false
		}
		keyNeedsQuoting := key_str[0] == '#' || key_str[0] == ' ' || key_str[0] == '\t' || key_str[len(key_str)-1] == ' ' || key_str[len(key_str)-1] == '\t'
		if keyNeedsQuoting {
			if _, e = writer.Write([]byte{'"'}); e != nil {
				return false
			}
		}
		if _, e = keyEscaper.WriteString(writer, key_str); e != nil {
			return false
		}
		if keyNeedsQuoting {
			if _, e = writer.Write([]byte{'"'}); e != nil {
				return false
			}
		}
		if _, e = writer.Write([]byte{'='}); e != nil {
			return false
		}
		val_str, val_is_string := val.(string)
		if !val_is_string {
			e = fmt.Errorf("Expected string, got %[1]T: %[1]v", val)
			return false
		}
		valNeedsQuoting := val_str[0] == ' ' || val_str[0] == '\t' || val_str[len(val_str)-1] == ' ' || val_str[len(val_str)-1] == '\t'
		if valNeedsQuoting {
			if _, e = writer.Write([]byte{'"'}); e != nil {
				return false
			}
		}
		if _, e = valueEscaper.WriteString(writer, val_str); e != nil {
			return false
		}
		if valNeedsQuoting {
			if _, e = writer.Write([]byte{'"'}); e != nil {
				return false
			}
		}
		if _, e = writer.Write([]byte{'\n'}); e != nil {
			return false
		}
		return true
	})
	return e
}
