下载
https://github.com/etcd-io/etcd/releases

$ls
Documentation			README.md			etcd				etcdctl
README-etcdctl.md		READMEv2-etcdctl.md		etcd-v3.4.1-darwin-amd64.zip

$./etcd
[WARNING] Deprecated '--logger=capnslog' flag is set; use '--logger=zap' flag instead
2019-09-28 00:51:04.911642 I | etcdmain: etcd Version: 3.4.1
2019-09-28 00:51:04.911784 I | etcdmain: Git SHA: a14579fbf

2019-09-28 00:51:05.958205 I | etcdserver: published {Name:default ClientURLs:[http://localhost:2379]} to cluster cdf818194e3a8c32
2019-09-28 00:51:05.958960 N | embed: serving insecure client requests on 127.0.0.1:2379, this is strongly discouraged!


$./etcdctl version
etcdctl version: 3.4.1
API version: 3.4

$./etcdctl put key1 v1
OK

<nil>[key:"/job/v3" create_revision:3 mod_revision:3 version:1 value:"push the box" ]

$./etcdctl del key1
1

$./etcdctl del /job
0

$./etcdctl del /job/v3
1

type:DELETE
 kv:key:"/job/v3" mod_revision:5   prevKey:key:"/job/v3" create_revision:3 mod_revision:3 version:1 value:"push the box"

 



$./etcdctl get key1
key1
v1




