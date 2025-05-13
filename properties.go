package properties

import (
	"bufio"
	"errors"
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
	return &Properties{}
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
 * If there is no property with this key, nil is returned.
 */
func (p *Properties) Get(key string) string {
	return p.values[key]
}

/*
 * Load properties in texte form from the given reader.
 */
func (p *Properties) Load(*bufio.Reader) error {
	return errors.New("not implemented") // TODO
}

/*
 * Output the properties in text form to the given writer.
 */
func (p *Properties) Write(*bufio.Writer) error {
	return errors.New("not implemented") // TODO
}
