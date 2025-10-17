package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Infof("Method: %s, URL: %s", r.Method, r.URL.String())
		for name, headers := range r.Header {
			for _, h := range headers {
				logx.Infof("Header: %s - %s", name, h)
			}
		}
		// 获取 GET 请求的查询参数
		logx.Infof("Query Params: %+v", r.URL.Query())
		// 获取Post form参数
		if r.Method == http.MethodPost && strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				logx.Infof("Error parsing form data: %v", err)
			} else {
				logx.Infof("Form Params: %+v", r.Form)
			}
		}
		// 获取 POST 请求的 JSON 数据
		if r.Method == http.MethodPost && strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logx.Infof("Error reading request body: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				logx.Infof("Error decoding JSON body: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			logx.Infof("JSON Params: %+v", jsonBody)
			// 重新设置请求体，以便后续处理可以继续使用
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next(w, r)
	}
}
