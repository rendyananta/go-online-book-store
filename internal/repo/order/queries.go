package order

const (
	queryGetUserOrdersPagination = `select id, grand_total, status, created_at, updated_at from orders 
					  where user_id = ? and deleted_at is null and id > ? order by id limit ?`

	queryGetOrderDetail = `select id, user_id, grand_total, status, created_at, updated_at from orders 
					  where id = ? and deleted_at is null`

	queryGetOrderLines = `select id, line_reference_type, line_reference_id, amount, quantity, subtotal from order_lines where order_id = ?`

	queryInsertOrder = `insert into orders (id, user_id, grand_total, status, created_at, updated_at) 
					values (?, ?, ?, ?, ?, ?)`

	queryInsertOrderLine = `insert into order_lines (id, order_id, line_reference_type, line_reference_id, amount, quantity, subtotal) 
					values (?, ?, ?, ?, ?, ?, ?)`
)
