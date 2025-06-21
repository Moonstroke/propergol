// Package properties provides a structure that centralizes and manipulates application properties.
package properties

import (
	"bufio"
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

// Parse properties in text form from the given reader.
func (p *Properties) Load(reader io.Reader) error {
	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanRunes)
	var lineNumber uint = 1
	var key string
	builder := strings.Builder{}
	// Indicates whether the scanner is currently parsing an escape sequence
	escaped := false
	// Indicates whether the current property member (key or value) is being parsed
	// (i.e. if we are no longer scanning leading whitespace)
	inMember := false
	// Indicates whether we are parsing the key or value (i.e. the separator has been met)
	inKey := true
	for s.Scan() {
		var c rune
		// string range iterates over runes. We just want the first one
		for _, r := range s.Text() {
			c = r
			break
		}
		if escaped {
			if c == '\n' {
				// Wrapped line
				lineNumber++
				inMember = false
			} else if !(c == '\\' || inKey && c == '=') {
				return propDefError{lineNumber, "illegal escape sequence \\" + string(c)}
			} else {
				builder.WriteRune(c)
			}
			escaped = false
		} else if c == '\\' {
			escaped = true
			inMember = true
		} else if c == '\n' {
			// End of logical line
			if inKey {
				// No separator found: ill-formed definition
				return propDefError{lineNumber, "no separator"}
			}
			p.Set(strings.TrimRight(key, " \t"), strings.TrimRight(builder.String(), " \t"))
			builder.Reset()
			inKey = true
			inMember = false
		} else if c == '=' && inKey {
			if !inMember {
				return propDefError{lineNumber, "empty key"}
			}
			// Actual separator met. Finalize the key and prepare to build the value
			key = builder.String()
			builder.Reset()
			inKey = false
			inMember = false
		} else if !inMember && !inKey && c == '#' {
			// (!inMember && !inKey) <=> at the beginning of the line (index 0 or in indentation whitespace)
			for t := s.Text(); s.Scan() && t != "\n"; {
				// Consume comment line
			}
		} else if inMember || c != ' ' && c != '\t' {
			// Skip leading whitespace
			builder.WriteRune(c)
			inMember = true
		}
	}
	// Process last line if no trailing EOL was found
	if inMember {
		if inKey {
			// No separator found: ill-formed definition
			return propDefError{lineNumber, "no separator"}
		}
		p.Set(strings.TrimRight(key, " \t"), strings.TrimRight(builder.String(), " \t"))
	}
	return s.Err()
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
