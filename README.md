antgo
=================

antgo使用经典的Gateway和Worker进程模型。Gateway进程负责维持客户端连接，并转发客户端的数据给Worker进程处理；Worker进程负责处理实际的业务逻辑，并将结果推送给对应的客户端。Gateway服务和Worker服务可以分开部署在不同的服务器上，实现分布式集群。还有Register负责内部组件管理和消息传递。

antgo提供非常方便的API，可以全局广播数据、可以向某个群体广播数据、也可以向某个特定客户端推送数据。配合Workerman的定时器，也可以定时推送数据。

启动
=======

```go run test/server.go```

```shell
telnet 127.0.0.1 23
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
```