package middleware

// import (
// 	"net/http"

// 	"golang.org/x/exp/slices"
// )

// // CorsMiddleware 是一个自定义的跨域中间件
// func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {

// 	// 定义允许的来源列表
// 	allowedOrigins := []string{
// 		"http://192.168.0.2:5173",  // 前端本地IP
// 		"http://192.168.0.12:8888", // 后端本机IP
// 		"http://localhost:5173",
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		// // 允许所有来源访问（可以根据需求限制特定域名）
// 		// w.Header().Set("Access-Control-Allow-Origin", "*")

// 		// 获取请求的 Origin 头部
// 		origin := r.Header.Get("Origin")
// 		if origin != "" && slices.Contains(allowedOrigins, origin) {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 		} else {
// 			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
// 			return
// 		}

// 		// 允许的 HTTP 方法
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		// 允许的自定义请求头
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		// 如果是预检请求（OPTIONS 方法），直接返回 200
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		// 调用下一个处理器
// 		next(w, r)
// 	}
// }
