## rabbitmq搭建

运行以下命令进行安装

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

访问 `http://(IP):156720` 即可进入rabbitmq控制面板 ，用户名和密码都是 `admin`