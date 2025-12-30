```bash
# consul 相关指令
## 列出注册过的服务实例
 curl -s http://127.0.0.1:8500/v1/agent/services
## 删除服务实例
 curl -X PUT http://127.0.0.1:8500/v1/agent/service/deregister/stock_6860223110692197945
```