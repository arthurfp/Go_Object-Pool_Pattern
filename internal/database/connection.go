package database

type Connection struct {
	ID int
}

func (c *Connection) Execute(query string) error {
	// Simulate query execution
	return nil
}
