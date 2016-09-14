# Server Configuration

## RHEL 6


Things to keep in mind

*  This allows ANY node that can access the SAN server to access this LUN

*  Firewall ports open - 3260, 860 (TCP and UDP)

First you create an "empty" LVM logical volume (in this case we'll use // /dev/sdb //)

	
	root@host# pvcreate /dev/sdb
	root@host# vgcreate vg_lun /dev/sdb
	root@host# lvcreate -l +100%VG -n lv_lun vg_lun


Now, install the iscsi utilities

	
	root@host# yum -y install scsi-target-utils


Make sure the service is on and starts on boot

	
	root@host# service tgtd start
	root@host# chkconfig tgtd on


Now edit the // /etc/tgt/targets.conf // file and add the following (using the LVM LV)

	
	`<target iqn.2013-10.com.4over:chrish.sbx.target1>`
	    backing-store /dev/vg_lun/lv_lun
	`</target>`


Make sure you make your target in the following format - // **iqn.YYYY-MM.DOMAIN.TLD:HOSTNAME.TARGETNAME** //

Use // **system-config-firewall** // to open ports.

Now make sure to restart the daemon

	
	root@host# service tgtd restart


## RHEL 7

In RHEL 7 there is a new way of creating an iSCSI initiator Called *Linux-IO Target* (LIO)

*Linux-IO Target* is based on a SCSI engine that implements the semantics of a SCSI target as described in the SCSI Architecture Model (SAM), and supports its comprehensive SPC-3/SPC-4 feature set in a fabric-agnostic way. The SCSI target core does not directly communicate with initiators and it does not directly access data on disk.

First install the software and start the service:

	
	yum -y install targetcli
	systemctl enable target


[Create an LVM volume](iscsi_notes#rhel_6) just like you would normally.

Now *targetcli* can be used interactively via a shell or from the command line. We'll use the shell and do an *ls* to list current config

	
	[root@nas ~]# targetcli 
	targetcli shell version 2.1.fb31
	Copyright 2011-2013 by Datera, Inc and others.
	For help on commands, type 'help'.
	
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 0]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 0]
	  o- loopback ......................................................................................................... [Targets: 0]
	/>


You have to tell the LIO iSCSI target software about the block device you want to use

*  cd // /backstores/block // (“cd” is optional)

*  create [lun_name] /dev/[device]

*  *ls* to check results

	
	/> /backstores/block create lun0 /dev/vg1/lun0
	Created block storage object lun0 using /dev/vg1/lun0.
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 1]
	  | | o- lun0 ...................................................................... [/dev/vg1/lun0 (8.8GiB) write-thru deactivated]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 0]
	  o- loopback ......................................................................................................... [Targets: 0]
	/>


Create an iSCSI Qualified Name (IQN) record (you can specify and IQN, but it's easier to let it autogenerate it for you)

	
	/> /iscsi create
	Created target iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81.
	Created TPG 1.
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 1]
	  | | o- lun0 ...................................................................... [/dev/vg1/lun0 (8.8GiB) write-thru deactivated]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 1]
	  | o- iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81 ............................................................. [TPGs: 1]
	  |   o- tpg1 ............................................................................................... [no-gen-acls, no-auth]
	  |     o- acls .......................................................................................................... [ACLs: 0]
	  |     o- luns .......................................................................................................... [LUNs: 0]
	  |     o- portals .................................................................................................... [Portals: 0]
	  o- loopback ......................................................................................................... [Targets: 0]
	/>


Now Create an iSCSI portal address. Unless you specify an address, it will listen on all addresses

	
	/> /iscsi/iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81/tpg1/portals create 
	Using default IP port 3260
	Binding to INADDR_ANY (0.0.0.0)
	Created network portal 0.0.0.0:3260.
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 1]
	  | | o- lun0 ...................................................................... [/dev/vg1/lun0 (8.8GiB) write-thru deactivated]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 1]
	  | o- iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81 ............................................................. [TPGs: 1]
	  |   o- tpg1 ............................................................................................... [no-gen-acls, no-auth]
	  |     o- acls .......................................................................................................... [ACLs: 0]
	  |     o- luns .......................................................................................................... [LUNs: 0]
	  |     o- portals .................................................................................................... [Portals: 1]
	  |       o- 0.0.0.0:3260 ..................................................................................................... [OK]
	  o- loopback ......................................................................................................... [Targets: 0]
	/> 


