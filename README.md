# 15-federated-learning-coordinator

纯Go联邦学习隐私计算协调与模型聚合服务。默认使用内存仓储，仓储、对象存储、时钟、密码学与通知均通过接口注入，可替换为生产实现。

## 启动

```bash
go run ./cmd/server
curl http://127.0.0.1:8080/healthz
```

主流程：创建协作组、注册参与方、创建轮次、提交更新、聚合、发布模型。参见`examples/flow.sh`。

## 质量门禁

```bash
gofmt -w . && go vet ./... && go test ./... && go build ./...
./scripts/check-lines.sh
```

非测试Go源码目标不少于2600行，排除测试、生成代码、依赖和构建产物。
