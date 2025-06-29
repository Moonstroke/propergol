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

func processByte(c byte, p *Properties, lineNumber *uint, key *string, builder *strings.Builder, escaped, inMember, inKey, skipLine *bool) error {
	if *skipLine {
		if c == '\n' {
			*skipLine = false
		}
	} else if *escaped {
		if c == '\n' {
			// Wrapped line
			*lineNumber++
			*inMember = false
		} else if !(c == '\\' || *inKey && c == '=') {
			return propDefError{*lineNumber, "illegal escape sequence \\" + string(c)}
		} else {
			builder.WriteByte(c)
		}
		*escaped = false
	} else if c == '\\' {
		*escaped = true
		*inMember = true
	} else if c == '\n' {
		// End of physical line (escaped line breaks already handled above)
		// not in a member => blank or empty line: no property to add.
		if *inMember {
			if *inKey {
				// No separator found: ill-formed definition
				return propDefError{*lineNumber, "no separator"}
			}
			p.Set(strings.TrimRight(*key, " \t"), strings.TrimRight(builder.String(), " \t"))
			builder.Reset()
			*inKey = true
			*inMember = false
		}
	} else if c == '=' && *inKey {
		if !*inMember {
			return propDefError{*lineNumber, "empty key"}
		}
		// Actual separator met. Finalize the key and prepare to build the value
		*key = builder.String()
		builder.Reset()
		*inKey = false
		*inMember = false
	} else if !*inMember && *inKey && c == '#' {
		// (!*inMember && *inKey) <=> at the beginning of the line (index 0 or in indentation whitespace)
		*skipLine = true
	} else if *inMember || c != ' ' && c != '\t' {
		// Skip leading whitespace
		builder.WriteByte(c)
		*inMember = true
	}
	return nil
}

// Parse properties in text form from the given reader.
func (p *Properties) Load(reader io.Reader) error {
	buffer := make([]byte, 1)
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
	// Indicates whether we are currently reading a comment line (to be skipped)
	skipLine := false
	var err error
	for _, err = reader.Read(buffer); err == nil; _, err = reader.Read(buffer) {
		if err := processByte(buffer[0], p, &lineNumber, &key, &builder, &escaped, &inMember, &inKey, &skipLine); err != nil {
			return err
		}
	}
	if escaped {
		return propDefError{lineNumber, "line wrapped without a continuation"}
	}
	// Process last line if no trailing EOL was found
	if inMember {
		if inKey {
			// No separator found: ill-formed definition
			return propDefError{lineNumber, "no separator"}
		}
		p.Set(strings.TrimRight(key, " \t"), strings.TrimRight(builder.String(), " \t"))
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
