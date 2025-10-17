package main

import (
	"flag"
	"fmt"

	"admin_backend/app/admin/internal/config"
	"admin_backend/app/admin/internal/handler"
	"admin_backend/app/admin/internal/middleware"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/pkg/common/response"
	"admin_backend/pkg/common/validate"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/admin.yaml", "the config file")

func main() {
	// 加载配置文件
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 配置log
	c.LoadLogConf()

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*")) // TODO 部署的时候注意跨域问题解决
	defer server.Stop()

	// 注册全局中间件
	server.Use(middleware.LoggerMiddleware)

	// 注册路由
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 添加验证器
	httpx.SetValidator(validate.DefaultValidator)

	// 使用拦截器，处理返回值
	httpx.SetOkHandler(response.OkHanandler)
	httpx.SetErrorHandlerCtx(response.ErrHandler(c.Name))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
