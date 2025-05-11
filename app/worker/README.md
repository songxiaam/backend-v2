# Worker 服务
Worker 服务是 MetaLand 系统的后台任务处理服务，负责处理异步任务、定时任务和同步区块链合约数据等工作。

## 核心功能
1. UCE (Unified Compute Engine) 引擎
   ADB数仓元数据计算引擎，用于处理各类简单规则计算任务
   支持动态规则配置和执行
   提供预测试功能
2. 任务调度系统
   - 同步区块链合约数据
   - 自动标签任务
   - 临时任务处理

## 项目结构
```text
   worker/
   ├── cmd/                    # 命令行工具
   │   └── ctl/               # 控制工具
   ├── etc/                   # 配置文件
   ├── internal/              # 内部实现
   │   ├── config/           # 配置定义
   │   ├── logic/            # 业务逻辑
   │   ├── server/           # gRPC 服务实现
   │   ├── svc/              # 服务上下文
   │   └── task/             # 任务实现
   ├── pb/                    # Proto 文件及生成的代码
   └── Dockerfile            # 容器构建文件
```

##   配置说明
主要配置项（`etc/metaLand.yaml`）：
```yaml
Name: metaLand.worker.rpc        # 服务名称
ListenOn: 0.0.0.0:9000     # gRPC 监听地址
Database:
    Startup:                  # 默认数据库配置
        DSN: ...
    Bounty:                  # 计算数据库配置
        DSN: ...
RedisServer:
    Default:                 # Redis 配置
        Host: ...
OSS:                      # 对象存储配置
    Endpoint: ...
    AccessKeyID: ...
    Bucket: ...
Worker:                   # Worker 特定配置
    Sync:
        DefaultShopID: 4
```

## API 接口
Worker 服务通过 `gRPC` 提供服务，主要接口包括：

1. UCE 相关接口
```protobuf
service CRMWorker {
    // 触发 UCE 任务
    rpc TriggerUCETask(TriggerUCETaskRequest) returns (TriggerUCETaskResponse);
    // 获取 UCE 计算模型定义
    rpc GetUCESchemas(GetUCESchemasRequest) returns (GetUCESchemasResponse);
    // 规则预测试
    rpc PreTestRule(PreTestRuleRequest) returns (PreTestRuleResponse);
}
```

2. 数据同步接口
```protobuf
service ContractWorker {
  // 同步合约数据
  rpc SyncWeworkContract(SyncWeworkContractRequest) returns (SyncWeworkContractResponse);
}
```

## 开发指南

### 1. 添加新任务
1. 在 internal/task 目录下创建新的任务包
2. 实现任务逻辑
3. 在 internal/task/core.go 中注册任务
   
示例：
```go
type NewTask struct {
    svcCtx *svc.ServiceContext
}

func NewNewTask(svcCtx *svc.ServiceContext) *NewTask {
    return &NewTask{svcCtx: svcCtx}
}

func (t *NewTask) Start() error {
// 实现任务逻辑
}
```

### 2. 添加新的 gRPC 接口
1. 在 `pb/meta_land_worker.proto` 中定义新接口
2. 执行 `make proto` 生成代码
3. 在 `internal/logic` 中实现接口逻辑

### 3. 命令行工具开发
   Worker 服务提供命令行工具用于管理和调试：

```shell
# 运行命令行工具
./ctl [command]

# 示例：触发数据同步
./ctl sync --type=wework
```

## 部署
1. 构建镜像 
```shell
docker build -t metaLand-worker . -f Dockerfile
```

2. 运行容器
```shell
docker run -d \
--name metaLand-worker \
-p 9000:9000 \
-v /path/to/config:/app/etc \
metaLand-worker
```

## 常见问题
1. 任务执行失败
    - 检查数据库连接
    - 检查 `Redis` 连接
    - 查看错误日志

2. 性能问题
   - 调整并发配置
   - 优化数据库查询
   - 检查内存使用

3. 配置问题
   - 确保配置文件格式正确
   - 检查必要的环境变量
   - 验证第三方服务配置

