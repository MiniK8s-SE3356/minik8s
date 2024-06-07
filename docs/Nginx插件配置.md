# Nginx插件配置

> author: Haocheng Wang  
> Version: 0.1


## 作用
Nginx插件只部署在控制平面，协同各节点上的/etc/hosts配置，作为minik8s的DNS/反向代理解决方案  

## 安装
> OS: Ubuntu 20.04  

``` shell
sudo apt update
sudo apt install nginx
```

## 检查

查看系统级进程状态：  
```shell
systemctl status nginx
```

测试nginx：
```shell
curl http://localhost
```