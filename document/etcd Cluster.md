# etcd Cluster
_以下命令行代码，全部基于V3版本_

## 搭建本地集群
使用 etcd 作为分布式键值存储或者测试，最简单的方式是搭建本地集群，以下是介绍如何快速搭建本地集群

### goreman 提供了快捷管理本地集群的命令
安装 goreman 程序来控制基于 Profile 的应用程序.
$ `go get github.com/mattn/goreman`

### 启动本地多成员集群.
$ `goreman -f Procfile start`

 _注： 这里所说的 Procfile 文件是来自 etcd 的 gitub 项目的根目录下的 Procfile 文件，但是需要修改一下，将里面的 bin/etcd 修改为 etcd_

### 查看已启动的集群状态
$ `etcdctl member list -w table`

 _注： 由于是本地测试，我在一台物理机上启动了3个etcd节点，分别用了不同的端口_

```响应结果
+------------------+---------+--------+------------------------+------------------------+------------+
|        ID        | STATUS  |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+---------+--------+------------------------+------------------------+------------+
| 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 |  http://127.0.0.1:2379 |      false |
| 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:22380 | http://127.0.0.1:22379 |      false |
| fd422379fda50e48 | started | infra3 | http://127.0.0.1:32380 | http://127.0.0.1:32379 |      false |
+------------------+---------+--------+------------------------+------------------------+------------+
```

### 相关的goreman命令
* 查看procfile中的etcd条目
  $ `goreman check`

  以下是命令响应结果：
  valid procfile detected (etcd1, etcd2, etcd3)

* 停止一个节点
  $ `goreman run stop etcd1`

* 启动一个节点
  $ `goreman run start etcd1`


## 搭建正式环境集群
为了方便演示，以下是在同一台物理机上启动集群，所以用了不同的端口。  
正式环境多物理机环境下，可以使用相同的端口不同的IP来搭建集群。
推荐多物理机环境搭建集群，但不推荐跨局域网搭建集群。

### etcd的启动参数解释：
* `--name`
  etcd集群中的节点名，这里可以随意，可区分且不重复就行。

* `--listen-peer-urls`
  监听的用于节点之间通信的url，可监听多个，集群内部将通过这些url进行数据交互(如选举，数据同步等)

* `--initial-advertise-peer-urls`
  建议用于节点之间通信的url，节点间将以该值进行通信。

* `--listen-client-urls`
  监听的用于客户端通信的url,同样可以监听多个。

* `--advertise-client-urls`
  建议使用的客户端通信url,该值用于etcd代理或etcd成员与etcd节点通信。

* `--initial-cluster-token etcd-cluster-1`
  节点的token值，设置该值后集群将生成唯一id,并为每个节点也生成唯一id,当使用相同配置文件再启动一个集群时，只要该token值不一样，etcd集群就不会相互影响。

* `--initial-cluster`
  也就是集群中所有的initial-advertise-peer-urls 的合集。

* `--initial-cluster-state new`
  新建集群的标志(new/existing))

### 创建一个3节点集群，依次启动3个节点
- 启动第一个节点
  $ `etcd --name infra1 --listen-client-urls http://127.0.0.1:2379 --advertise-client-urls http://127.0.0.1:2379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster infra3=http://127.0.0.1:32380,infra2=http://127.0.0.1:22380,infra1=http://127.0.0.1:12380 --initial-cluster-state new`

  _参数--initial-cluster-state new表示是新的集群,参数--initial-cluster表明了集群的容量_


- 启动第二个节点
  $ `etcd --name infra2 --listen-client-urls http://127.0.0.1:22379 --advertise-client-urls http://127.0.0.1:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster infra3=http://127.0.0.1:32380,infra2=http://127.0.0.1:22380,infra1=http://127.0.0.1:12380 --initial-cluster-state new`

- 启动第三个节点
  $ `etcd --name infra3 --listen-client-urls http://127.0.0.1:32379 --advertise-client-urls http://127.0.0.1:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster infra3=http://127.0.0.1:32380,infra2=http://127.0.0.1:22380,infra1=http://127.0.0.1:12380 --initial-cluster-state new`

- 查看集群状态
  $ `etcdctl member list -w table`

  ```命令响应结果可以看到成功启动了3节点的集群
  +------------------+---------+--------+------------------------+------------------------+------------+
  |        ID        | STATUS  |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
  +------------------+---------+--------+------------------------+------------------------+------------+
  | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 |  http://127.0.0.1:2379 |      false |
  | 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:22380 | http://127.0.0.1:22379 |      false |
  | fd422379fda50e48 | started | infra3 | http://127.0.0.1:32380 | http://127.0.0.1:32379 |      false |
  +------------------+---------+--------+------------------------+------------------------+------------+
  ```

