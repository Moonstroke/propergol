// Package properties provides a structure that centralizes and manipulates application properties.
package properties

import (
	"fmt"
	"io"
	"strings"
)

// This structure represents a mapping of keys to values.
// It is intended to be used to centralize configuration data of an application.
// The property keys and values are represented as string objects.
type Properties struct {
	values map[string]string
}

// Create an empty instance of the Properties structure.
func New() *Properties {
	return &Properties{make(map[string]string)}
}

// Assign the given value to the property with the specified key.
// If no property with this key exists, it is added;
// otherwise, the value is replaced by the one given and the former value is discarded.
func (p *Properties) Set(key string, value string) {
	p.values[key] = value
}

// Retrieve the value of the property with the specified key.
// If there is no property with this key, the empty string is returned.
func (p *Properties) Get(key string) (string, bool) {
	val, present := p.values[key]
	return val, present
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
	case '\\', '=':
		return c, true
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
	// Indicates whether the scanner is currently parsing an escape sequence
	escaped bool
	// Indicates whether the current property member (key or value) is being parsed
	// (i.e. if we are no longer scanning leading whitespace)
	inMember bool
	// Indicates whether we are parsing the key or value (i.e. the separator has been met)
	inKey bool
	// Indicates whether we are currently reading a comment line (to be skipped)
	skipLine bool
}

func processByte(c byte, p *Properties, state *loadState) error {
	switch {
	case state.skipLine:
		if c == '\n' {
			state.skipLine = false
		}
	case state.escaped:
		if c == '\n' {
			// Wrapped line
			state.lineNumber++
			state.inMember = false
		} else {
			u, ok := unescape(c)
			if !ok {
				return propDefError{state.lineNumber, "illegal escape sequence \\" + string(c)}
			}
			state.builder.WriteByte(u)
		}
		state.escaped = false
	case c == '\\':
		state.escaped = true
		state.inMember = true
	case c == '\n':
		// End of physical line (escaped line breaks already handled above)
		// not in a member => blank or empty line: no property to add.
		if state.inMember {
			if state.inKey {
				// No separator found: ill-formed definition
				return propDefError{state.lineNumber, "no separator"}
			}
			p.Set(strings.TrimRight(state.key, " \t"), strings.TrimRight(state.builder.String(), " \t"))
			state.builder.Reset()
			state.inKey = true
			state.inMember = false
		}
	case c == '=' && state.inKey:
		if !state.inMember {
			return propDefError{state.lineNumber, "empty key"}
		}
		// Actual separator met. Finalize the key and prepare to build the value
		state.key = state.builder.String()
		state.builder.Reset()
		state.inKey = false
		state.inMember = false
	case !state.inMember && state.inKey && c == '#':
		// (!state.inMember && state.inKey) <=> at the beginning of the line (index 0 or in indentation whitespace)
		state.skipLine = true
	case state.inMember || c != ' ' && c != '\t':
		// Skip leading whitespace
		state.builder.WriteByte(c)
		state.inMember = true
	}
	return nil
}

// Parse properties in text form from the given reader.
func (p *Properties) Load(reader io.Reader) error {
	buffer := make([]byte, 1)
	state := loadState{
		lineNumber: 1,
		inKey:      true,
	}
	var err error
	for _, err = reader.Read(buffer); err == nil; _, err = reader.Read(buffer) {
		if err := processByte(buffer[0], p, &state); err != nil {
			return err
		}
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
		p.Set(strings.TrimRight(state.key, " \t"), strings.TrimRight(state.builder.String(), " \t"))
	}
	if err == io.EOF {
		return nil
	}
	return err
}

// Output the properties in text form to the given writer.
func (p *Properties) Store(writer io.Writer) error {
	keyEscaper := strings.NewReplacer("=", "\\=", "\\", "\\\\", "\n", "\\\n")
	valueEscaper := strings.NewReplacer("\\", "\\\\", "\n", "\\\n")
	for key, val := range p.values {
		if _, e := keyEscaper.WriteString(writer, key); e != nil {
			return e
		}
		if _, e := writer.Write([]byte{'='}); e != nil {
			return e
		}
		if _, e := valueEscaper.WriteString(writer, val); e != nil {
			return e
		}
		if _, e := writer.Write([]byte{'\n'}); e != nil {
			return e
		}
	}
	return nil
}
