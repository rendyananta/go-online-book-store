package user

const (
	queryGetUserByEmail = `select id, email, password from users where email = ?`
	queryGetUserByID    = `select id, email, password from users where id = ?`
	queryInsertUser     = `insert into users (id, name, email, password, created_at, updated_at) values (?, ?, ?, ?, ?, ?) returning id`
)
