package userrepo

import user "roulettept/domain/user/model"

// /admin/users の検索条件
type UserListFilter struct {
	Role     *user.UserRole
	IsActive *bool
	Q        string
}
