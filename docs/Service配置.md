# ClusterIP iptables design

## 配置
### pod
- 主机cloudos1
  - os1_my1  10.5.75.2:3000
  - os1_my2  10.5.75.3:3000
- 主机cloudos2
  - os2_my1  10.5.88.2:3000
  - os2_my2  10.5.88.3:3000
### service
- service1 10.100.100.1.7070
  - os1_my1服务
  - os2_my1服务 
- service2 10.100.100.2:6060
  - os1_my2服务
  - os2_my2服务   

### iptables 命令(超级用户下)
#### filter table配置
``` sh

```

``` sh
# 查看配置规则
iptables -t filter -L -nv --line-number
```

#### nat table配置
``` sh
# 新建链
iptables -t nat -N KUBE-SERVICES
iptables -t nat -N KUBE-SVC-1
iptables -t nat -N KUBE-SVC-2
iptables -t nat -N KUBE-SEP-os1my1
iptables -t nat -N KUBE-SEP-os1my2
iptables -t nat -N KUBE-SEP-os2my1
iptables -t nat -N KUBE-SEP-os2my2

# 配置KUBE-MARK-MASQ
iptables -t nat -N KUBE-MARK-MASQ
iptables -t nat -A KUBE-MARK-MASQ -j MARK --set-mark 0x4000

# 配置KUBE-POSTROUTING
iptables -t nat -N KUBE-POSTROUTING
iptables -t nat -A POSTROUTING -j KUBE-POSTROUTING
iptables -t nat -A KUBE-POSTROUTING  -m mark --mark 0x4000 -j MASQUERADE

# 将KUBE-SERVICES插入REROUTING和OUTPUT
iptables -t nat -A PREROUTING -j KUBE-SERVICES
iptables -t nat -A OUTPUT -j KUBE-SERVICES

# 将KUBE-SVC-XXX插入KUBE-SERVICES中，并设定匹配规则
iptables -t nat -A KUBE-SERVICES -d 10.100.100.1 -p tcp -m tcp --dport 7070 -j KUBE-SVC-1
iptables -t nat -A KUBE-SERVICES -d 10.100.100.2 -p tcp -m tcp --dport 6060 -j KUBE-SVC-2

# 为KUBE-SVC-1引导至KUBE-SEP-os1my1和KUBE-SEP-os2my1以DNAT,并建立随机分配策略（负载均衡）
iptables -t nat -A KUBE-SVC-1 -m statistic --mode random --probability 0.5 -j KUBE-SEP-os1my1
iptables -t nat -A KUBE-SVC-1 -j KUBE-SEP-os2my1
# 为KUBE-SVC-2引导至KUBE-SEP-os1my2和KUBE-SEP-os2my2以DNAT,并建立随机分配策略（负载均衡）
iptables -t nat -A KUBE-SVC-2 -m statistic --mode random --probability 0.5 -j KUBE-SEP-os1my2
iptables -t nat -A KUBE-SVC-2 -j KUBE-SEP-os2my2

# 在KUBE-SEP-XXX进行SNAT准备
iptables -t nat -A KUBE-SEP-os1my1 -s 10.5.75.2 -j KUBE-MARK-MASQ
iptables -t nat -A KUBE-SEP-os1my2 -s 10.5.75.3 -j KUBE-MARK-MASQ
iptables -t nat -A KUBE-SEP-os2my1 -s 10.5.88.2 -j KUBE-MARK-MASQ
iptables -t nat -A KUBE-SEP-os2my2 -s 10.5.88.3 -j KUBE-MARK-MASQ

# 在KUBE-SEP-XXX进行实际的DNAT
iptables -t nat -A KUBE-SEP-os1my1 -p tcp -m tcp -j DNAT --to-destination 10.5.75.2:3000
iptables -t nat -A KUBE-SEP-os1my2 -p tcp -m tcp -j DNAT --to-destination 10.5.75.3:3000
iptables -t nat -A KUBE-SEP-os2my1 -p tcp -m tcp -j DNAT --to-destination 10.5.88.2:3000
iptables -t nat -A KUBE-SEP-os2my2 -p tcp -m tcp -j DNAT --to-destination 10.5.88.3:3000

# 需要在POSTROUTING中加入SNAT,否则从宿主机层面访问服务且负载均衡的pod不在本机上时，从flannle接口发出的包的source不符，无法寻址
iptables -t nat -I POSTROUTING 1 -s 192.168.1.6 -j MASQUERADE

```
> 对于最后一条SANT规则的效果，执行`tcpdump -ni flannel.1`，然后`curl 10.100.100.1：7070`，在负载均衡到外机时可以看到source IP的变化：  
> ```
> 11:27:01.379880 IP 192.168.1.6.56094 > 10.5.88.2.3000: Flags [S], seq 1416448032, win 64860, options [mss 1410,sackOK,TS val 3771774629 ecr 0,nop,wscale 7], length 0  
> 11:27:02.392892 IP 192.168.1.6.56094 > 10.5.88.2.3000: Flags [S], seq 1416448032, win 64860, options [mss 1410,sackOK,TS val 3771775642 ecr 0,nop,wscale 7], length 0  
> 11:27:04.408937 IP 192.168.1.6.56094 > 10.5.88.2.3000: Flags [S], seq 1416448032, win 64860, options [mss 1410,sackOK,TS val 3771777658 ecr 0,nop,wscale 7], length 0  
> 11:27:08.472861 IP 192.168.1.6.56094 > 10.5.88.2.3000: Flags [S], seq 1416448032, win 64860, options [mss 1410,sackOK,TS val 3771781722 ecr 0,nop,wscale 7], length 0  
> 11:27:16.664923 IP 192.168.1.6.56094 > 10.5.88.2.3000: Flags [S], seq 1416448032, win 64860, options [mss 1410,sackOK,TS val 3771789914 ecr 0,nop,wscale 7], length 0  
> 11:31:45.929381 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [S], seq 2113060412, win 65280, options [mss 1360,sackOK,TS val 4189490513 ecr 0,nop,wscale 7], length 0  
> 11:31:45.930865 IP 10.5.88.2.3000 > 10.5.75.0.55646: Flags [S.], seq 502294427, ack 2113060413, win 64704, options [mss 1360,sackOK,TS val 1393832131 ecr 4189490513,nop,wscale 7], length 0  
> 11:31:45.930917 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [.], ack 1, win 510, options [nop,nop,TS val 4189490515 ecr 1393832131], length 0  
> 11:31:45.930971 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [P.], seq 1:79, ack 1, win 510, options [nop,nop,TS val 4189490515 ecr 1393832131], length 78  
> 11:31:45.931188 IP 10.5.88.2.3000 > 10.5.75.0.55646: Flags [.], ack 79, win 505, options [nop,nop,TS val 1393832132 ecr 4189490515], length 0  
> 11:31:45.933121 IP 10.5.88.2.3000 > 10.5.75.0.55646: Flags [P.], seq 1:963, ack 79, win 505, options [nop,nop,TS val 1393832134 ecr 4189490515], length 962  
> 11:31:45.933157 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [.], ack 963, win 503, options [nop,nop,TS val 4189490517 ecr 1393832134], length 0  
> 11:31:45.933331 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [F.], seq 79, ack 963, win 503, options [nop,nop,TS val 4189490517 ecr 1393832134], length 0  
> 11:31:45.933536 IP 10.5.88.2.3000 > 10.5.75.0.55646: Flags [F.], seq 963, ack 80, win 505, options [nop,nop,TS val 1393832134 ecr 4189490517], length 0  
> 11:31:45.933564 IP 10.5.75.0.55646 > 10.5.88.2.3000: Flags [.], ack 964, win 503, options [nop,nop,TS val 4189490518 ecr 1393832134], length 0  
> ```

