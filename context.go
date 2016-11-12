package antgo

type Context interface {
	SetSession(id string, data string)
	GetSession(id string) string
	SessionDecode(str string) map[string]interface{}
	SessionEncode(data map[string]interface{}) string
}
