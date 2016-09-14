Note: Red Hat Online Documentation - https://access.redhat.com/knowledge/docs/Red_Hat_Storage/

# Create Bricks with XFS

Created "bricks" with LVM (leaving room for snapshots) and fromat with [XFS](xfs)

__Notes:__

*  Install XFS filesystem (with EPEL Repo if needed)

*  GlusterFS uses extended attributes on files, you must increase the inode size to 512 bytes (default is 256 bytes) 

	
	glusterN# pvcreate /dev/sdN
	glusterN# vgcreate gfs_brick /dev/sdN
	glusterN# lvcreate -L +nG -n lv_brick gfs_brick
	glusterN# mkfs.xfs -i size=512 /dev/gfs_brick/lv_brick 
	glusterN# mount.xfs  /dev/gfs_brick/lv_brick /exp
	glusterN# echo "/dev/gfs_brick/lv_brick	xfs	defaults	0 0" >> /etc/fstab


Make sure that the service is up and it starts on boot

	
	glusterN# service glusterd start
	glusterN# chkconfig glusterd on


# Probing Peers

Make sure each server knows about each other... (NOTE: Don't probe yourself)

	
	glusterN# gluster peer probe glusterN


# Creating Volumes

Create volume - Distributed Replicated i.e. RAID 10

`glusterN# gluster volume create export replica 2 transport tcp gluster1:/exp gluster2:/exp gluster3:/exp gluster4:/exp`

Create volume - Replicated i.e. RAID 1

`glusterN# gluster volume create demo1 replica 2 transport tcp gluster1:/demo1 gluster2:/demo1`

Start the volume

`glusterN# gluster volume start export`

Setup Fast Failover

`<code>`glusterN# gluster volume set export network.frame-timeout 60
glusterN# gluster volume set export network.ping-timeout 20`</code>`

Test Client

`<code>`[root@gluster-client ~]# time df -h 
Filesystem            Size  Used Avail Use% Mounted on
/dev/mapper/vg_glusterclient-lv_root
                      	5.5G  759M  4.5G  15% /
tmpfs                 499M     0  499M   0% /dev/shm
/dev/sda1             485M   31M  429M   7% /boot
gluster1:/export       16G   65M   16G   1% /export
	
real	0m0.003s
user	0m0.000s
sys	0m0.000s
`</code>`
SHUTDOWN gluster1 and gluster4

`<code>`[root@gluster-client ~]# time df -h 
Filesystem            Size  Used Avail Use% Mounted on
/dev/mapper/vg_glusterclient-lv_root
                      	5.5G  759M  4.5G  15% /
tmpfs                 499M     0  499M   0% /dev/shm
/dev/sda1             485M   31M  429M   7% /boot
gluster1:/export       16G   65M   16G   1% /export
	
real	0m22.867s
user	0m0.000s
sys	0m0.000s

`</code>`

# Mounting

Mouting with native client...

You can specify the following options when using the mount -t glusterfs command. Note that you need to separate all options with commas.

__//backup-volfile-servers=server-name//__ - name of the backup volfile server to mount the client. If this option is added while mounting fuse client, when the first volfile server fails, then the server specified in backup-volfile-servers option is used as volfile server to mount the client.

__//fetch-attempts=number//__ - number of attempts to fetch volume files while mounting a volume. This option is useful when you mount a server with multiple IP addresses or when round-robin DNS is configured for the server name. (NOTE: Don't copy and paste the below line...I used line-break for "easy reading"...but the example is good nonetheless)

	client# mount -t glusterfs \
		-o backup-volfile-servers=gluster3,log-level=WARNING,log-file=/var/log/gluster.log,acl \
		gluster1:/export /export  \

Common Options:


*  *log-level* ~ Logs only specified level or higher severity messages in the log-file.

*  *log-file* ~ Logs the messages in the specified file.

*  *ro* ~ Mounts the file system as read only.

*  *acl* ~ Enables POSIX Access Control List on mount.

*  *selinux* ~ Enables handling of SELinux xattrs through the mount point.

*  *background-qlen=n* ~ Enables FUSE to handle n number of requests to be queued before subsequent requests are denied. Default value of n is 64.

*  *enable-ino32* ~ this option enables file system to present 32-bit inodes instead of 64- bit inodes.

# GEO-Replication (Copy to a remote location)

First Make sure you have [NTP](ntp_notes#client_set_up) set up on both gluster systems

On the "master" and "slave" system...create SSH keys and copy them over.

	
	glusterMaster# ssh-keygen -t rsa
	glusterMaster# ssh-copy-id -i ~/.ssh/id_rsa.pub glusterN
	glusterSlave# ssh-keygen -t rsa
	glusterSlave# ssh-copy-id -i ~/.ssh/id_rsa.pub glusterN


Now (on the master) set up geo-replication

`<code>`glusterMaster# gluster volume geo-replication demo1 root@gluster2::demo2 start
glusterMaster# gluster volume geo-replication demo1 root@gluster2::demo2 status`</code>`
# Installing UFO Frontend

UFO (provided by Openstack Swift) is a web "front end" for the GlusterFS volume that provides a REST API similar to Amazon S3 Bucket Storage.

This works for RHSS version 2.0

Basic Steps are \\

1) Install Gluster Nodes \\
2) Make your LVM volume with // pvcreate --dataalignment 2560k /dev/vd? // and format bricks with // mkfs.xfs -i size=1024,maxpct=0  /path/to/vol // \\
3) Probe your peer and create your Gluster Volume (and start it!) \\
4) Start Memcache \\
5) Configure Swift \\

