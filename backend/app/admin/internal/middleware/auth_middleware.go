package middleware

import (
	"net/http"

	"admin_backend/app/admin/internal/config"
	"admin_backend/pkg/common/casbin"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/utils/jwt"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	Config     config.Config
	JWTManager *jwt.JWTManager
}

func NewAuthMiddleware(c config.Config, jwtManager *jwt.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		Config:     c,
		JWTManager: jwtManager,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.Error(w, xerr.NewErrCodeMsg(xerr.TokenInvalid, "请先登录"))
			return
		}

		// 使用JWTManager进行验证，包含黑名单检查
		claims, err := m.JWTManager.ValidateAccessTokenWithBlacklist(r.Context(), tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.Error(w, err)
			return
		}

		// // 3. 租户隔离验证
		// if err := m.validateTenant(r, claims); err != nil {
		// 	httpx.Error(w, xerr.NewErrCodeMsg(http.StatusForbidden, err.Error()))
		// 	return
		// }

		// // 4. 权限验证（可选）
		// if len(m.requiredRoles) > 0 {
		// 	if err := m.validateRoles(claims); err != nil {
		// 		httpx.Error(w, xerr.NewErrCodeMsg(http.StatusForbidden, "权限不足"))
		// 		return
		// 	}
		// }

		if err := m.CheckPermission(r, claims); err != nil {
			w.WriteHeader(http.StatusForbidden)
			httpx.Error(w, xerr.NewErrCodeMsg(http.StatusForbidden, "权限不足"))
			return
		}

		// 将租户ID和用户ID存入context
		ctx := r.Context()
		ctx = contextutil.SetUserIDToCtx(ctx, claims.UserID)
		ctx = contextutil.SetTenantCodeToCtx(ctx, claims.TenantCode)
		ctx = contextutil.SetRoleCodeToCtx(ctx, claims.RoleCode)

		next(w, r.WithContext(ctx))
	}
}

// // 租户隔离验证
// func (m *AuthMiddleware) validateTenant(r *http.Request, claims *jwt.Claims) error {
// 	reqTenant := r.Header.Get("x")
// 	if reqTenant != "" && claims.TenantID != reqTenant {
// 		return errors.New("租户身份不匹配")
// 	}
// 	return nil
// }

// 角色权限验证（可选）
func (m *AuthMiddleware) validateRoles(claims *jwt.Claims) error {
	// // 示例：从数据库或缓存中获取用户角色
	// userRoles := getUserRolesFromDB(claims.UserID)

	// // 检查用户是否具备所需角色
	// for _, requiredRole := range m.requiredRoles {
	// 	if !contains(userRoles, requiredRole) {
	// 		return error.New("用户角色不符合要求")
	// 	}
	// }

	return nil
}

// 检查权限
func (m *AuthMiddleware) CheckPermission(r *http.Request, claims *jwt.Claims) error {

	userID := claims.UserID
	tenantCode := claims.TenantCode
	resource := r.Header.Get("X-Resource")
	action := r.Header.Get("X-Action")
	permType := r.Header.Get("X-Type")

	if resource != "" { // 资源路径不为空时
		manager := casbin.GetCasbinManager()
		if ok, err := manager.CheckPermission(userID, tenantCode, resource, action, permType); err != nil {
			return err
		} else if !ok {
			return xerr.NewErrCodeMsg(xerr.ForbiddenError, "用户没有权限")
		}
	}

	return nil
}
