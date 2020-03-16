package cleanup

import (
	"fmt"
	"testing"
)

func TestDatabase_Sequential(t *testing.T) {
	pool:=initializeConnectionPool(2)

	db := &Database{}
	db.Open(pool)
	defer db.Close(pool)

	for i := 0; i < 100; i++ {
		func(i int) {
			t.Run(fmt.Sprintf("test for %d", i), func(t *testing.T) {
				res := db.Square(i)
				if res != i*i {
					t.Errorf("Expected %v, got %v", i*i, res)
				}
			})
		}(i)
	}
}

func TestDatabase_Parallel(t *testing.T) {
	pool := initializeConnectionPool(10)

	db := &Database{}
	db.Open(pool)
	defer db.Close(pool)

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("test for %d", i), func(t *testing.T) {
			t.Parallel()

			res := db.Square(i)
			if res != i*i {
				t.Errorf("Expected %v, got %v", i*i, res)
			}
		})
	}
}

func TestDatabase_DuplicateSetup(t *testing.T) {
	pool := initializeConnectionPool(10)

	for i := 0; i < 100; i++ {
		func(i int) {
			t.Run(fmt.Sprintf("test for %d", i), func(t *testing.T) {
				t.Parallel()

				db := &Database{}
				db.Open(pool)
				defer db.Close(pool)

				res := db.Square(i)
				if res != i*i {
					t.Errorf("Expected %v, got %v", i*i, res)
				}
			})
		}(i)
	}
}

func TestDatabase_Cleanup(t *testing.T) {
	pool := initializeConnectionPool(10)

	db := &Database{}
	db.Open(pool)
	t.Cleanup(func() {
		db.Close(pool)
	})

	for i := 0; i < 100; i++ {
		func(i int) {
			t.Run(fmt.Sprintf("test for %d", i), func(t *testing.T) {
				t.Parallel()

				res := db.Square(i)
				if res != i*i {
					t.Errorf("Expected %v, got %v", i*i, res)
				}
			})
		}(i)
	}
}

func initializeConnectionPool(size int) ConnPool {
	connections := make(chan *Connection, size)

	for i := 0; i < size; i++ {
		connections <- &Connection{}
	}

	return ConnPool{Conns: connections}
}
