package basic

import "github.com/micro-plat/lib4go/encoding/base64"

type auth struct {
	userName string
	auth     string
}

func newAuthorization(m map[string]string) []*auth {
	pairs := make([]*auth, 0, len(m))
	for user, password := range m {
		value := createAuth(user, password)
		pairs = append(pairs, &auth{
			userName: user,
			auth:     value,
		})
	}
	return pairs
}
func createAuth(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.Encode(base)
}
