package useradmin

import "errors"

var (
	// controller が参照しているやつ
	ErrNotFound        = errors.New("not found")
	ErrSelfRoleChange  = errors.New("self role change not allowed")
	ErrSelfDeactivate  = errors.New("self deactivate forbidden")
	ErrAlreadyInactive = errors.New("already inactive")
)
