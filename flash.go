package spell

import (
	"bytes"
	"encoding/gob"
)

type Flash map[string]string

func NewFlash() Flash {
	return make(Flash)
}

func (f Flash) Error(msg string) {
	f["error"] = msg
}

func (f Flash) Success(msg string) {
	f["success"] = msg
}

func (f Flash) encode() (string, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(f)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (f Flash) decode(data string) error {
	b := bytes.NewBufferString(data)
	dec := gob.NewDecoder(b)

	return dec.Decode(&f)
}
