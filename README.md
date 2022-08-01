# apinto-import
## 导入配置
- 当前只考虑未对接etcd的版本（v0.6.4）升级到对接etcd的版本（v0.7.0）
- 升级原因：
  - 由于早期选型时对etcd的源码调研不够完整，旧版本只使用了etcd的raft协议实现，eosc重新实现了kv的处理，该实现的可靠性、性能未经证明
  - 本次升级抛弃了eosc自己实现的kv，使用etcd内置的kv实现，该模型及代码经历过海量用户长时间、大规模的使用，可靠性、性能方面具有一定保证，并且方便以后同步升级到etcd的新版本
  
## 快速使用
1、下载并解压导入程序
```
wget https://github.com/eolinker/apinto-import/releases/download/v1.0.0/apinto-import-v1.0.0.linux.x64.tar.gz && tar -zxvf apinto-import-v1.0.0.linux.x64.tar.gz && cd apinto-import
```
2、导入配置数据
```
./apinto-import import -path "{压缩包名称}" -apinto-address {apinto访问地址}
```
示例：
```
./apinto-import import -path "export_2022-07-29 161215.zip" -apinto-address http://127.0.0.1:9400
```

## 升级流程
### 单节点升级
1、将旧数据导出，浏览器访问接口：{ip}:{port}/export

2、关闭旧节点
```
./apinto stop
```
3、下载并解压新版本节点
 ```
cd {存放目录} && wget https://github.com/eolinker/apinto/releases/download/v0.7.0/apinto-v0.7.0.linux.x64.tar.gz && tar -zxvf apinto-v0.7.0.linux.x64.tar.gz && cd apinto
```
4、启动新节点
```
./apinto start
```
5、下载并解压导入程序
```
wget https://github.com/eolinker/apinto-import/releases/download/v1.0.0/apinto-import-v1.0.0.linux.x64.tar.gz && tar -zxvf apinto-import-v1.0.0.linux.x64.tar.gz && cd apinto-import
```
6、导入配置数据（当前版本只支持zip类型文件）
```
./apinto-import import -path "{压缩包名称}" --apinto-address {apinto访问地址}
```
示例：
```
./apinto-import import "export_2022-07-29 161215.zip" --apinto-address http://127.0.0.1:9400
```
### 集群节点升级
1、下载并解压导入程序
```
wget https://github.com/eolinker/apinto-import/releases/download/v1.0.0/apinto-import-v1.0.0.linux.x64.tar.gz && tar -zxvf apinto-import-v1.0.0.linux.x64.tar.gz && cd apinto-import
```
2、将旧数据导出，浏览器访问接口：{ip}:{port}/export

3、进入任意节点服务器（下述描述为节点A），让该节点离开集群
```
./apinto leave
```
4、关闭节点A
```
./apinto stop
```
5、下载并解压新版本节点
```
cd {存放目录} && wget https://github.com/eolinker/apinto/releases/download/v0.7.0/apinto-v0.7.0.linux.x64.tar.gz && tar -zxvf apinto-v0.7.0.linux.x64.tar.gz && cd apinto
```
6、启动新节点A
```
./apinto start
```
7、导入配置数据
```
./apinto-import import -path "{压缩包名称}" -apinto-address {apinto访问地址}
```
示例：
```
./apinto-import import -path "export_2022-07-29 161215.zip" -apinto-address http://127.0.0.1:9400
```
8、进入到剩余的其他节点，依次执行步骤3、4、5、6
9、新节点加入节点A所在集群
```
./apinto join --ip {新节点广播ip} --addr={节点A请求地址}
```
示例：
```
./apinto join --ip 10.18.0.1 --addr=10.18.0.2:9400
```
