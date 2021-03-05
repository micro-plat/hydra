package mock

type Middle struct {
	FindPathV  bool
	ClearAuthV bool
	RoutePath  string
}

func (m *Middle) Next()                    {}
func (m *Middle) Find(path string) bool    { return m.FindPathV }
func (m *Middle) Service(string)           {}
func (m *Middle) ClearAuth(c ...bool) bool { return m.ClearAuthV }
func (m *Middle) GetRouterPath() string    { return m.RoutePath }
