# How To Configure A Network Bridge In Red Hat / Fedora

A network bridge is a forwarding technique very useful when you have to deal with virtualization and you want to give your virtual machines direct access to your real network, without using NAT.

In this example, I’m going to use a bridge (br0) to access a wired network interface (eth0).

__Installation__

Use yum to install the packages needed: 

	
	root@host# yum install bridge-utils


__Configuration__

Edit // /etc/sysconfig/network-scripts/ifcfg-eth0 // and write this (changing the **HWADDR** for the MAC address of your network card):

	
	DEVICE=eth0
	HWADDR=00:11:22:33:44:55
	ONBOOT=yes
	BRIDGE=br0


Edit // /etc/sysconfig/network-scripts/ifcfg-br0 // with this content (change the IP related fields to fit your needs):

	
	DEVICE=br0
	TYPE=Bridge
	ONBOOT=yes
	DELAY=0
	BOOTPROTO=static
	BROADCAST=192.168.1.255
	IPADDR=192.168.1.100
	NETMASK=255.255.255.0
	NETWORK=192.168.1.0
	GATEWAY=192.168.1.1


Add these lines to // /etc/sysctl.conf // in order to disable packet filtering in the bridge:

	
	net.bridge.bridge-nf-call-ip6tables = 0
	net.bridge.bridge-nf-call-iptables = 0
	net.bridge.bridge-nf-call-arptables = 0`</code>`
	
	This improves the bridge’s performance. I recommend to use packet filtering in the computers which connect through the bridge, but not in the bridge itself.
	
	Apply the syscttl changes: 
	
	`<code>`
	root@host# sysctl -p /etc/sysctl.conf`</code>`
	
	Restart your network interfaces:
	
	`<code>`
	root@host# service network restart`</code>`
	