To configure swift, edit the // /etc/swift/proxy-server.conf // file

Under // [filter:tempauth] // add a line...

	
	user_vol_foo = bar .admin


The format of this is // user_VOLUMENAME_USERNAME = PASSWORD .admin // where: \\

*  VOLUMENAME is the name of your gluster volume

*  USERNAME is the user you are defining

*  PASSWORD is a password you pick

Under  // [filter:cache] // add the line

	
	memcache_servers = 192.168.3.171:11211,192.168.3.172:11211


The format is // `<ip>`:`<port>` // of your memcache servers (in this case they are running on the same gluster nodes) separated by a comma

Next, save this file and copy it over to all the other gluster nodes in your cluster

	
	for i in gluster2.gln.4over.com gluster3.gln.4over.com gluster4.gln.4over.com
	do
	  scp /etc/swift/proxy-server.conf $i:/etc/swift/proxy-server.conf
	done


Now start the services (On EACH gluster node)

	
	chkconfig memcached on
	chkconfig gluster-swift-proxy on 
	chkconfig gluster-swift-account on
	chkconfig gluster-swift-container on
	chkconfig gluster-swift-object on
	service memcached start
	swift-init main start


Now you can test this with the API. First you must get the token.

	
	[chrish@aardvark ~]$ curl -v -H X-Storage-User:vol:foo -H X-Storage-Pass:bar -k http://gluster1.gln.4over.com:8080/auth/v1.0

	* About to connect() to gluster1.gln.4over.com port 8080
	*   Trying 192.168.3.171... connected
	* Connected to gluster1.gln.4over.com (192.168.3.171) port 8080
	> GET /auth/v1.0 HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster1.gln.4over.com:8080
	> Accept: */*
	> X-Storage-User:vol:foo
	> X-Storage-Pass:bar
	> 
	< HTTP/1.1 200 OK
	< X-Storage-Url: http://127.0.0.1:8080/v1/AUTH_vol
	< X-Storage-Token: AUTH_tk05808c024b2743aba407ede547558dc5
	< X-Auth-Token: AUTH_tk05808c024b2743aba407ede547558dc5
	< Content-Length: 0
	< Date: Tue, 02 Jul 2013 16:27:50 GMT

	* Connection #0 to host gluster1.gln.4over.com left intact
	* Closing connection #0


Now with this token you can use "PUT" to put and "GET" to get objects.

sample get

	
	[chrish@aardvark ~]$ curl -v -X GET -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster1.gln.4over.com:8080/v1/AUTH_vol/

	* About to connect() to gluster1.gln.4over.com port 8080
	*   Trying 192.168.3.171... connected
	* Connected to gluster1.gln.4over.com (192.168.3.171) port 8080
	> GET /v1/AUTH_vol/ HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster1.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 204 No Content
	< X-Account-Container-Count: 0
	< X-Account-Object-Count: 0
	< X-Bytes-Used: 0
	< X-Object-Count: 0
	< X-Account-Bytes-Used: 0
	< X-Type: Account
	< X-Container-Count: 0
	< Accept-Ranges: bytes
	< Content-Length: 0
	< Date: Tue, 02 Jul 2013 16:30:13 GMT

	* Connection #0 to host gluster1.gln.4over.com left intact
	* Closing connection #0


Create a container

	
	[chrish@aardvark ~]$ curl -v -X PUT -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster1.gln.4over.com:8080/v1/AUTH_vol/foo

	* About to connect() to gluster1.gln.4over.com port 8080
	*   Trying 192.168.3.171... connected
	* Connected to gluster1.gln.4over.com (192.168.3.171) port 8080
	> PUT /v1/AUTH_vol/foo HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster1.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 201 Created
	< Content-Length: 18
	< Content-Type: text/html; charset=UTF-8
	< Date: Tue, 02 Jul 2013 16:31:28 GMT
	201 Created


List your containers

	
	[chrish@aardvark ~]$ curl -v -X GET -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster1.gln.4over.com:8080/v1/AUTH_vol/

	* About to connect() to gluster1.gln.4over.com port 8080
	*   Trying 192.168.3.171... connected
	* Connected to gluster1.gln.4over.com (192.168.3.171) port 8080
	> GET /v1/AUTH_vol/ HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster1.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 200 OK
	< X-Account-Container-Count: 1
	< X-Account-Object-Count: 0
	< X-Bytes-Used: 0
	< X-Object-Count: 0
	< X-Account-Bytes-Used: 0
	< X-Type: Account
	< X-Container-Count: 1
	< Accept-Ranges: bytes
	< Content-Length: 4
	< Content-Type: text/plain; charset=utf-8
	< Date: Tue, 02 Jul 2013 16:31:41 GMT
	foo

	* Connection #0 to host gluster1.gln.4over.com left intact
	* Closing connection #0


Put a file and call it "myfile.txt"

	
	[chrish@aardvark ~]$ curl -v -X PUT -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster1.gln.4over.com:8080/v1/AUTH_vol/foo/myfile.txt -T ./_4over.com_zonexfr.txt 

	* About to connect() to gluster1.gln.4over.com port 8080
	*   Trying 192.168.3.171... connected
	* Connected to gluster1.gln.4over.com (192.168.3.171) port 8080
	> PUT /v1/AUTH_vol/foo/myfile.txt HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster1.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> Content-Length: 67269
	> Expect: 100-continue
	> 
	< HTTP/1.1 100 Continue
	< HTTP/1.1 201 Created
	< Content-Length: 118
	< Content-Type: text/html; charset=UTF-8
	< Etag: ba5f37b8b01fd10e93ad41d37c2cf08e
	< Last-Modified: Tue, 02 Jul 2013 16:33:12 GMT
	< Date: Tue, 02 Jul 2013 16:33:13 GMT
	`<html>`
	 `<head>`
	  `<title>`201 Created`</title>`
	 `</head>`
	 `<body>`
	  `<h1>`201 Created`</h1>`
	  `<br />``<br />`
	
	
	
	 `</body>`
	Connection #0 to host gluster1.gln.4over.com left intact

	* Closing connection #0


Now see that file

	
	[chrish@aardvark ~]$ curl -v -X GET -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster2.gln.4over.com:8080/v1/AUTH_vol/foo/

	* About to connect() to gluster2.gln.4over.com port 8080
	*   Trying 192.168.3.172... connected
	* Connected to gluster2.gln.4over.com (192.168.3.172) port 8080
	> GET /v1/AUTH_vol/foo/ HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster2.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 200 OK
	< X-Container-Object-Count: 1
	< X-Container-Bytes-Used: 67269
	< Accept-Ranges: bytes
	< Content-Length: 11
	< Content-Type: text/plain; charset=utf-8
	< Date: Tue, 02 Jul 2013 16:34:45 GMT
	myfile.txt

	* Connection #0 to host gluster2.gln.4over.com left intact
	* Closing connection #0
	[chrish@aardvark ~]$ curl -v -X GET -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://gluster2.gln.4over.com:8080/v1/AUTH_vol/foo/

	* About to connect() to gluster2.gln.4over.com port 8080
	*   Trying 192.168.3.172... connected
	* Connected to gluster2.gln.4over.com (192.168.3.172) port 8080
	> GET /v1/AUTH_vol/foo/ HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: gluster2.gln.4over.com:8080
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 200 OK
	< X-Container-Object-Count: 1
	< X-Container-Bytes-Used: 67269
	< Accept-Ranges: bytes
	< Content-Length: 11
	< Content-Type: text/plain; charset=utf-8
	< Date: Tue, 02 Jul 2013 17:21:07 GMT
	myfile.txt

	* Connection #0 to host gluster2.gln.4over.com left intact
	* Closing connection #0


Delete objects/containers

	
	[chrish@aardvark ~]$ curl -v -X DELETE -H X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5 -k http://ufo.dev.4over.com/v1/AUTH_vol/foo

	* About to connect() to ufo.dev.4over.com port 80
	*   Trying 192.168.11.210... connected
	* Connected to ufo.dev.4over.com (192.168.11.210) port 80
	> DELETE /v1/AUTH_vol/foo HTTP/1.1
	> User-Agent: curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5
	> Host: ufo.dev.4over.com
	> Accept: */*
	> X-Auth-Token:AUTH_tk05808c024b2743aba407ede547558dc5
	> 
	< HTTP/1.1 204 No Content
	< Content-Length: 0
	< Content-Type: text/html; charset=UTF-8
	< Date: Tue, 02 Jul 2013 19:09:06 GMT

	* Connection #0 to host ufo.dev.4over.com left intact
	* Closing connection #0


