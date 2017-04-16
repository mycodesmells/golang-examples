package sessions

type Session struct {
	VisitCount int `json:"visitCount"`
}

type Store interface {
	Get(string) (Session, error)
	Set(string, Session) error
}
