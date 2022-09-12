package spell

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
