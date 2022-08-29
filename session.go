package spell

import (
	"bytes"
	"encoding/gob"
)

type Session map[string]interface{}

func NewSession() Session {
	return make(Session)
}

func (s Session) encode() (string, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(s)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (s Session) decode(data string) error {
	b := bytes.NewBufferString(data)
	dec := gob.NewDecoder(b)

	return dec.Decode(&s)
}
