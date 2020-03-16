package cleanup

import (
	"sync"
	"time"
)

type Memory struct {
	*sync.Mutex
}

func (m *Memory) Use() {
	m.Lock()
	time.Sleep(50 * time.Millisecond)
	m.Unlock()
}

type Connection struct {}

type ConnPool struct {
	Conns chan(*Connection)
}

func (cp ConnPool) GetConn() *Connection {
	return <-cp.Conns
}
func (cp ConnPool) FreeConn(conn *Connection) {
	cp.Conns <- conn
}

type Database struct {
	values     map[int]string
	connection *Connection
}

func (d *Database) Square(n int) int {
	if d.connection == nil {
		return -2
	}

	time.Sleep(100 * time.Millisecond)

	if n > 100 {
		return -1
	}
	return n * n
}

func (d *Database) Open(cp ConnPool) {
	time.Sleep(2 * time.Second)
	d.connection = cp.GetConn()

	d.values = make(map[int]string)
}

func (d *Database) Close(cp ConnPool) error {
	time.Sleep(100 * time.Millisecond)
	cp.FreeConn(d.connection)
	d.connection = nil

	return nil
}