``` sh
# 查看配置规则
iptables -t nat -L -nv --line-number
```

## NodePort
NodePort所有配置基于ClusterIP的配置

### 暴露端口
- 主机cloudos1
  - service1 -- port 34211
  - service2 -- port 34212 
- 主机cloudos2
  - service1 -- port 34211
  - service2 -- port 34212

### nat table配置
``` sh
# 创建新链KUBE-NODEPORTS
iptables -t nat -N KUBE-NODEPORTS
# KUBE-NODEPORTS插入KUBE-SERVICES最后一行,在所有的ClusterIP之后处理
iptables -t nat -A KUBE-SERVICES -m addrtype --dst-type LOCAL -j KUBE-NODEPORTS
# 在KUBE-NODEPORTS中添加规则以保证SNAT
iptables -t nat -A KUBE-NODEPORTS -p tcp -m tcp --dport 34211 -j KUBE-MARK-MASQ
iptables -t nat -A KUBE-NODEPORTS -p tcp -m tcp --dport 34212 -j KUBE-MARK-MASQ
# 在KUBE-NODEPORTS中添加规则，引导不同的端口访问流向不同的服务链
iptables -t nat -A KUBE-NODEPORTS -p tcp -m tcp --dport 34211 -j KUBE-SVC-1
iptables -t nat -A KUBE-NODEPORTS -p tcp -m tcp --dport 34212 -j KUBE-SVC-2
```
``` sh
# 查看配置规则
iptables -t nat -L -nv --line-number
```