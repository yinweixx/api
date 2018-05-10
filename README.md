# <center> 安云云计算 API 网关 </center> 
## 概览
   网关使用go语言编写，运行在docker容器，使用alpine linux的镜像编译；使用轻量级的协程支持高并发的API访问；能集成到平台的容器管理服务，提供API的动态集群弹性部署；支持平台统一认证鉴权平台，提供细粒度的API访问控制；集成SNMP，提供标准的API运维数据到运维平台；

## 特性
1. Docker容器部署
2. 从元数据服务器获取API定义并且自动部署API
3. 将当前节点注册到分布式DNS服务
4. 服务的DNS查询，服务元数据缓存
5. 服务的调用支持多种方式(JSON-RPC、NATS、WS、RESTFUL）
6. 管理总线的接入，发送和接收处理服务总线的消息
7. 监控（CPU、内存、网络、连接数、会话数）
8. 从全局配置文件获取启动配置参数（端口、最大会话数等配置）
## 平台API
   API使用[RAML 1.0](https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md)定义，打包后（zip）通过API管理门户上传，API管理门户解压并且解析所有包内的RAML文件，生成kv键值对存入平台全局配置中间件etcd；当API网关启动后，会通过etcd的[v3客户端](https://github.com/coreos/etcd/tree/master/clientv3)获取所有的API定义并且部署；在API有添加或者API有变更时，API管理平台会发送变更事件到所有的API网关节点，网关节点会重新部署所有的API。
## API网关的调度
   API网关本身无调度的功能；网关在启动后会运行多个crontab监控协程，将相关的CPU、内存以及会话连接数据发送到监控平台网关；当API网关的CPU、内存以及会话连接数超过预定的阈值，也会主动生成报警事件通过事件消息系统发送给容器调度平台；监控平台也会根据预先配置好的策略生成事件通过事件消息系统发送给容器调度平台，容器调度平台最终综合所有的事件后对API网关进行调度，创建新的网关加入到API服务集群或者删除现有的空闲的网关。
## 服务的查询
   API网关使用DNS工具包通过查询服务的DNS [SRV记录](https://en.wikipedia.org/wiki/SRV_record)来完成服务的查找；如果API网关查找到了服务的SRV记录，将会把记录写到Redis缓存，并根据配置设置缓存的[TTL](https://en.wikipedia.org/wiki/Time_to_live)；如果缓存的服务SRV记录过期API网关会重新通过DNS查找服务的SRV记录；如果查找不到服务的DNS记录，API网关会发送事件通知容器调度平台创建该服务，容器服务在容器创建完后会通知API网关，如果创建失败或者服务不存在，API网关会将该服务从缓存中剔除，并且发送创建失败的事件。