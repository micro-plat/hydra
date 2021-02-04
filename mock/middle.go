package mock

type Middle struct {
	FindPathV  bool
	ClearAuthV bool
}

func (m *Middle) Next()                    {}
func (m *Middle) Find(path string) bool    { return m.FindPathV }
func (m *Middle) Service(string)           {}
func (m *Middle) ClearAuth(c ...bool) bool { return m.ClearAuthV }
