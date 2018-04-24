# Linux Routing

Linux Routing Notes in no paticular order

* [Static Routes](#static-routes)
* [Persistant Static Routes](#persistant-static-routes)
* [Linux As A Router](#linux-as-a-router)

# Static Routes


In Linux  you would use the ip route command

Use `ip route add` command like this

```
root@host# ip route add 192.168.11.0/24 via 192.168.11.1 dev eth1
```

You can also add it like this

```
root@host# route add -net 192.168.11.0 netmask 255.255.255.0 gw 192.168.11.1 dev eth1
```

Use `system-config-network` to add routes if you are using GUI (These routes will be permanent).
  * On a terminal run `system-config-network`
  * On the 'Network Configuration' window select the interface. For example - Nickname - eth0 and click on "Edit" button.. Click on "Route" tab.
  * Click on "Route" tab.
  * Click on "Add" and specify "Address, Subnet mask, Gateway" and then 'OK'.
  * On main menu click on "File" menu and Save.


# Persistant Static Routes

If you have added a second interface (`eth1`) and have static routes set up - for example...

```
root@host# ip route add 192.168.117.0/24 via 192.168.117.1 dev eth1
root@host# ip route add 192.168.217.0/24 via 192.168.117.1 dev eth1
```

You can make these routes persistant on boot...create the `/etc/sysconfig/network-scripts/route-eth1` and put your static routes there. For each route incement the number.

Example

```
ADDRESS0=192.168.117.0
NETMASK0=255.255.255.0
GATEWAY0=192.168.117.1
ADDRESS1=192.168.217.0
NETMASK1=255.255.255.0
GATEWAY1=192.168.117.1
```

Also make sure routing is turned on in the kernel - in the `/etc/sysctl.conf` file make sure the `net.ipv4.ip_forward` looks like this...

```
net.ipv4.ip_forward = 1
```

Then reload with...

```
root@host# sysctl -p
```

# Linux As A Router

Specific steps for different OSes below

* [RHEL 6](#rhel-6)
* [RHEL 7](#rhel-7)

## RHEL 6

Linux can be easily configured to share an internet connection using iptables. All you need to two network interface cards.

Example:
  * Your internal (LAN) network connected via eth0 with static ip address 192.168.1.254
  * Your external (WAN) network is connected via eth1 with static ip address 192.168.2.1


Login as the root user. Open `/etc/sysctl.conf` file and add/edit the following entry...

```
net.ipv4.conf.default.forwarding=1
```

Restart networking

```
root@host# service network restart
```

Next Enable IP masquerading (via IPTABLES)

```
root@host# service iptables stop
root@host# iptables -t nat -A POSTROUTING -o eth1 -j MASQUERADE
root@host# service iptables save
root@host# service iptables restart
```

This assumes that the interface eth1 may have public IP address or IP assigned by ISP. eth1 may be connected to a dedicated DSL / ADSL / WAN / Cable router.

## RHEL 7

These are quicknotes

Add the following to `/etc/sysctl.conf`

```
net.ipv4.ip_forward = 1
```

Then run the following
```
sysctl -p
```

Add `direct` rules to `firewalld`. Add the `--permanent` option to keep these rules across restarts.

```
firewall-cmd --direct --add-rule ipv4 nat POSTROUTING 0 -o eth_ext -j MASQUERADE
firewall-cmd --direct --add-rule ipv4 filter FORWARD 0 -i eth_int -o eth_ext -j ACCEPT
firewall-cmd --direct --add-rule ipv4 filter FORWARD 0 -i eth_ext -o eth_int -m state --state RELATED,ESTABLISHED -j ACCEPT
```
