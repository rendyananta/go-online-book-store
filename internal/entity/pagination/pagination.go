package pagination

type Param struct {
	PerPage int
	LastID  string
}

type Base struct {
	Data    interface{} `json:"data"`
	PerPage int         `json:"per_page"`
	LastID  string      `json:"last_id"`
}

type PageInfo struct {
	PerPage int    `json:"per_page"`
	LastID  string `json:"last_id"`
}
