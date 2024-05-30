#!/bin/bash

# 环境Ubuntu 20.04
# 请在sudo权限下运行该脚本

echo $SHELL

# 更新apt
apt update
# 安装依赖工具
apt install curl -y
apt install net-tools -y


# DNS插件：Nginx
apt install nginx -y


# # 异步消息队列插件：RabbitMQ
# apt-get install erlang -y
# apt-get install rabbitmq-server -y
# ## 启动控制面板
# rabbitmq-plugins enable rabbitmq_management
# ## 添加用户并设置权限
# rabbitmqctl add_user admin admin
# rabbitmqctl set_user_tags admin administrator
# rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"


# 分布式一致性持久化存储插件：Etcd
## 安装
ETCD_VER=v3.4.32
### choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test
curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
## 为etcd创建系统级进程
## 移入/usr/bin
cp /tmp/etcd-download-test/etcd /usr/bin/
cp /tmp/etcd-download-test/etcdctl /usr/bin/
## 创建/etc/etcd/etcd.conf
chmod -R 777 /var/lib/etcd/
mkdir /etc/etcd
ETCD_CONFIG_FILE="/etc/etcd/etcd.conf"
MY_IPV4=`ip addr show | grep 'inet ' | grep -v '127.0.0.1' | grep -v 'docker0' | grep -v 'flannel' | awk '$1 == "inet" {print $2}' | cut -d '/' -f 1`
MY_HOSTNAME=`hostname`
echo "name: $MY_HOSTNAME" > $ETCD_CONFIG_FILE
echo "data-dir: /var/lib/etcd/$MY_HOSTNAME.etcd" >> $ETCD_CONFIG_FILE
echo "initial-advertise-peer-urls: http://$MY_IPV4:2380" >> $ETCD_CONFIG_FILE
echo "listen-peer-urls: http://$MY_IPV4:2380" >> $ETCD_CONFIG_FILE
echo "listen-client-urls: http://$MY_IPV4:2379,http://127.0.0.1:2379,http://$MY_IPV4:4001,http://127.0.0.1:4001" >> $ETCD_CONFIG_FILE
echo "advertise-client-urls: http://$MY_IPV4:2379" >> $ETCD_CONFIG_FILE
echo "initial-cluster-token: etcd-cluster-1" >> $ETCD_CONFIG_FILE
echo "initial-cluster: $MY_HOSTNAME=http://$MY_IPV4:2380" >> $ETCD_CONFIG_FILE
echo "initial-cluster-state: new" >> $ETCD_CONFIG_FILE
echo "" >> $ETCD_CONFIG_FILE
## 创建/etc/systemd/system/etcd.service文件
ETCD_SERVICE_FILE="/etc/systemd/system/etcd.service"
echo "[Unit]" > $ETCD_SERVICE_FILE
echo "Description=etcd - highly-available key value store" >> $ETCD_SERVICE_FILE
echo "Documentation=GitHub - etcd-io/etcd: Distributed reliable key-value store for the most critical data of a distribu" >> $ETCD_SERVICE_FILE
echo "Documentation=man:etcd" >> $ETCD_SERVICE_FILE
echo "After=network.target" >> $ETCD_SERVICE_FILE
echo "Wants=network-online.target" >> $ETCD_SERVICE_FILE
echo "" >> $ETCD_SERVICE_FILE
echo "[Service]" >> $ETCD_SERVICE_FILE
echo "Environment=DAEMON_ARGS=" >> $ETCD_SERVICE_FILE
echo "Environment=ETCD_NAME=%H" >> $ETCD_SERVICE_FILE
echo "Environment=ETCD_DATA_DIR=/var/lib/etcd/default" >> $ETCD_SERVICE_FILE
echo "EnvironmentFile=-/etc/default/%p" >> $ETCD_SERVICE_FILE
echo "Type=notify" >> $ETCD_SERVICE_FILE
echo "User=root" >> $ETCD_SERVICE_FILE
echo "PermissionsStartOnly=true" >> $ETCD_SERVICE_FILE
echo "#ExecStart=/bin/sh -c \"GOMAXPROCS=$(nproc) /usr/bin/etcd \$DAEMON_ARGS\"" >> $ETCD_SERVICE_FILE
echo "ExecStart=/usr/bin/etcd \$DAEMON_ARGS --config-file=/etc/etcd/etcd.conf" >> $ETCD_SERVICE_FILE
echo "Restart=on-abnormal" >> $ETCD_SERVICE_FILE
echo "#RestartSec=10s" >> $ETCD_SERVICE_FILE
echo "LimitNOFILE=65536" >> $ETCD_SERVICE_FILE
echo "" >> $ETCD_SERVICE_FILE
echo "[Install]" >> $ETCD_SERVICE_FILE
echo "WantedBy=multi-user.target" >> $ETCD_SERVICE_FILE
echo "Alias=etcd.service" >> $ETCD_SERVICE_FILE
## 重新加载systemd配置并启动service系统级服务
systemctl daemon-reload
systemctl start etcd
# systemctl status etcd

