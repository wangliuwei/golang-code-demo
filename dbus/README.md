使用dbus api操作linux防火墙
Firewalld是一个基于动态区域的防火墙守护进程，自 2009 年左右开始开发，目前为Fedora 18 以及随后的 RHEL7 和 CentOS 7 中的默认防火墙机制。

Firewalld被配置为systemd D-Bus 服务。请注意下面的“Type=dbus”指令。

cat /usr/lib/systemd/system/firewalld.service

[Unit]
Description=firewalld - dynamic firewall daemon
Before=network.target
Before=libvirtd.service
Before=NetworkManager.service
Conflicts=iptables.service ip6tables.service ebtables.service

[Service]
EnvironmentFile=-/etc/sysconfig/firewalld
ExecStart=/usr/sbin/firewalld --nofork --nopid $FIREWALLD_ARGS
ExecReload=/bin/kill -HUP $MAINPID
//# supress to log debug and error output also to /var/log/messages
StandardOutput=null
StandardError=null
Type=dbus
BusName=org.fedoraproject.FirewallD1

[Install]
WantedBy=basic.target
Alias=dbus-org.fedoraproject.FirewallD1.service
 

知道了firewalld服务是基于D-Bus的，就可以通过D-Bus来操作防火墙。

查看dbus注册的服务是否包含firewalld，这里需要注意的是，firewalld依赖dbus服务，每次启动firewalld时注册到dbus总线内。所以需要先启动​​dbus-daemon​​与 ​​firewalld ​​ 服务。

dbus-send --system --dest=org.freedesktop.DBus --type=method_call --print-reply \
/org/freedesktop/DBus org.freedesktop.DBus.ListNames | grep FirewallD
 

查看得知 ​​org.fedoraproject.FirewallD1​​ 为firewalld接口

查看接口所拥有的方法、属性、信号等信息

dbus-send --system --dest=org.fedoraproject.FirewallD1 --print-reply \
/org/fedoraproject/FirewallD1 org.freedesktop.DBus.Introspectable.Introspect
 

获得zone

firewall-cmd --get-zones

dbus-send --system \
--dest=org.fedoraproject.FirewallD1 \
--print-reply \
--type=method_call /org/fedoraproject/FirewallD1 \
org.fedoraproject.FirewallD1.zone.getZones
 

查看zone内的条目信息

firewall-cmd --zone=public --list-all

dbus-send --system --dest=org.fedoraproject.FirewallD1 --print-reply --type=method_call \
/org/fedoraproject/FirewallD1 org.fedoraproject.FirewallD1.getZoneSettings string:"public"
