# Sub interface


First find the interface you want to “split” using ifconfig (usually eth0 but make sure)

	root@host# ifconfig 
	eth0    Link encap:Ethernet  HWaddr 00:30:48:2E:33:84  
          	inet addr:192.168.11.93  Bcast:192.168.11.255  Mask:255.255.255.0
          	inet6 addr: fe80::230:48ff:fe2e:3384/64 Scope:Link
          	UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          	RX packets:24442366 errors:0 dropped:0 overruns:0 frame:0
          	TX packets:23894502 errors:0 dropped:0 overruns:0 carrier:0
          	collisions:0 txqueuelen:1000 
          	RX bytes:1495936958 (1.3 GiB)  TX bytes:2904282128 (2.7 GiB)

Next find the configuration file to “clone”: Usually under /etc/sysconfig/network-scripts

	root@host# cd /etc/sysconfig/network-scripts
	root@host# ll -d ifcfg-eth0
	-rw-r--r-- 1 root root 210 Dec 22 14:42 ifcfg-eth0

Copy the file and append a “:`<last-octect>`” to the file (the last octect of the IP you are using)

	root@host# cp “ifcfg-eth0” “ifcfg-eth0:95”
	root@host# ll -d ifcfg-eth0*
	-rw-r--r-- 1 root root 210 Dec 22 14:42 ifcfg-eth0
	-rw-r--r-- 1 root root 217 Jan 14 09:11 ifcfg-eth0:95

Edit the copied file & comment out the HWADDR and change the IPADDR to the IP you want to use; And change the DEVICE to a new interface (I put in bold what you need to change for easy reading)

	# Intel Corporation 82546GB Gigabit Ethernet Controller
	DEVICE=eth0:95
	BOOTPROTO=static
	BROADCAST=192.168.11.255
	####HWADDR=00:30:48:2E:33:84
	IPADDR=192.168.11.95
	NETMASK=255.255.255.0
	NETWORK=192.168.11.0
	ONBOOT=yes

If there is a UUID entry — comment it out. Now you can ifup the interface name to make it “live”

	root@host# ifup “eth0:95”

Since it’s set to // ONBOOT=yes // in the configuration file; this will “plumb” on boot. 