## NOTES FOR Version 2.1

Version 2.1 I had to...

1) Copy the example files to "real" files

	
	cd /etc/swift
	cp object-server.conf-gluster object-server.conf
	cp container-server.conf-gluster container-server.conf
	cp account-server.conf-gluster account-server.conf
	cp swift.conf-gluster swift.conf
	cp proxy-server.conf-gluster proxy-server.conf
	cp fs.conf-gluster fs.conf
	cp object-expirer.conf-gluster object-expirer.conf


2) Add this to the // /etc/swift/proxy-server.conf // file

	
	[filter:tempauth]
	use = egg:swift#tempauth
	user_4api_teapot = quadrapus .admin
	user_hermes_hermes = hermes .admin


3) Make sure the // [pipeline:main] // line in the // /etc/swift/proxy-server.conf //  file includes the **tempauth** entry. It should look something like this

	
	[pipeline:main]
	pipeline = catch_errors healthcheck proxy-logging cache proxy-logging tempauth proxy-server


4) You need to create a "ring" config (whatever that is) (note that you specify the Gluster Vol Name and all in one line)

	
	gluster-swift-gen-builders 4api hermes


5) In 2.1 the services have different names

	
	chkconfig memcached on
	chkconfig openstack-swift-proxy on
	chkconfig openstack-swift-account on
	chkconfig openstack-swift-container on
	chkconfig openstack-swift-object on
	chkconfig openstack-swift-object-expirer on


