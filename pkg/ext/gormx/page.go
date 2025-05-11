package gormx

type PagedResult[T any] struct {
	Total int64 `json:"total"`
	Size  int   `json:"size"`
	Page  int   `json:"page"`
	List  []T   `json:"list"`
}

type PageParams struct {
	Size int `json:"size"`
	Page int `json:"page"`
}
