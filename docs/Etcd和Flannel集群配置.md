# Etcd和Flannel集群配置

> Version: 0.1

## IP
### clouds1
- 192.168.1.6
- 10.119.13.187

### clouds2
- 192.168.1.11
- 10.119.13.190


## 前置要求
### 开启内核ipv4转发功能
``` sh
echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
​
sysctl -p

重启服务器后失效，需要重新执行开启转发
```
### 清除iptables底层默认规则，并开启允许转发功能
``` sh
iptables -P INPUT ACCEPT
​
iptables -P FORWARD ACCEPT
​
iptables -F
​
iptables -L -n
```
### 关闭防火墙，如果开启防火墙，则最好打开2379和4001端口
``` sh
[root@node-1 ~]  systemctl disable firewalld.service
[root@node-1 ~]  systemctl stop firewalld.service
```


## ETCD
### install
``` sh
sudo apt update
sudo apt-get install etcd -y
```

### clouds1
<!-- $ etcd --name cloudos1 --initial-advertise-peer-urls http://192.168.1.6:2380 \  
  --listen-peer-urls http://192.168.1.6:2380 \  
  --listen-client-urls http://192.168.1.6:2379,http://127.0.0.1:2379 \  
  --advertise-client-urls http://192.168.1.6:2379 \  
  --initial-cluster-token etcd-cluster-1 \  
  --initial-cluster cloudos1=http://192.168.1.6:2380,cloudos2=http://192.168.1.11:2380 \  
  --initial-cluster-state new -->  

vim /etc/systemd/system/etcd.service   
修改一行： `ExecStart=/usr/bin/etcd $DAEMON_ARGS --config-file=/etc/etcd/etcd.conf`   

chmod -R 777 /var/lib/etcd/  

/etc/etcd/etcd.conf:  
name:  cloudos1   
data-dir:  /var/lib/etcd/cloudos1.etcd  
initial-advertise-peer-urls:  http://192.168.1.6:2380  
listen-peer-urls:   http://192.168.1.6:2380  
listen-client-urls:   http://192.168.1.6:2379,http://127.0.0.1:2379,http://192.168.1.6:4001,http://127.0.0.1:4001  
advertise-client-urls:  http://192.168.1.6:2379  
initial-cluster-token:  etcd-cluster-1  
initial-cluster:  cloudos1=http://192.168.1.6:2380,cloudos2=http://192.168.1.11:2380  
initial-cluster-state:   new  

### clouds2
<!-- $ etcd --name cloudos2 --initial-advertise-peer-urls http://192.168.1.11:2380 \  
  --listen-peer-urls http://192.168.1.11:2380 \  
  --listen-client-urls http://192.168.1.11:2379,http://127.0.0.1:2379 \  
  --advertise-client-urls http://192.168.1.11:2379 \  
  --initial-cluster-token etcd-cluster-1 \  
  --initial-cluster cloudos1=http://192.168.1.6:2380,cloudos2=http://192.168.1.11:2380 \  
  --initial-cluster-state new -->

vim /etc/systemd/system/etcd.service   
修改一行： `ExecStart=/usr/bin/etcd $DAEMON_ARGS --config-file=/etc/etcd/etcd.conf`   

chmod -R 777 /var/lib/etcd/  

/etc/etcd/etcd.conf:  
name:  cloudos2   
data-dir:  /var/lib/etcd/cloudos2.etcd    
initial-advertise-peer-urls:  http://192.168.1.11:2380   
listen-peer-urls:  http://192.168.1.11:2380   
listen-client-urls:  http://192.168.1.11:2379,http://127.0.0.1:2379,http://192.168.1.11:4001,http://127.0.0.1:4001   
advertise-client-urls:  http://192.168.1.11:2379   
initial-cluster-token:  etcd-cluster-1  
initial-cluster:  cloudos1=http://192.168.1.6:2380,cloudos2=http://192.168.1.11:2380   
initial-cluster-state:  new  

### check  
``` sh
# 重新加载system service配置
systemctl daemon-reload
# 启动
systemctl start etcd
# 设置为开机启动
systemctl enable etcd
# 查看启动状态
systemctl status etcd
# 查看进程服务
netstat -tnlp | grep -E  "4001|2380"
ps -ef | grep etcd
# 检查集群连通性
etcdctl -C http://192.168.1.6:2379 cluster-health
etcdctl -C http://192.168.1.6:4001 cluster-health
etcdctl -C http://192.168.1.11:2379 cluster-health
etcdctl -C http://192.168.1.11:4001 cluster-health
# 查看下etcd集群的状态
etcdctl member list 
```


## FLANNEL 需要先按照上述要求启动etcd
### install
``` sh
sudo apt update
sudo apt install flannel
```

### 配置node
``` sh
etcdctl put /coreos.com/network/config '{ "Network": "10.5.0.0/16", "Backend": {"Type": "vxlan"}}'

flannel

cat /run/flannel/subnet.env
FLANNEL_NETWORK=10.5.0.0/16
FLANNEL_SUBNET=10.5.72.1/24
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false

vim /etc/docker/daemon.json

{
  "bip": "10.5.75.1/24",
  "mtu": 1500
}

sudo systemctl restart docker
```

### check

ping即可


