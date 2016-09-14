# HAProxy Basics

HAProxy, which stands for High Availability Proxy, is a popular open source software TCP/HTTP Load Balance, proxying solution, and SSL Offloading which can be run on Linux, Solaris, and FreeBSD. Its most common use is to improve the performance and reliability of a server environment by distributing the workload across multiple servers (e.g. web, application, database). It is used in many high-profile environments, including: GitHub, Imgur, Instagram, and Twitter.

I ran this on an EL7 machine but I assume the steps will be the same on any system as long as the haproxy version is the same

Pretty good "how to" found [HERE](https///www.digitalocean.com/community/tutorials/an-introduction-to-haproxy-and-load-balancing-concepts)

## Installation

It was part of the default YUM Repo

	
	[root@sandbox ~]# yum -y install haproxy


Make sure it starts on boot

	
	[root@sandbox ~]# systemctl enable haproxy
	[root@sandbox ~]# systemctl start haproxy


## Configuration

You basically need 3 configurations...


*  ***frontend*** - A frontend defines how requests should be forwarded to backend. Their definitions are composed of the following
    * a set of IP addresses and a port (e.g. 10.1.1.7:80, *:443, etc.)
    * ACLs
    * *use_backend rules*, which define which backends to use depending on which ACL conditions are matched, and/or a *default_backend* rule that handles every other case

*  ***backend*** - A backend is a set of servers that receives forwarded requests. In its most basic form, a backend can be defined by
    * which load balance algorithm to use (few of the commonly used algorithms below)
      * *roundrobin* - Round Robin selects servers in turns. This is the default algorithm.
      * *leastconn* - Selects the server with the least number of connections; it is recommended for longer sessions. Servers in the same backend are also rotated in a round-robin fashion.
      * *source* - This selects which server to use based on a hash of the source IP i.e. your user's IP address. This is one method to ensure that a user will connect to the same server.
      * Complete Info on balance algorithms can be found [HERE](http://cbonte.github.io/haproxy-dconv/configuration-1.5.html#4-balance)
    * a list of servers and ports

*  ***stats*** - If you want to enable HAProxy stats, which can be useful in determining how HAProxy is handling incoming traffic, you will want to add the configuration

In the end I added the following in the // /etc/haproxy/haproxy.cfg // to load balance on port 80 against two servers listening on port 8080; I also added my stats on port 8888

	
	#---------------------------------------------------------------------
	# round robin balancing for glusterfs
	#---------------------------------------------------------------------
	frontend glusterfs
	    bind 172.16.1.2:80
	    default_backend glusterfs-backend
	
	backend glusterfs-backend
	    balance roundrobin
	    mode http
	    option httpclose
	    option forwardfor
	    server gluster01 172.16.1.101:8080 check
	    server gluster02 172.16.1.102:8080 check
	
	listen stats
	    bind 172.16.1.2:8888
	    stats enable
	    stats uri /haproxy?stats
	    stats realm Strictly\ Private
	    stats auth chx:pulse
	#---------------------------------------------------------------------
	##
	##


You want to open up the ports you specified in the config. Basically ports 80 and 8888 (one for the *frontend* and the other for the *stats* page)

	
	[root@sandbox ~]# firewall-cmd --permanent --add-service=http --add-port=8888/tcp
	[root@sandbox ~]# firewall-cmd --reload 


Restart the service

	
	[root@sandbox ~]# systemctl restart haproxy


## Misc Notes

Misc HAProxy notes in no particular order

### Host Many Sites

You can host many sites off of one port depending on the HTTP_HOST header

	
	frontend http-in
	        bind *:80
	
	        # Define hosts
	        acl host_bacon hdr(host) -i ilovebacon.com
	        acl host_milkshakes hdr(host) -i bobsmilkshakes.com
	
	        ## figure out which one to use
	        use_backend bacon_cluster if host_bacon
	        use_backend milshake_cluster if host_milkshakes
	
	backend baconcluster
	        balance leastconn
	        option httpclose
	        option forwardfor
	        cookie JSESSIONID prefix
	        server node1 10.0.0.1:8080 cookie A check
	        server node1 10.0.0.2:8080 cookie A check
	        server node1 10.0.0.3:8080 cookie A check
	
	
	backend milshake_cluster
	        balance leastconn
	        option httpclose
	        option forwardfor
	        cookie JSESSIONID prefix
	        server node1 10.0.0.4:8080 cookie A check
	        server node1 10.0.0.5:8080 cookie A check
	        server node1 10.0.0.6:8080 cookie A check
	


After we bind to port 80, we set up two acls. The *hdr* (short for header) checks the hostname header. We also specify *-i* to make sure its case insensitive, then provide the domain name that we want to match. So now we effectively have two variables; *host_bacon* and *host_milkshakes*. Then we tell HAProxy what backend to use by checking to see if the variable is true or not.
### One Block Config


You can achieve this in one config too; but you probably want to use the above instead of this (but it works)

	
	#---------------------------------------------------------------------
	# round robin balancing for glusterfs
	#---------------------------------------------------------------------
	listen glusterfs 0.0.0.0:80
	    mode http
	    stats enable
	    stats uri /haproxy?stats
	    stats realm Strictly\ Private
	    stats auth chx:pulse
	    balance roundrobin
	    option httpclose
	    option forwardfor
	    server gluster01 172.16.1.101:8080 check
	    server gluster02 172.16.1.102:8080 check
	#---------------------------------------------------------------------


### SELinux

I had an error with SELinux saying (in // /var/log/messages //)

	
	/var/log/messages-20150107:Jan  5 18:21:28 sandbox setroubleshoot: SELinux is preventing /usr/sbin/haproxy from name_connect access on the tcp_socket . For complete SELinux messages. run sealert -l ce6d448c-0272-4563-a4f8-c8c9a940cbcb
	/var/log/messages-20150107:Jan  5 18:21:28 sandbox python: SELinux is preventing /usr/sbin/haproxy from name_connect access on the tcp_socket .
	/var/log/messages-20150107-
	/var/log/messages-20150107-*****  Plugin connect_ports (85.9 confidence) suggests   *********************
	/var/log/messages-20150107-
	/var/log/messages-20150107:If you want to allow /usr/sbin/haproxy to connect to network port 5002
	/var/log/messages-20150107-Then you need to modify the port type.
	/var/log/messages-20150107-Do
	/var/log/messages-20150107-# semanage port -a -t PORT_TYPE -p tcp 5002
	/var/log/messages-20150107-    where PORT_TYPE is one of the following: commplex_link_port_t, commplex_main_port_t, dns_port_t, dnssec_port_t, fmpro_internal_port_t, http_cache_port_t, http_port_t, kerberos_port_t, ocsp_port_t, rtp_media_port_t.
	/var/log/messages-20150107-
	/var/log/messages-20150107-*****  Plugin catchall_boolean (7.33 confidence) suggests   ******************
	/var/log/messages-20150107-
	/var/log/messages-20150107-If you want to allow nis to enabled
	--
	/var/log/messages-20150107:If you want to allow haproxy to connect any
	/var/log/messages-20150107:Then you must tell SELinux about this by enabling the 'haproxy_connect_any' boolean.
	/var/log/messages-20150107-
	/var/log/messages-20150107-Do
	/var/log/messages-20150107:setsebool -P haproxy_connect_any 1
	/var/log/messages-20150107-
	/var/log/messages-20150107-*****  Plugin catchall (1.35 confidence) suggests   **************************
	/var/log/messages-20150107-
	/var/log/messages-20150107:If you believe that haproxy should be allowed name_connect access on the  tcp_socket by default.
	/var/log/messages-20150107-Then you should report this as a bug.
	/var/log/messages-20150107-You can generate a local policy module to allow this access.
	/var/log/messages-20150107-Do
	/var/log/messages-20150107-allow this access for now by executing:
	/var/log/messages-20150107:# grep haproxy /var/log/audit/audit.log | audit2allow -M mypol
	/var/log/messages-20150107-# semodule -i mypol.pp


Just do what it says 

	
	root@sandbox ~]# setsebool -P haproxy_connect_any 1


### Command Line

Check config file 

	
	[root@sandbox ~]# haproxy -f /etc/haproxy/haproxy.cfg -c
	Configuration file is valid



### SSL

Quick SSL notes (passthrough)

	
	
	### Global settings Above ^
	## maybe do this in the global settings?
	#ssl-server-verify none
	#---------------------------------------------------------------------
	# round robin balancing for OSE Broker
	#---------------------------------------------------------------------
	frontend osebroker
	    bind 172.16.1.120:80
	    default_backend osebroker-backend
	
	frontend osebroker-ssl
	    bind 172.16.1.120:443
	    mode tcp
	    default_backend osebroker-backendssl
	
	backend osebroker-backend
	    balance roundrobin
	    mode http
	    option httpclose
	    option forwardfor
	    server broker1 172.16.1.121:80 check
	    server broker2 172.16.1.122:80 check
	    server broker3 172.16.1.123:80 check
	
	backend osebroker-backendssl
	    mode tcp
	    balance source
	    server broker1-ssl 172.16.1.121:443 check # ssl verify none
	    server broker2-ssl 172.16.1.122:443 check # ssl verify none
	    server broker3-ssl 172.16.1.123:443 check # ssl verify none
	
	listen stats
	    bind 172.16.1.120:8888
	    stats enable
	    stats uri /haproxy?stats
	    stats realm Strictly\ Private
	    stats auth chx:pulse
	#---------------------------------------------------------------------
	##
	##


### Checks

To check on a different port

	
	backend bk_myapp
	 [...]
	 option httpchk get /healthz
	 http-check expect status 200
	 server srv1 10.0.0.1:80 check port 1936
	 server srv1 10.0.0.1:80 check port 1936

