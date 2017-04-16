package sessions

import "errors"

type memoryStore struct {
	sessions map[string]Session
}

func NewMemoryStore() Store {
	return &memoryStore{
		sessions: make(map[string]Session),
	}
}

func (m memoryStore) Get(id string) (Session, error) {
	session, ok := m.sessions[id]
	if !ok {
		return Session{}, errors.New("session not found")
	}

	return session, nil
}

func (m *memoryStore) Set(id string, session Session) error {
	m.sessions[id] = session
	return nil
}