Now you export this LUN

	
	/> /iscsi/iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81/tpg1/luns create /backstores/block/lun0 
	Created LUN 0.
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 1]
	  | | o- lun0 ........................................................................ [/dev/vg1/lun0 (8.8GiB) write-thru activated]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 1]
	  | o- iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81 ............................................................. [TPGs: 1]
	  |   o- tpg1 ............................................................................................... [no-gen-acls, no-auth]
	  |     o- acls .......................................................................................................... [ACLs: 0]
	  |     o- luns .......................................................................................................... [LUNs: 1]
	  |     | o- lun0 ..................................................................................... [block/lun0 (/dev/vg1/lun0)]
	  |     o- portals .................................................................................................... [Portals: 1]
	  |       o- 0.0.0.0:3260 ..................................................................................................... [OK]
	  o- loopback ......................................................................................................... [Targets: 0]
	/>


For testing purposes we'll remove ACLs. This is dangerous! Don't do this in production, you'll want to set ACLs.

	
	/> cd /iscsi/iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81/tpg1
	/iscsi/iqn.20...7efe7b81/tpg1> set attribute authentication=0
	Parameter authentication is now '0'.
	/iscsi/iqn.20...7efe7b81/tpg1> set attribute demo_mode_write_protect=0
	Parameter demo_mode_write_protect is now '0'.
	/iscsi/iqn.20...7efe7b81/tpg1> set attribute generate_node_acls=1
	Parameter generate_node_acls is now '1'.
	/iscsi/iqn.20...7efe7b81/tpg1> set attribute cache_dynamic_acls=1
	Parameter cache_dynamic_acls is now '1'.
	/iscsi/iqn.20...7efe7b81/tpg1> cd /
	/> ls
	o- / ......................................................................................................................... [...]
	  o- backstores .............................................................................................................. [...]
	  | o- block .................................................................................................. [Storage Objects: 1]
	  | | o- lun0 ........................................................................ [/dev/vg1/lun0 (8.8GiB) write-thru activated]
	  | o- fileio ................................................................................................. [Storage Objects: 0]
	  | o- pscsi .................................................................................................. [Storage Objects: 0]
	  | o- ramdisk ................................................................................................ [Storage Objects: 0]
	  o- iscsi ............................................................................................................ [Targets: 1]
	  | o- iqn.2003-01.org.linux-iscsi.nas.x8664:sn.8ab37efe7b81 ............................................................. [TPGs: 1]
	  |   o- tpg1 .................................................................................................. [gen-acls, no-auth]
	  |     o- acls .......................................................................................................... [ACLs: 0]
	  |     o- luns .......................................................................................................... [LUNs: 1]
	  |     | o- lun0 ..................................................................................... [block/lun0 (/dev/vg1/lun0)]
	  |     o- portals .................................................................................................... [Portals: 1]
	  |       o- 0.0.0.0:3260 ..................................................................................................... [OK]
	  o- loopback ......................................................................................................... [Targets: 0]
	/>


Save config by typing *saveconfig* (exit the shell as well)

	
	/> saveconfig 
	Last 10 configs saved in /etc/target/backup.
	Configuration saved to /etc/target/saveconfig.json
	/> exit
	Global pref auto_save_on_exit=true
	Last 10 configs saved in /etc/target/backup.
	Configuration saved to /etc/target/saveconfig.json


This will create the // /etc/target/saveconfig.json // file...inspect if you wish

	
	less /etc/target/saveconfig.json


Make sure *target* is started/restarted

	
	[root@nas ~]# systemctl restart target
	[root@nas ~]# systemctl status target
	target.service - Restore LIO kernel target configuration
	   Loaded: loaded (/usr/lib/systemd/system/target.service; enabled)
	   Active: active (exited) since Thu 2014-04-17 12:23:55 PDT; 1s ago
	  Process: 2010 ExecStop=/usr/bin/targetctl clear (code=exited, status=0/SUCCESS)
	  Process: 2017 ExecStart=/usr/bin/targetctl restore (code=exited, status=0/SUCCESS)
	 Main PID: 2017 (code=exited, status=0/SUCCESS)
	
	Apr 17 12:23:54 nas systemd[1]: Starting Restore LIO kernel target configuration...
	Apr 17 12:23:55 nas systemd[1]: Started Restore LIO kernel target configuration.


Set up the firewall thusly

	
	root@host# firewall-cmd --permanent --add-port=3260/tcp
	root@host# firewall-cmd --reload

# Client Configuration