NOTE: I mounted the UFO brick with the following options...

	
	[root@gluster1.la3.4over.com ~]# grep xfs /etc/fstab 
	UUID="f5eb892e-255a-4acb-83fa-e99bf3b958a3"     /netapp_brick   xfs     rw,noatime,nodiratime,inode64   0 0


# CTDB Setup

In replicated volume environment, you can configure Cluster Trivial Database (CTDB) to provide high availability for NFS and SMB exports. CTDB adds virtual IP addresses (VIPs) and a heartbeat service to Red Hat Storage Server.

What you'll need


*  A glusterfs volume

*  Floating IP address

*  Patience

First create a volume for CTDB

	
	gluster volume create ctdb replica 2 gluster1.gln.4over.com:/brick/ctdb gluster2.gln.4over.com:/brick/ctdb


Now, Update the ***META=all*** to the newly created volume name on all Red Hat Storage servers which require IP failover in the hook scripts available at // **/var/lib/glusterd/hooks/1/start/post/S29CTDBsetup.sh** // and **// /var/lib/glusterd/hooks/1/stop/pre/S29CTDB-teardown.sh //** (MAKE SURE YOU RUN THESE COMMANDS ON ALL NODES IN YOUR CLUSTER)

	
	cp /var/lib/glusterd/hooks/1/start/post/S29CTDBsetup.sh ~/S29CTDBsetup.sh.bak
	cp /var/lib/glusterd/hooks/1/stop/pre/S29CTDB-teardown.sh ~/S29CTDB-teardown.sh.bak
	sed -i 's/^META\=\"all\"/META\=\"ctdb\"/g' /var/lib/glusterd/hooks/1/start/post/S29CTDBsetup.sh
	sed -i 's/^META\=\"all\"/META\=\"ctdb\"/g' /var/lib/glusterd/hooks/1/stop/pre/S29CTDB-teardown.sh


