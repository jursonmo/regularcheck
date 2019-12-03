#### tcp keepalive 检测
```
net.ipv4.tcp_keepalive_time=7200
net.ipv4.tcp_keepalive_intvl=75
net.ipv4.tcp_keepalive_probes=9
```
每隔两小时心跳测量一次，如果失败，就75检测一次，连续9次都失败就认为断开连接了。
即第一次心跳检测失败后，降低心跳的检测间隔。这是在心跳的频繁程度和快速检测出连接超时的一种折中方式。
工作中也经常要用到这种方式的心跳检测。所以写个简单库：regularcheck。