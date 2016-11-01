package chatroom

type Handler struct {
	worker *Worker
}

func (*Handler) GetSecret() string {
    return ""
}

func (*Handler) SetSecret(secret string) {

}

func (*Handler) GetSession(end_id string) string {
    return ""
}

func (*Handler) SetSession(end_id string, session string) {

}
func (*Handler) InitSession(end_id string, session string) {

}

func (*Handler) SendToEnd(end_id string, message string) {

}

func (*Handler) SendToUid(uid string, message string) {

}

func (*Handler) SendToGroup(uid string, message string) {

}

func (*Handler) BindUid(end_id string, uid string) {

}

func (*Handler) UnbindUid(end_id string, uid string) {

}

func (*Handler) JoinGroup(end_id string, group string) {

}

func (*Handler) LeaveGroup(end_id string, group string) {

}
