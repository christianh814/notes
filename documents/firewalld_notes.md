# FirewallD Notes

These are my `firewalld` in no paticular order

* [Overview](#overview)
* [Change Firewall Rules On Boot](#change-firewall-rules-on-boot)
* [Misc Commands](#misc-commands)

## Overview

In future releases of Red Hat — `iptables` will be replaced with `firewalld`

These are just some basic commands that will “get you up and running”

Start/Stop/Status/Restart the service

```
root@host# systemctl start firewalld.service
root@host# systemctl stop firewalld.service
root@host# systemctl status firewalld.service
root@host# systemctl restart firewalld.service
```

Enable/Disable the firewall from starting on boot

```
root@host# systemctl enable firewalld.service
root@host# systemctl disable firewalld.service
```

Allow a service through the firewall (only current session; but does not persist through reboot)

```
root@host# firewall-cmd --add-service=ssh
root@host# firewall-cmd --add-port=80/tcp
```

Allow a service through the firewall for 10 seconds only (only current session; but does not persist through reboot)

```
root@host# firewall-cmd --add-service=samba --timeout=10
```

Disallow a service through the firewall (only current session; but does not persist through reboot)

```
root@host# firewall-cmd --remove-service=ipp-client
```

## Change Firewall Rules On Boot

In `firewalld` it loads a default configuration (`public.xml`) allowing only certian ports on boot

These configurations live under...

```
root@host# ll /usr/lib/firewalld/
drwxr-x---. 2 root root 4096 Nov 28 15:10 icmptypes
drwxr-x---. 2 root root 4096 Nov 28 15:10 services
drwxr-x---. 2 root root 4096 Nov 28 15:10 zones
```

Now in the `/etc/firewalld/zones/` directory you can place an XML that "overwrites" the defaults with whatever is in there. So to be "safe" we copy over the "default" and modify it

```
root@host# cp /usr/lib/firewalld/zones/public.xml /etc/firewalld/zones/
root@host# ll /etc/firewalld/zones/
-rw-r-----. 1 root root 375 Nov 29 15:08 public.xml
```

Now change the file from this...

```
<?xml version="1.0" encoding="utf-8"?>
<zone>
<short>Public<short>
<description>For use in public areas. You do not trust the other computers on networks to not harm your computer. Only selected incoming connections are accepted.<description>
<service name="ssh"/>
<service name="mdns"/>
<service name="dhcpv6-client"/>
<zone>
```

...to this

```
<?xml version="1.0" encoding="utf-8"?>
<zone>
<short>Public<short>
<description>For use in public areas. You do not trust the other computers on networks to not harm your computer. Only selected incoming connections are accepted.<description>
<service name="ssh"/>
<service name="mdns"/>
<service name="dhcpv6-client"/>
<port protocol="tcp" port="80"/>
<zone>
```

Adding whatever port/protocol combonation you need

You can do this with the command line using the `--permanent` switch (reload after)

```
root@host# firewall-cmd --permanent --add-service=ssh
root@host# firewall-cmd --permanent --add-port=80/tcp
root@host# firewall-cmd --reload
```

If you want to add it to a specific zone...

```
root@host# firewall-cmd --permanent --zone=public --add-service=ssh
root@host# firewall-cmd --permanent --zone=public --add-port=80/tcp
root@host# firewall-cmd --reload
```

## Misc Commands

__List Default Zone__

List the currently set default zone
```
[root@fedoraxp ~]# firewall-cmd --get-default-zone
public
```

__List Interface Attached to Zone__

Different interfaces can have differnent rules based on Zones
```
[root@fedoraxp ~]# firewall-cmd --zone=public --list-interfaces
p3p1
```

__List Zones__

List all zones

```
[root@fedoraxp ~]# firewall-cmd --list-all-zones
```

__List Services__

List current services allowed in a specific zone
```
[root@fedoraxp ~]# firewall-cmd --zone=public --list-services
mdns dhcpv6-client ssh
```

List current services in the current zone

```
[root@fedoraxp ~]# firewall-cmd --list-services
mdns dhcpv6-client ssh
```

More detail

```
[root@fedoraxp ~]# firewall-cmd --list-all
public
  interfaces: p3p1
  services: mdns dhcpv6-client ssh
  ports:
  forward-ports:
  icmp-blocks:
```

__Rich Rules__

For more "complex" rules; you'll have to use rich-rules
```
firewall-cmd --add-rich-rule="rule family="ipv4" source address="10.0.0.0/8" service name="ssh" accept"
```

-30-