### 动态扩容集群 
- 扩容集群
  $ `etcdctl --endpoints=http://127.0.0.1:2379 member add infra4 --peer-urls=http://127.0.0.1:42380`

  ```执行命令后的正确响应内容如下：
  Member 72a6d8ad06c5d803 added to cluster ef37ad9dc622a7c4

  ETCD_NAME="infra4"
  ETCD_INITIAL_CLUSTER="infra4=http://127.0.0.1:42380,infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380"
  ETCD_INITIAL_ADVERTISE_PEER_URLS="http://127.0.0.1:42380"
  ETCD_INITIAL_CLUSTER_STATE="existing"
  ```

- 查看扩容后的集群状态
  $ `etcdctl member list -w table`

  ```可以看到响应结果，第4个节点已经扩容成功，等待启动
  +------------------+-----------+--------+------------------------+------------------------+------------+
  |        ID        |  STATUS   |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
  +------------------+-----------+--------+------------------------+------------------------+------------+
  | 72a6d8ad06c5d803 | unstarted |        | http://127.0.0.1:42380 |                        |      false |
  | 8211f1d0f64f3269 |   started | infra1 | http://127.0.0.1:12380 |  http://127.0.0.1:2379 |      false |
  | 91bc3c398fb3c146 |   started | infra2 | http://127.0.0.1:22380 | http://127.0.0.1:22379 |      false |
  | fd422379fda50e48 |   started | infra3 | http://127.0.0.1:32380 | http://127.0.0.1:32379 |      false |
  +------------------+-----------+--------+------------------------+------------------------+------------+
  ```

- 启动加入的第4个节点
  $ `etcd --name infra4 --listen-client-urls http://127.0.0.1:42379 --advertise-client-urls http://127.0.0.1:42379 --listen-peer-urls http://127.0.0.1:42380 --initial-advertise-peer-urls http://127.0.0.1:42380 --initial-cluster-token etcd-cluster-1 --initial-cluster infra4=http://127.0.0.1:42380,infra3=http://127.0.0.1:32380,infra2=http://127.0.0.1:22380,infra1=http://127.0.0.1:12380 --initial-cluster-state existing`

  _参数--initial-cluster-state existing表示已有的集群中增加节点_

- 查看启动新节点后的集群状态
  $ `etcdctl member list -w table`
  
  ```响应结果可以看到新增节点已成功启动并加入了集群
  +------------------+---------+--------+------------------------+------------------------+------------+
  |        ID        | STATUS  |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
  +------------------+---------+--------+------------------------+------------------------+------------+
  | 72a6d8ad06c5d803 | started | infra4 | http://127.0.0.1:42380 | http://127.0.0.1:42379 |      false |
  | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 |  http://127.0.0.1:2379 |      false |
  | 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:22380 | http://127.0.0.1:22379 |      false |
  | fd422379fda50e48 | started | infra3 | http://127.0.0.1:32380 | http://127.0.0.1:32379 |      false |
  +------------------+---------+--------+------------------------+------------------------+------------+
  ```

- 移除集群中的节点
  $ `etcdctl --endpoints=http://127.0.0.1:2379 member remove 72a6d8ad06c5d803`

  ```成功的响应为：
  Member 72a6d8ad06c5d803 removed from cluster ef37ad9dc622a7c4
  ```

  _此处以移除第4个节点为例，参数对应第4个节点的ID `72a6d8ad06c5d803`_

  $ `etcdctl member list -w table`

  ```响应结果可以看到成功移除了第4个节点
  +------------------+---------+--------+------------------------+------------------------+------------+
  |        ID        | STATUS  |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
  +------------------+---------+--------+------------------------+------------------------+------------+
  | 8211f1d0f64f3269 | started | infra1 | http://127.0.0.1:12380 |  http://127.0.0.1:2379 |      false |
  | 91bc3c398fb3c146 | started | infra2 | http://127.0.0.1:22380 | http://127.0.0.1:22379 |      false |
  | fd422379fda50e48 | started | infra3 | http://127.0.0.1:32380 | http://127.0.0.1:32379 |      false |
  +------------------+---------+--------+------------------------+------------------------+------------+
  ```


## ETCD gRPC代理
- etcd proxy
 etcd提供了proxy功能，即代理功能，etcd可以代理的方式来运行。
 etcd代理可以运行在每一台主机，在这种代理模式下，etcd的作用就是一个反向代理，把客户端的etcd请求转发到真正的etcd集群。这种方式既加强了集群的弹性，又不会降低集群的写的性能。

- 启动一个代理节点监听23790端口
  $ `etcd grpc-proxy start --endpoints=127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379 --listen-addr=127.0.0.1:23790 --advertise-client-url=127.0.0.1:23790`

- put到代理节点测试
  $ `etcdctl --endpoints=http://127.0.0.1:23790 put /foo bar`
  
  成功响应 ok

- 向集群中的其他节点查询key测试代理结果
  $ `etcdctl --endpoints=http://127.0.0.1:22379 --prefix --keys-only=false get /`

  ```成功取到结果
  /foo
  bar
  ```

## 本地测试可视化工具
[ETCD Manager](https://etcdmanager.io/)
