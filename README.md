# one

名为一，意为起点，这个项目是一个Gin项目，通过暴露的网络接口提供以下服务：

- 报表一览服务
- 孤块监测

## 如何使用

服务运行：

- go build
- 填写 db/miner.json 和 db/dingdingbot.json
- ./one --initdb

用户使用：

```
# 获取报表，在浏览器中访问
127.0.0.1:8080/miner/report/download

# 查看孤块状态，在浏览器中访问
127.0.0.1:8080/orphanblock/view/last5block/human
127.0.0.1:8080/orphanblock/view/all/human
```

# 后续更新

- 增量增加miner信息
- 配置钉钉机器人通知服务