Now that those scripts are in place, you can start the volume (there should be a *** /gluster/lock *** volume created)

	
	gluster  volume  start ctdb
	df -h /gluster/lock
	   Filesystem            Size  Used Avail Use% Mounted on
	   gluster1.gln.4over.com:ctdb
	                          10G  5.0G  5.0G  50% /gluster/lock


Create ***/gluster/lock/ctdb*** file and add the following entries

	
	CTDB_RECOVERY_LOCK=/gluster/lock/lockfile
	CTDB_PUBLIC_ADDRESSES=/etc/ctdb/public_addresses
	CTDB_MANAGES_SAMBA=yes
	CTDB_NODES=/etc/ctdb/nodes


Create **// /gluster/lock/nodes //** file and list the IPs of Red Hat Storage servers that are in your cluster.

	
	192.168.3.171
	192.168.3.172


Create ***/gluster/lock/public_addresses *** file and list the Virtual IPs that CTDB should create. Replace eth0 with the interface available on that node for CTDB to use. This should be your floating IP

	
	192.168.3.143/24 eth0


Run the following commands on all Red Hat Storage servers which require IP failover to create symbolic links

	
	mv /etc/sysconfig/ctdb /etc/sysconfig/ctdb.bak
	ln -s /gluster/lock/ctdb /etc/sysconfig/ctdb
	ln -s /gluster/lock/nodes /etc/ctdb/nodes
	ln -s /gluster/lock/public_addresses /etc/ctdb/public_addresses 


