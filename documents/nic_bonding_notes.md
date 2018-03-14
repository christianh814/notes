# NIC Bonding Notes

How to create a bonding interface

* [EL 6](#el-6)
* [EL 7](#el-7)

## EL 6

The Linux bonding driver provides a method for aggregating multiple network interfaces into a single logical “bonded” interface. The behavior of the bonded interfaces depends upon the mode; generally speaking, modes provide either hot standby or load balancing services. Additionally, link integrity monitoring may be performed.

* [Create a Bond0 Configuration File](#create-a-bond0-configuration-file)
* [Modify eth0 and eth1 config files](#modify-eth0-and-eth1-config-files)
* [Load bond driver and module](#load-bond-driver-and-module)
* [Test configuration](#test-configuration)

### Create a Bond0 Configuration File

Red Hat Enterprise Linux (and its clone such as CentOS) stores network configuration in /etc/sysconfig/network-scripts/ directory. First, you need to create a bond0 config file as follows:

```
root@host# vi /etc/sysconfig/network-scripts/ifcfg-bond0
```

Append the following lines:

```
DEVICE=bond0
IPADDR=192.168.1.20
NETWORK=192.168.1.0
NETMASK=255.255.255.0
USERCTL=no
BOOTPROTO=none
ONBOOT=yes
MTU=9000
BONDING_OPTS="mode=4 miimon=100"
```

You need to replace IP address with your actual setup. The `MTU=9000` is optional (if you want "jumbo" TCP Frames). The `BONDING_OPTS="mode=4 miimon=100"` means that you're using "mode 4" balancing (compatibale with most switches) and to check health every 100 milliseconds. It specifies (in milliseconds) how often MII link monitoring occurs. This is useful if high availability is required because MII is used to verify that the NIC is active. Save and close the file.

### Modify eth0 and eth1 config files

Open both configuration using a text editor such as vi/vim, and make sure file read as follows for eth0 interface

```
root@host# vi /etc/sysconfig/network-scripts/ifcfg-eth0
```

Modify/append directive as follows:

```
DEVICE=eth0
HWADDR=00:10:18:F3:30:14
USERCTL=no
ONBOOT=yes
MASTER=bond0
SLAVE=yes
BOOTPROTO=none
```

Open eth1 configuration file using vi text editor, enter:

```
root@host# vi /etc/sysconfig/network-scripts/ifcfg-eth1
```

Make sure file read as follows for eth1 interface:

```
DEVICE=eth1
HWADDR="00:10:18:F3:30:16"
USERCTL=no
ONBOOT=yes
MASTER=bond0
SLAVE=yes
BOOTPROTO=none
```
Save and close the file.

### Load bond driver and module

Make sure bonding module is loaded when the channel-bonding interface (bond0) is brought up. You need to modify kernel modules configuration file:

```
root@# vi /etc/modprobe.d/bonding.conf
```

Add the following lines:

```
#
# Add "bond" alias
alias bond0 bonding
#
#
```
Save file and exit to shell prompt.

### Test configuration

First, load the bonding module, enter:

```
root@host# modprobe bonding
```

Restart the networking service in order to bring up bond0 interface, enter:

```
root@host# service network restart
```

Make sure everything is working. Type the following cat command to query the current status of Linux kernel bounding driver, enter:

```
root@host# cat /proc/net/bonding/bond0
```

## EL 7

The steps above will work with EL7 - but as an alternative `nmcli` offers [network teaming](nmcli_notes.md#configuring-networking-teaming)
