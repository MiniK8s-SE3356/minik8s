# NFS配置

> author: Haocheng Wang  
> version: 1.0


## 作用
NFS在不改变用户视角的文件系统使用方式的情况下，打通了集群内各个节点的存储，使得多机PV/PVC成为可能

## 节点信息
- cloudos1
  - IP: 192.168.1.6 
  - role: NFS server 
- cloudos2
  - IP: 192.168.1.11
  - role: NFS client
  
## 具体配置步骤
设计中，NFS server开放的目录为`/var/lib/minik8s/volumes`，NFS client将这个路径挂载到自己的`/var/lib/minik8s/volumes`下

> 注意，配置前请配置网络流默认策略为ACCEPT（可见flannel配置部分）

### NFS server配置
``` shell
sudo apt update
sudo apt install nfs-kernel-server
sudo mkdir /var/lib/minik8s/volumes -p
```
随后vim进入nfs配置文件
``` shell
sudo vim /etc/exports
```
添加条目如下,然后保存退出
``` 
/var/lib/minik8s/volumes  192.168.1.11(rw,sync,no_root_squash,no_subtree_check)
```
重启nfs系统进程
``` shell
sudo systemctl restart nfs-kernel-server
```


### NFS client配置
``` shell
sudo apt update
sudo apt install nfs-common
sudo mkdir /var/lib/minik8s/volumes -p
sudo mount 192.168.1.6:/var/lib/minik8s/volumes /var/lib/minik8s/volumes
```

可以设置每次启动时自动挂载，在`/etc/fstab` 里添加如下一条
``` 
192.168.1.6:/var/lib/minik8s/volumes   /var/lib/minik8s/volumes   nfs auto,nofail,noatime,nolock,intr,tcp,actimeo=1800 0 0
```

## 测试
在挂载的目录下进行文件操作，可以看到两机一致即可