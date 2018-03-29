# DHCP Notes

What I know over dhcpd in EL

* [Sample dhcpd config File](#sample-dhcpd-config-file)

## Sample dhcpd config File

Config where it won't allow unknown hosts to get an IP address. You will need to add a host name option in the configuration file with the mac address.

Edit the `/etc/dhcpd.conf` file and add the mac address. You can also set "static" dhcp with the `fixed-address` directive

```
        authoritative;
	ddns-update-style interim;
	default-lease-time 14400;
	max-lease-time 14400;

        	option routers                  192.168.2.1;
        	option broadcast-address        192.168.2.255;
        	option subnet-mask              255.255.255.0;
        	option domain-name-servers      192.168.2.225, 192.168.2.223;
        	option domain-name              “example.com”;

        	subnet 192.168.2.0 netmask 255.255.255.0 {
             	pool {
                	range 192.168.2.205 192.168.2.212;
                	# Systems laptop
			host systems-laptop { hardware ethernet 00:1c:26:03:6f:80; }
                        #
                        host chood-laptop { hardware ethernet 08:ed:b9:2d:82:4a; fixed-address 192.168.2.101; }
			#
			deny unknown-clients;
             	}
	}
```

 Just remember to put your entry ABOVE the `deny unknown-clients`; option. once you added an entry; restart the dhcpd service

```
root@dhcp-srv# service dhcpd restart
```

You can check current leases by veiwing the `/var/lib/dhcpd/dhcpd.leases` file with vi or cat

```
root@dhcp-srv# cat /var/lib/dhcpd/dhcpd.leases
```

Remember that the host option dosen’t have to match DNS. But the name has to be uniq to the file.

