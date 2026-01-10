```bash
# consul 相关指令
## 列出注册过的服务实例
 curl -s http://127.0.0.1:8500/v1/agent/services
## 删除服务实例
 curl -X PUT http://127.0.0.1:8500/v1/agent/service/deregister/stock_6860223110692197945
```

```bash
# stripe客户端 监听 payment端口
 stripe listen --forward-to 10.11.71.154:8284/api/webhook
```

```bash
# 在系统中永久添加环境变量
## 打开当前用户的.bashrc文件
vim ~/.bashrc
添加变量： 在文件末尾添加以下行： export MY_VARIABLE="value" 如果是添加路径到 PATH： export PATH="$PATH:/your/custom/path"
## 使配置立即生效
source ~/.bashrc
## 查看环境变量
echo $MY_VARIABLE
```
```bash
8431e37c27988fb12e389d6af6cc6a10
```