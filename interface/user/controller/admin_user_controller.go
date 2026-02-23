package usercontroller

import (
	"net/http"
	"strconv"
	"strings"

	user "roulettept/domain/user/model"
	"roulettept/usecase/useradmin"

	"github.com/labstack/echo/v4"
)

type AdminUserController struct {
	uc *useradmin.Service
}

func NewAdminUserController(uc *useradmin.Service) *AdminUserController {
	return &AdminUserController{uc: uc}
}

func (h *AdminUserController) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	var role *user.UserRole
	if r := strings.TrimSpace(c.QueryParam("role")); r != "" {
		rr := user.UserRole(r)
		role = &rr
	}

	var isActive *bool
	if v := strings.TrimSpace(c.QueryParam("is_active")); v != "" {
		b := v == "true"
		isActive = &b
	}

	q := c.QueryParam("q")

	out, err := h.uc.ListUsers(c.Request().Context(), useradmin.ListInput{
		Page:     page,
		Limit:    limit,
		Role:     role,
		IsActive: isActive,
		Q:        q,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}
	return c.JSON(http.StatusOK, okResp(out))
}

type updateRoleReq struct {
	Role string `json:"role"`
}

func (h *AdminUserController) UpdateRole(c echo.Context) error {
	targetID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	actorID, _ := c.Get("user_id").(int64)

	var req updateRoleReq
	if err := c.Bind(&req); err != nil || strings.TrimSpace(req.Role) == "" {
		return c.JSON(http.StatusBadRequest, errResp("VALIDATION_ERROR", "invalid body"))
	}

	role := user.UserRole(strings.TrimSpace(req.Role))

	err := h.uc.UpdateRole(c.Request().Context(), actorID, targetID, role)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, okResp(map[string]any{}))
	case useradmin.ErrSelfRoleChange:
		return c.JSON(http.StatusConflict, errResp("CONFLICT", "self_change_not_allowed"))
	case useradmin.ErrNotFound:
		return c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found"))
	default:
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}
}

func (h *AdminUserController) Deactivate(c echo.Context) error {
	targetID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	actorID, _ := c.Get("user_id").(int64)

	err := h.uc.Deactivate(c.Request().Context(), actorID, targetID)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, okResp(map[string]any{}))
	case useradmin.ErrSelfDeactivate:
		return c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "self_deactivate_forbidden"))
	case useradmin.ErrAlreadyInactive:
		return c.JSON(http.StatusConflict, errResp("CONFLICT", "already_inactive"))
	case useradmin.ErrNotFound:
		return c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found"))
	default:
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}
}

func (h *AdminUserController) ForceLogout(c echo.Context) error {
	targetID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	actorID, _ := c.Get("user_id").(int64)

	err := h.uc.ForceLogout(c.Request().Context(), actorID, targetID)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, okResp(map[string]any{}))
	case useradmin.ErrNotFound:
		return c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found"))
	default:
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}
}
