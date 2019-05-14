# Network Manager Notes

These are notes for NetworkMangaer in no paticular order

* [Configuring Networking with nmcli](#configuring-networking-with-nmcli)
* [Configuring Networking Manually](#configuring-networking-manually)
* [Configuring Hostname](#configuring-hostname)
* [Configuring Networking Teaming](#configuring-networking-teaming)
* [Configuring Software Bridges ](#configuring-software-bridges)
* [Misc Notes](#misc)

## Configuring Networking with nmcli

The `NetworkManager` is a daemon that monitors and manages network settings. In addition to the daemon, there is a GNOME Notification Area applet that provides network status information. Command-line and graphical tools talk to `NetworkManage` and save configuration files in the `/etc/sysconfig/network-scripts` directory.

In general...
  * `device` - a network interface such as `eth0` or `wpls0`
  * `connection` - a configuration used for a device which is made up of a collection of settings


Multiple connections may exist for a device, but only one may be active at a time. For example, a system may normally be connected to a network with settings provided by DHCP. Occasionally, that system needs to be connected to a lab or data center network, which only uses static networking. Instead of changing the configuration manually, each configuration can be stored as a separate connection.

* [Create a static connection](#create-a-static-connection)
* [Modify the new connection](#modify-the-new-connection)
* [Display and activate the new connection](#display-and-activate-the-new-connection)
* [Test the connectivity using the new network addresses](#test-the-connectivity-using-the-new-network-addresses)
* [Disable Autoconnect](#disable-autoconnect)

### View network settings using nmcli

Show all connections.
```
[user@host ~]$ nmcli con show
NAME         UUID                                  TYPE            DEVICE
System eth0  5fb06bd0-0bb0-7ffb-45f1-d6edd65f3e03  802-3-ethernet  eth0
```

Display all configuration settings for the active connection.
```
[user@host ~]$ nmcli con show "System eth0"
connection.id:                          System eth0
connection.uuid:                        5fb06bd0-0bb0-7ffb-45f1-d6edd65f3e03
connection.interface-name:              eth0
connection.type:                        802-3-ethernet
connection.autoconnect:                 yes
connection.timestamp:                   1394813303
connection.read-only:                   no
connection.permissions:
...
IP4.ADDRESS[1]:                         ip = 172.25.45.11/16, gw = 172.25.0.1
IP4.DNS[1]:                             172.25.254.254
IP4.DOMAIN[1]:                          example.com
...
```

Show device status.

```
[user@host ~]$ nmcli dev status
DEVICE  TYPE      STATE      CONNECTION
eth0    ethernet  connected  System eth0
lo      loopback  unmanaged  --
```
Display the settings for the eth0 device.

```
[user@host ~]$ nmcli dev show eth0
GENERAL.DEVICE:                         eth0
GENERAL.TYPE:                           ethernet
GENERAL.HWADDR:                         52:54:00:00:00:0B
GENERAL.MTU:                            1500
GENERAL.STATE:                          100 (connected)
GENERAL.CONNECTION:                     System eth0
GENERAL.CON-PATH:                       /org/freedesktop/NetworkManager/ActiveConnection/1
WIRED-PROPERTIES.CARRIER:               on
IP4.ADDRESS[1]:                         ip = 172.25.45.11/16, gw = 172.25.0.1
IP4.DNS[1]:                             172.25.254.254
IP4.DOMAIN[1]:                          example.com
IP6.ADDRESS[1]:                         ip = fe80::5054:ff:fe00:b/64, gw = ::
```

### Create a static connection

Use the `con add` option to create a new connection. Assign an IPv4 address, network prefix, and default gateway. Name the new connection `static-eth0`

```
[root@host ~]# nmcli con add con-name "static-eth0" ifname eth0 type ethernet ip4 172.25.45.11/16 gw4 172.25.0.1
Connection 'static-eth0' (f3e8dd32-3c9d-48f6-9066-551e5b6e612d) successfully added.
```

### Modify the new connection

You can modify properties of the new connection (ip address, gateway, DNS settings, etc) using the `con mod` arugment. To add the DNS setting edit the `ipv4.dns` property

```
[root@host ~]# nmcli con mod "static-eth0" ipv4.dns 172.25.254.254
```

### Display and activate the new connection

View all connections.

```
[root@host ~]# nmcli con show
NAME         UUID                                  TYPE            DEVICE
static-eth0  f3e8dd32-3c9d-48f6-9066-551e5b6e612d  802-3-ethernet  --
System eth0  5fb06bd0-0bb0-7ffb-45f1-d6edd65f3e03  802-3-ethernet  eth0
```

View the active connection.

```
[root@host ~]# nmcli con show --active
System eth0  5fb06bd0-0bb0-7ffb-45f1-d6edd65f3e03  802-3-ethernet  eth0
```

Activate the new connection.

```
[root@host ~]# nmcli con up "static-eth0"
Connection successfully activated (D-Bus active path: /org/freedesktop/NetworkManager/ActiveConnection/3)
```

View the active connection.

```
[root@host ~]# nmcli con show --active
NAME         UUID                                  TYPE            DEVICE
static-eth0  f3e8dd32-3c9d-48f6-9066-551e5b6e612d  802-3-ethernet  eth0
```

### Test the connectivity using the new network addresses

Verify the IP address.

```
[root@host ~]# ip addr show eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 52:54:00:00:00:0b brd ff:ff:ff:ff:ff:ff
    inet 172.25.45.11/16 brd 172.25.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::5054:ff:fe00:b/64 scope link
       valid_lft forever preferred_lft forever
```

Verify the default gateway.

```
[root@host ~]# ip route
default via 172.25.0.1 dev eth0  proto static  metric 1024
```

Ping the DNS address.

```
[root@host ~]# ping -c3 172.25.254.254
PING 172.25.254.254 (172.25.254.254) 56(84) bytes of data.
64 bytes from 172.25.254.254: icmp_seq=1 ttl=64 time=0.419 ms
64 bytes from 172.25.254.254: icmp_seq=2 ttl=64 time=0.598 ms
64 bytes from 172.25.254.254: icmp_seq=3 ttl=64 time=0.503 ms

--- 172.25.254.254 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 1999ms
rtt min/avg/max/mdev = 0.419/0.506/0.598/0.077 ms
```

### Disable Autoconnect

Configure the original connection so that it does not start at boot and verify that the static connection is used when the system reboots.

Disable the original connection from autostarting at boot.

```
[root@host ~]# nmcli con mod "System eth0" connection.autoconnect no
```

Reboot the system.

```
[root@host ~]# reboot
```

View the active connection.

```
[root@host ~]# nmcli con show --active
NAME         UUID                                  TYPE            DEVICE
static-eth0  f3e8dd32-3c9d-48f6-9066-551e5b6e612d  802-3-ethernet  eth0
```

## Configuring Networking Manually

It is also possible to configure the network by editing interface configuration files. Interface configuration files control the software interfaces for individual network devices. These files are usually named `/etc/sysconfig/network-scripts/ifcfg-<name>`, where `<name>` refers to the name of the device or connection that the configuration file controls.

* [Configuration Options for ifcfg File](#configuration-options-for-ifcfg-file)
* [Reloading Config](#reloading-config)

### Configuration Options for ifcfg File

The following are standard variables found in the file used for static or dynamic configuration.

![images/table.nmcli.jpg](table.nmcli.jpg)

In the static settings, variables for IP address, prefix, and gateway have a number at the end. This allows multiple sets of values to be assigned to the interface. The DNS variable also has a number which is used to specify the order of lookup when multiple servers are specified.

Good documentation can be found on the system in the `/usr/share/doc/initscripts-*/sysconfig.txt` file

### Reloading Config

After modifying the configuration files, run `nmcli con reload` to make `NetworkManager` read the configuration changes. The interface still needs to be restarted for changes to take effect.

 ```
[root@host ~]# nmcli con reload
[root@host ~]# nmcli con down "System eth0"
[root@host ~]# nmcli con up "System eth0"
```


## Configuring Hostname

 The hostname command displays or temporarily modifies the system's fully qualified host name.

```
[root@host ~]# hostname
```

A static host name may be specified in the `/etc/hostname` file.

### Using hostnamectl

Display the host name status.

```
[root@host ~]# hostnamectl status
   Static hostname: n/a
Transient hostname: host.example.com
         Icon name: computer
           Chassis: n/a
        Machine ID: 9f6fb63045a845d79e5e870b914c61c9
           Boot ID: d4ec3a2e8d3c48749aa82738c0ea946a
  Operating System: Red Hat Enterprise Linux Server 7.0 (Maipo)
       CPE OS Name: cpe:/o:redhat:enterprise_linux:7.0:beta:server
            Kernel: Linux 3.10.0-97.el7.x86_64
      Architecture: x86_64
```

Set a static host name to match the current transient host name.

Change the host name and host name configuration file.

```
[root@host ~]# hostnamectl set-hostname host.example.com
```
View the configuration file providing the host name at network start.

```
[root@host ~]# cat /etc/hostname
serverX.example.com
```

Display the host name status.

```
[root@host ~]# hostnamectl status
   Static hostname: serverX.example.com
         Icon name: computer
           Chassis: n/a
        Machine ID: 9f6fb63045a845d79e5e870b914c61c9
           Boot ID: d4ec3a2e8d3c48749aa82738c0ea946a
  Operating System: Red Hat Enterprise Linux Server 7.0 (Maipo)
       CPE OS Name: cpe:/o:redhat:enterprise_linux:7.0:beta:server
            Kernel: Linux 3.10.0-97.el7.x86_64
      Architecture: x86_64
```

## Configuring Networking Teaming

Network teaming is method for linking NICs together logically to allow for failover or higher throughput. Teaming is a new implementation that does not affect the older nic_bonding_notes in the Linux kernel; it offers an alternate implementation. Red Hat Enterprise Linux 7 supports channel bonding for backward compatability. Network teaming provides better performance and is more extensible because of its modular design.

Red Hat Enterprise Linux 7 implements network teaming with a small kernel driver and a user-space daemon, `teamd` The kernel handles network packets efficiently and `teamd` handles logic and interface processing. Software, called `runners` implement load balancing and active-backup logic, such as roundrobin. The following `runners` are available to `teamd`:

* `broadcast` a simple runner which transmits each packet from all ports.
* `roundrobin` a simple runner which transmits packets in a round-robin fashing from each of the ports.
* `activebackup` this is a failover runner which watches for link changes and selects an active port for data transfers.
* `loadbalance` this runner monitors traffic and uses a hash function to try to reach a perfect balance when selecting ports for packet transmission.
* `lacp` implements the 802.3ad Link Aggregation Control Protocol. Can use the same transmit port selection possibilities as the loadbalance runner.

Steps are outlined as:

* [Create Team Interface](#create-team-interface)
* [Assign Interfaces To Team Interface](#assign-interfaces-to-team-interface)
* [Team Interface Management](#team-interface-management)
* [Configuration Settings](#configuration-settings)

### Create Team Interface

Create an `active-backup` teaming interface called `team0` and assign its IPv4 settings. (although most likely you'll use the `round-robin` runner config)

```
root@host# nmcli con add type team con-name team0 ifname team0 config '{"runner": {"name": "activebackup"}}'
```

Assign an ipv4 address to the interface

```
root@host# nmcli con mod team0 ipv4.addresses '192.168.0.100/24'
root@host# nmcli con mod team0 ipv4.method manual
```


### Assign Interfaces To Team Interface

Next Assign port interfaces for `team0` Use the `team-slave` type.

```
root@host# nmcli con add type team-slave con-name team0-port1 ifname eno1 master team0
root@host# nmcli con add type team-slave con-name team0-port2 ifname eno2 master team0
```

Now check the state

```
root@host# teamdctl team0 state
setup:
  runner: activebackup
ports:
  eno1
    link watches:
      link summary: up
      instance[link_watch_0]:
        name: ethtool
        link: up
  eno2
    link watches:
      link summary: up
      instance[link_watch_0]:
        name: ethtool
        link: up
runner:
  active port: eno1
```

### Team Interface Management

The following are management commands that should be self explanatory - so I put a short description

Disable interface device

```
root@host# nmcli dev dis eno1
```

Bring up port conneciton

```
root@host# nmcli con up team0-port1
```

### Configuration Settings

Initial network team configuration is set when the team interface is created. The default runner is roundrobin, but a different runner can be chosen by specifying a JSON string when the team is created with the team.config subcommand. Default values for runner parameters are used when they are not specified.
A different runner can be assigned to an existing team, or runner parameters can be adjusted using the nmcli con mod command. The configuration changes can be specified as a JSON string (in the case of simple changes) or the name of a file with a more complex JSON configuration can be given.

 `nmcli con mod <IFACE> team.config <JSON-configuration-file-or-string>`


The following example shows how to assign different priorities to port interfaces in an active-backup team:

```
root@host#  cat /tmp/team.conf
{
    "device": "team0",
    "mcast_rejoin": {
        "count": 1
    },
    "notify_peers": {
        "count": 1
    },
    "ports": {
        "eth1": {
	    "prio": -10,
	    "sticky": true,
            "link_watch": {
                "name": "ethtool"
            }
        },
        "eth2": {
	    "prio": 100,
            "link_watch": {
                "name": "ethtool"
            }
        }
    },
    "runner": {
        "name": "activebackup"
    }
}
root@host# nmcli con mod team0 team.config /tmp/team.conf
```

Note: Any changes made do not go into effect until the next time the team interface is brought up.

The link_watch settings in the configuration file determines how the link state of the port interfaces are monitored. The default looks like the following, and uses functionality similar to the `ethtool` command to check the link of each interface:

```
"link_watch": {
    "name": "ethtool"
}
```

Another way to check link state is to periodically use an ARP ping packet to check for remote connectivity. Local and remote IP addresses and timeouts would have to be specified. A configuration that would accomplish that would look similar to the following:

```
"link_watch":{
    "name": "arp_ping",
    "interval": 100,
    "missed_max": 30,
    "source_host": "192.168.23.2",
    "target_host": "192.168.23.1"
},
```

Note: Be aware that omitted options revert to their default values when they are not specified in the JSON file.

The teamnl and teamdctl commands are very useful for troubleshooting network teams. These commands only work on network teams that are up. The following examples show some typical uses for these commands.
Display the team ports of the team0 interface:

```
root@host# teamnl team0 ports
 4: eth2: up 0Mbit HD
 3: eth1: up 0Mbit HD
```

Display the currently active port of team0:

```
root@host# teamnl team0 getoption activeport
3
```

Set the option for the active port of team0:

```
root@host# teamnl team0 setoption activeport 3
```


Use teamdctl to display the current JSON configuration for team0:

```
root@host# teamdctl team0 config dump
{
    "device": "team0",
    "mcast_rejoin": {
        "count": 1
    },
    "notify_peers": {
        "count": 1
    },
    "ports": {
        "eth1": {
            "link_watch": {
                "name": "ethtool"
            },
            "prio": -10,
            "sticky": true
        },
        "eth2": {
            "link_watch": {
                "name": "ethtool"
            },
            "prio": 100
        }
    },
    "runner": {
        "name": "activebackup"
    }
}
```

## Configuring Software Bridges

A network bridge is a link-layer device that forwards traffic between networks based on MAC addresses. It learns what hosts are connected to each network, builds a table of MAC addresses, then makes packet forwarding decisions based on that table.

A software bridge can be used in a Linux environment to emulate a hardware bridge. The most common application for software bridges is in virtualization applications for sharing a hardware NIC among one or more virtual NICs.

* [Create Bridge](#create-bridge)
* [Assign Bridge IP Address](#assign-bridge-ip-address)
* [Assign Interface To Bridge](#assign-interface-to-bridge)

### Create Bridge

To create a software bridge

```
root@host# nmcli con add con-name br0 type bridge ifname br0
```

### Assign Bridge IP Address

To assign the bridge IP address

```
root@host# nmcli con mod br0 ipv4.addresses "192.168.0.5/24"
root@host# nmcli con mod br0 ipv4.method manual
```

### Assign Interface To Bridge

Assign an interface to that network bridge

```
root@host# nmcli con add con-name br0-port1 type bridge-slave ifname eno1 master br0
```

After that bring up the interface

```
root@host# nmcli con up br0
```

Test with `brctl`

```
root@host# brctl show
```

## Misc

This is how you set up a DHCP interface with `PEERDNS=no`

```
nmcli con mod "Wired connection 1" ipv4.dns "192.168.1.2 192.168.1.3"
nmcli con mod "Wired connection 1" ipv4.ignore-auto-dns yes
```