You can connect to an iscsi disk with the "iscsiadm" command. First you must make sure
it's installed. Use the "yum" command...
        root@host# yum -y whatprovides *bin/iscsiadm
        Loaded plugins: fastestmirror
        Loading mirror speeds from cached hostfile
        .
        .
        .
        
            iscsi-initiator-utils-6.2.0.872-10.el5.x86_64 : iSCSI daemon and utility programs
            Repo        : base
            Matched from:
            Filename    : /sbin/iscsiadm
            
        root@host# yum -y install iscsi-initiator-utils

I found this from the man page; just search for "example"
        root@host# iscsiadm --mode discoverydb --type sendtargets --portal 192.168.122.49 --discover 
          192.168.122.49:3260,1 iqn.2011-09.com.example:for.all

Now "login" to connect to the target
        root@host# iscsiadm --mode node --targetname iqn.2011-09.com.example:for.all  --portal 192.168.122.49:3260 --login

Now you can start the iscsi daemon (make sure it starts on boot)
        root@host# /etc/init.d/iscsi start
        root@host# /etc/init.d/iscsid start
        root@host# chkconfig iscsi on
        root@host# chkconfig iscsid on

You can look in "/var/log/messages" to see what the device is called (usually /dev/sdX). You can now
partition it and mount it like a regular drive. If you want it persistant on reboots make sure you add
it to the /etc/fstab file. Note: You may want to use the UUID since the dev name will change.

You can now partition/format this device.

Mount this target, remembering using the "_netdev" mount option in the // /etc/fstab //

	
	/dev/sda /data/mount ext4 _netdev 0 0


# With Multipathing

This is assuming that you have 2 NICs to the ISCSI target

First install the // device-mapper-multipath // package.

	
	yum -y install device-mapper-multipath


Now, login to your target TWICE (once for each path/nic)

	
	iscsiadm --mode discoverydb --type sendtargets --portal 172.16.1.41 --discover 
	iscsiadm --mode discoverydb --type sendtargets --portal 172.16.2.41 --discover 
	iscsiadm --mode node --targetname iqn.2011-09.com.example:for.all  --portal 172.16.1.41:3260 --login
	iscsiadm --mode node --targetname iqn.2011-09.com.example:for.all  --portal 172.16.2.41:3260 --login


Once you did this successfully, create a multipath file (clean it up if you wish)

	
	mpathconf --enable --user_friendly_names n --find_multipaths y --with_module y --with_chkconfig y
	sed -i '/^#/d;/^$/d' /etc/multipath.conf 


Now, you need to grab a line from the example file

	
	grep getuid_callout /usr/share/doc/device-mapper-multipath-0.4.9/multipath.conf.defaults | sort -u | head -1 >> /etc/multipath.conf 


Add the // --replace-whitespace // option in the // getuid_callout // portion; and in the end...the file should look something like this...

	
	[root@mpath-client ~]# cat /etc/multipath.conf 
	defaults {
	        find_multipaths yes
	        getuid_callout "/lib/udev/scsi_id --replace-whitespace --whitelisted --device=/dev/%n"
	        user_friendly_names no
	        path_grouping_policy multibus
	        path_selector "service-time 0"
	}
	blacklist {
	}


I've also added...


*  **path_grouping_policy multibus** ~ Uses both paths at the same time (instead of just one w/ failover) 

*  **path_selector "service-time 0"** ~ The next I/O request goes to the lest busy path

Now you can find the name of your device with the multipath command (noted here as **1IET_00010001** )...

	
	[root@mpath-client ~]# multipath -ll
	1IET_00010001 dm-2 IET,VIRTUAL-DISK
	size=11G features='0' hwhandler='0' wp=rw
	|-+- policy='round-robin 0' prio=1 status=active
	| `- 4:0:0:1 sdb 8:16 active ready running
	`-+- policy='round-robin 0' prio=1 status=enabled
	  `- 3:0:0:1 sdc 8:32 active ready running


Laydown your filesystem

	
	mkfs.ext4 /dev/mapper/1IET_00010001 


Find out your UUID

	
	[root@mpath-client ~]# blkid  /dev/mapper/1IET_00010001
	/dev/mapper/1IET_00010001: UUID="2f02223d-503c-4eb6-86a8-3800c2fb9795" TYPE="ext4"


And, as with all iscsi device, make sure you put // _netdev // to mount it on boot.

	
	[root@mpath-client ~]# grep '_netdev' /etc/fstab 
	UUID=2f02223d-503c-4eb6-86a8-3800c2fb9795       /data   ext4    defaults,_netdev 0 0


More info on multipathing can be found on the [Red Hat Doc Site](https///access.redhat.com/site/documentation/en-US/Red_Hat_Enterprise_Linux)
