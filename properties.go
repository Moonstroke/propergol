package properties

import (
	"bufio"
)

/*
 * This structure represents a mapping of keys to values.
 * It is intended to be used to centralize configuration data of an application.
 * The property keys and values are represented as string objects.
 */
type Properties struct {
	// TODO
}

/*
 * Assign the given value to the property with the specified key.
 */
func (p *Properties) Set(key string, value string) {
	// TODO
}

/*
 * Retrieve the value of the property with the specified key.
 */
func (p *Properties) Get(key string) string {
	return "" // TODO
}

/*
 * Load properties in texte form from the given reader.
 */
func (p *Properties) Load(bufio.Reader) error {
	return nil // TODO
}

/*
 * Output the properties in text form to the given writer.
 */
func (p *Properties) Write(bufio.Writer) error {
	return nil // TODO
}