In the end your nodes should look like this...

	
	[root@glusterX ~]# ll /etc/ctdb | egrep 'public_address|nodes'
	lrwxrwxrwx  1 root root    19 Feb 18 23:53 nodes -> /gluster/lock/nodes
	lrwxrwxrwx  1 root root    30 Feb 18 23:53 public_addresses -> /gluster/lock/public_addresses
	
	[root@glusterX ~]# ll /etc/sysconfig/ctdb
	lrwxrwxrwx 1 root root 18 Feb 18 23:53 /etc/sysconfig/ctdb -> /gluster/lock/ctdb
	
	[root@glusterX ~]# ll /gluster/lock/
	total 2
	-rw-r--r-- 1 root root 141 Feb 18 23:21 ctdb
	-rw-r--r-- 1 root root  28 Feb 18 23:38 nodes
	-rw-r--r-- 1 root root  22 Feb 18 23:48 public_addresses


NOW, start the service on each node

	
	service ctdb start
	chkconfig ctdb on


Make sure that samba is off on each node

	
	chkconfig smb off


Verify that CTDB is running using the following command(s):

	
	ctdb status
	ctdb ip
	ctdb ping -n all
	ctdb listnodes

# Misc Commands

__Grow XFS filesystem (with LVM)__

	
		glusterN# df -hF xfs
			Filesystem            Size  Used Avail Use% Mounted on
			/dev/mapper/wrkflow_brick-lv_workflow
	                      		5.0G   33M  5.0G   1% /bricks/workflow_storehost
		glusterN# vgs
	  		VG            #PV #LV #SN Attr   VSize  VFree
	  		VolGroup        1   2   0 wz--n- 31.51g    0 
	  		wrkflow_brick   5   1   0 wz--n-  9.99t 9.99t
		glusterN# lvs
	  		LV          VG            Attr     LSize  Pool Origin Data%  Move Log Copy%  Convert
	  		lv_root     VolGroup      -wi-ao-- 25.63g                                           
	  		lv_swap     VolGroup      -wi-ao--  5.88g                                           
	  		lv_workflow wrkflow_brick -wi-ao--  5.00g
		glusterN# # lvextend -L 5T wrkflow_brick/lv_workflow
	  		Extending logical volume lv_workflow to 5.00 TiB
	  		Logical volume lv_workflow successfully resize
		glusterN# xfs_growfs /bricks/workflow_storehost
			meta-data=/dev/mapper/wrkflow_brick-lv_workflow isize=512    agcount=4, agsize=327680 blks
	         	 	 =                       sectsz=512   attr=2
			data     =                       bsize=4096   blocks=1310720, imaxpct=25
	         		 =                       sunit=0      swidth=0 blks
			naming   =version 2              bsize=4096   ascii-ci=0
			log      =internal               bsize=4096   blocks=2560, version=2
	         		 =                       sectsz=512   sunit=0 blks, lazy-count=1
			realtime =none                   extsz=4096   blocks=0, rtextents=0
			data blocks changed from 1310720 to 1342177280
		glusterN# df -hF xfs
			Filesystem            Size  Used Avail Use% Mounted on
			/dev/mapper/wrkflow_brick-lv_workflow
	                      		5.0T  161M  5.0T   1% /bricks/workflow_storehost


When you get an error when trying to re-use bricks that you have removed from a volume that you're no loger using...

	
	{path} or a prefix of it is already part of a volume


You can fix this by re-setting attributes 

For the directory (or any parent directories) that was formerly part of a volume, simply (ON EVERY PEER):

	
	root@gluserN# brick_path=/path/to/brick
	root@gluserN# setfattr -x trusted.glusterfs.volume-id $brick_path
	root@gluserN# setfattr -x trusted.gfid $brick_path
	root@gluserN# rm -rf $brick_path/.glusterfs


You should be able to add the bricks to your new volume now

__Fuse Module__

I've encountered an error where the volume won't mount. The log gives a "fuse not found error"

	
	[2013-01-04 15:27:15.232943] E [mount.c:596:gf_fuse_mount] 0-glusterfs-fuse: cannot open /dev/fuse (No such file or directory)


You must make sure the fuse module is loaded.

	
	root@host# modprobe fuse
	root@host# mount /share/store
	root@host# df -h | grep gluster
	glusterfs#gluster1.la3.4over.com:/workflow2

