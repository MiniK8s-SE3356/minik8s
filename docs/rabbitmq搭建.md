# Rabbitmq搭建

Rabbmitmq后续作为一个docker容器启动,启动指令：
``` shell
docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

<!-- 运行以下命令进行安装

```
sudo apt-get install erlang
sudo apt-get install rabbitmq-server
```

启动控制面板

```
sudo rabbitmq-plugins enable rabbitmq_management

```

添加用户并设置权限

```
sudo rabbitmqctl add_user admin admin
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

访问 `http://(IP):15672` 即可进入rabbitmq控制面板 ，用户名和密码都是 `admin` -->