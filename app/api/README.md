## 开发指南

### 添加新 API
1. 在对应的分组目录下创建或修改 .api 文件
2. 定义 API 路由、请求/响应结构
3. 运行 `goctl` 生成代码
4. 实现对应的 `handler` 逻辑

示例 API 定义:
```api
@server(
    group: metaLand/startup
    prefix: api/startup
    middleware: OIDCAuthMiddleware
)

service metaLand {
    @doc "查询项目列表"
    @handler ListStartups
    get /startups  (ListStartupsRequest) returns (ListResponse)
}


type (
    ListStartupsRequest {
        Page int `json:"page"`
        Size int `json:"size"`
    }

    ListResponse {

    }
)
```
生成代码
```shell
goctl api go -api  api/meta_land.api -dir . -style go_zero
```


## 配置说明
配置文件位于 `etc/crm.yaml`，主要包含：

- HTTP 服务配置
- 数据库配置
- 缓存配置
- 认证配置
- 第三方服务配置

## 启动服务
```shell
go run meta_land.go -f etc/meta_land.yaml
```
