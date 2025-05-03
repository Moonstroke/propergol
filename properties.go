package properties

import (
	"bufio"
)

type Properties struct {
	// TODO
}

func (p *Properties) Set(key string, value string) {
	// TODO
}

func (p *Properties) Get(key string) string {
	return "" // TODO
}

func (p *Properties) Load(bufio.Reader) error {
	return nil // TODO
}

func (p *Properties) Write(bufio.Writer) error {
	return nil // TODO
}
