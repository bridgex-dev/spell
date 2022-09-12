package spell

type Session map[string]interface{}

func NewSession() Session {
	return make(Session)
}