# 网络CNI插件：Flannel
## 提前配置
echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
sysctl -p
iptables -P INPUT ACCEPT
iptables -P FORWARD ACCEPT
iptables -F
## 安装,由于最新flannel下载链接需要代理，故此处随工程发布了一个已经下载好的压缩包，解压即可
mkdir ./flannel
tar -zxvf ./plugin/flannel0.25.2/flannel-v0.25.2-linux-amd64.tar.gz -C ./flannel
## etcd中内容添加
etcdctl put /coreos.com/network/config '{ "Network": "10.5.0.0/16", "Backend": {"Type": "vxlan"}}'
## 创建flannle系统进程
cp ./flannel/flanneld /usr/bin/
FLANNEL_SERVICE_FILE="/etc/systemd/system/flannel.service"
echo "[Unit]" > $FLANNEL_SERVICE_FILE
echo "Description=flannel CNI" >> $FLANNEL_SERVICE_FILE
echo "" >> $FLANNEL_SERVICE_FILE
echo "[Service]" >> $FLANNEL_SERVICE_FILE
echo "User=root" >> $FLANNEL_SERVICE_FILE
echo "ExecStart=/usr/bin/flanneld " >> $FLANNEL_SERVICE_FILE
## 重新加载systemd配置并启动service系统级服务
systemctl daemon-reload
systemctl start flannel
# systemctl status flannel

# 虚拟化容器插件：Docker
## 安装
apt-get install ca-certificates curl -y
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update
apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
# 根据flannel配置，为docker重新分配子网段
# cat /run/flannel/subnet.env
# source /run/flannel/subnet.env
FLANNEL_SUBNET=`cat /run/flannel/subnet.env |grep SUBNET | awk -F'=' '{print $2}'`
FLANNEL_MTU=`cat /run/flannel/subnet.env |grep MTU | grep -oE '[0-9]+'`
DOCKER_DAEMON_FILE="/etc/docker/daemon.json"
echo "{" > $DOCKER_DAEMON_FILE
echo "  \"bip\": \"$FLANNEL_SUBNET\"," >> $DOCKER_DAEMON_FILE
echo "  \"mtu\": $FLANNEL_MTU" >> $DOCKER_DAEMON_FILE
echo "}" >> $DOCKER_DAEMON_FILE
echo "" >> $DOCKER_DAEMON_FILE
cat $DOCKER_DAEMON_FILE
systemctl restart docker
systemctl status docker

# 异步消息队列插件：RabbitMQ
# rabbitmq启动在docker中
docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management

# 语言运行时：Golang
rm -rf /usr/local/go
tar -C /usr/local -xzf ./plugin/go1.22.3/go1.22.3.linux-amd64.tar.gz


