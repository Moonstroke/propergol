package properties

import (
	"bufio"
	"errors"
	"strings"
)

/*
 * This structure represents a mapping of keys to values.
 * It is intended to be used to centralize configuration data of an application.
 * The property keys and values are represented as string objects.
 */
type Properties struct {
	values map[string]string
}

/*
 * Create an empty instance of the Properties structure.
 */
func New() *Properties {
	return &Properties{make(map[string]string)}
}

/*
 * Assign the given value to the property with the specified key.
 * If no property with this key exists, it is added;
 * otherwise, the value is replaced by the one given and the former value is discarded.
 */
func (p *Properties) Set(key string, value string) {
	p.values[key] = value
}

/*
 * Retrieve the value of the property with the specified key.
 * If there is no property with this key, the empty string is returned.
 */
func (p *Properties) Get(key string) (string, bool) {
	val, present := p.values[key]
	return val, present
}

func splitLine(line string) (string, string, bool) {
	var key string
	builder := strings.Builder{}
	escaped := false
	inKey := true
	for _, c := range line {
		if c == '\\' {
			escaped = true
		} else if escaped {
			// TODO check if c is in white list ('\\', '=', others?)
			builder.WriteRune(c)
			escaped = false
		} else if c == '=' && inKey {
			// Actual separator met. Finalize the key and prepare to build the value
			key = builder.String()
			builder.Reset()
			inKey = false
		} else {
			builder.WriteRune(c)
		}
	}
	if inKey {
		// No separator found: illegal definition. Return what we can anyway
		return builder.String(), "", false
	}
	return key, builder.String(), true
}

/*
 * Parse properties in text form from the given reader.
 */
func (p *Properties) Load(reader *bufio.Reader) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		// TODO count line numbers
		line := s.Text()
		for line[len(line)-1] == '\\' {
			if !s.Scan() {
				return errors.New("invalid property definition: no continuation line")
			}
			contLine := s.Text()
			line = line[:len(line)-1] + strings.TrimLeft(contLine, " \t")
		}
		key, value, ok := splitLine(line)
		if !ok {
			return errors.New("invalid property definition: no separator")
		}
		p.Set(strings.Trim(key, " \t"), strings.Trim(value, " \t"))
	}
	return s.Err()
}

func escapeKey(key string) string {
	return key // TODO escape separators, backslashes, line breaks
}

func escapeValue(value string) string {
	return value // TODO escape backslashes, line breaks
}

/*
 * Output the properties in text form to the given writer.
 */
func (p *Properties) Store(writer *bufio.Writer) error {
	for key, val := range p.values {
		if _, e := writer.WriteString(escapeKey(key)); e != nil {
			return e
		}
		if e := writer.WriteByte('='); e != nil {
			return e
		}
		if _, e := writer.WriteString(escapeValue(val)); e != nil {
			return e
		}
		if e := writer.WriteByte('\n'); e != nil {
			return e
		}
	}
	return nil
}
