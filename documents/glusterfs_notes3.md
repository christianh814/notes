# GlusterFS Notes

These are glusterfs notes in no paticular order. Old 3.1 notes can be found [here](glusterfs_notes.md)

* [Installation](#installation)
* [Building a Trusted Storage Pool](#building-a-trusted-storage-pool)
* [Creating Bricks](#creating-bricks)
* [Building a Volume](#building-a-volume)
* [Client](#client)
* [NFS](#nfs)
* [CIFS](#cifs)
* [Volume Options](#volume-options)
* [ACLs And Quotas](#acls-and-quotas)
* [Resizing Volumes](#resizing-volumes)
* [Configuring IP Failover](#configuring-ip-failover)
* [Configuring NFS Ganesha](#configuring-nfs-ganesha)
* [Georeplication](#georeplication)
* [Basic Troubleshooting](#basic-troubleshooting)

## Installation

 If Red Hat Gluster Storage is going to be layered on top of Red Hat Enterprise Linux 7, then the following repositories need to be attached and enabled on the system:

* `rhel-7-server-rpms`
* `rh-gluster-3-for-rhel-7-server-rpms`
* `rh-gluster-3-nfs-for-rhel-7-server-rpms`
* `rh-gluster-3-samba-for-rhel-7-server-rpms`

You can also install it from the ISO from [access.redhat.com](https://access.redhat.com)

First make sure you enable the above repos. Register the system

```
subscription-manager register
```

List the available subscriptions:

```
subscription-manager list --available
```

In the output for the previous command, find the pool ID for an OpenShift Enterprise subscription and attach it:

```
subscription-manager attach --pool=${pool_id}
```

Disable all repositories and enable only the required ones:

```
subscription-manager repos  --disable=*
subscription-manager repos \
--enable=rhel-7-server-rpms \
--enable=rh-gluster-3-for-rhel-7-server-rpms \
--enable=rh-gluster-3-nfs-for-rhel-7-server-rpms \
--enable=rh-gluster-3-samba-for-rhel-7-server-rpms
```


Install with

```
yum -y install redhat-storage-server
```

Start the service

```
systemctl start --now glusterd
```
## Building a Trusted Storage Pool

To build a storage pool (i.e. a storage cluster); you need to probe other servers. First, check if it's active on both servers

```
[root@servera ~]# systemctl is-active glusterd
active

[root@serverb ~]# systemctl is-active glusterd
active
```

Next, configure the firewall

```
[root@servera ~]# firewall-cmd --permanent --add-service=glusterfs
success
[root@servera ~]# firewall-cmd --reload
success

[root@serverb ~]# firewall-cmd --permanent --add-service=glusterfs
success
[root@serverb ~]# firewall-cmd --reload
success
```

Now you can probe other nodes (but don't probe yourself)

```
[root@servera ~]# gluster peer probe serverb.lab.example.com
peer probe: success.
```

Verify

```
[root@servera ~]# gluster peer probe serverb.lab.example.com
peer probe: success.
[root@servera ~]# gluster peer status
Number of Peers: 1

Hostname: serverb.lab.example.com
Uuid: 94d5ce6f-02fe-471a-b887-20768212b7c5
State: Peer in Cluster (Connected)

[root@servera ~]# gluster pool list
UUID					Hostname               	State
94d5ce6f-02fe-471a-b887-20768212b7c5	serverb.lab.example.com	Connected
0b3fb2fc-6d99-48e7-a218-92a7009073cd	localhost              	Connected
```
## Creating Bricks

Bricks are usally created with an `lvm` thinpool. After that, you format it with `xfs` with an inode size of `512` for metadata.

This implicitly needs to be done an *ALL* servers in the pool (substituting brick names where it makes sense)

First, check to see if you have a volume group (if not create one)
```
[root@servera ~]# vgs
  VG        #PV #LV #SN Attr   VSize  VFree
  rhgs        1   2   0 wz--n-  9.51g 40.00m
  vg_bricks   1   0   0 wz--n- 20.00g 20.00g
```

In our example we will create a 10 GiB LVM thin pool called `thinpool`, inside the `vg_bricks` volume group. The `-T` means to create this as a thin provisioned volume.
```
[root@servera ~]# lvcreate -L +10G -T vg_bricks/thinpool
  Logical volume "thinpool" created.
```

Now, we will create a logical volume with a virtual size of 2 GiB called `brick-a1`, using the `vg_bricks/thinpool` LVM thin pool. The `-V` means I'm creating a "virtual" size from this thin pool (`-T`)

```
[root@servera ~]# lvcreate -V +2G -T vg_bricks/thinpool -n brick-a1
  Logical volume "brick-a1" created.

```

(Again, on the other server I'd call this `brick-b1` for example. Needs to be different)

Now, we will create an `xfs` filesystem with `512` inode size.
```
[root@servera ~]# mkfs.xfs -i size=512 /dev/vg_bricks/brick-a1
meta-data=/dev/vg_bricks/brick-a1 isize=512    agcount=8, agsize=65520 blks
         =                       sectsz=512   attr=2, projid32bit=1
         =                       crc=0        finobt=0
data     =                       bsize=4096   blocks=524160, imaxpct=25
         =                       sunit=16     swidth=16 blks
naming   =version 2              bsize=4096   ascii-ci=0 ftype=0
log      =internal log           bsize=4096   blocks=2560, version=2
         =                       sectsz=512   sunit=16 blks, lazy-count=1
realtime =none                   extsz=4096   blocks=0, rtextents=0
```

Now create mountpoints for them

```
[root@servera ~]# mkdir -p /bricks/brick-a1
```

Add them to the fstab
```
[root@servera ~]# grep brick-a1 /etc/fstab
/dev/vg_bricks/brick-a1i	/bricks/brick-a1	xfs	defaults 1 2
```

Now create a `brick` subdirectory on the new brick for your files.
```
[root@servera ~]# mkdir /bricks/brick-a1/brick
```

You will need to configure the default SELinux label for the brick directory you just created to be `glusterd_brick_t`, and then recusively relabel the directory
```
[root@servera ~]# semanage fcontext -a -t glusterd_brick_t /bricks/brick-a1/brick
[root@servera ~]# restorecon -vR /bricks/
restorecon reset /bricks/brick-a1 context system_u:object_r:unlabeled_t:s0->system_u:object_r:default_t:s0
restorecon reset /bricks/brick-a1/brick context unconfined_u:object_r:unlabeled_t:s0->unconfined_u:object_r:glusterd_brick_t:s0
```

## Building a Volume

There are different kinds of volumes you can create, we'll go over some (but not all). Volumes are built from multiple bricks found on each server. You also run these command on one server only (since you have a trusted pool set up already)

* [Distributed](#distributed)
* [Replicated](#replicated)
* [Distributed Replicated](#distributed-replicated)

### Distributed

Distributed volumes "stripe" data accross all bricks (no fault tolerance)

To create a distributed volume
```
[root@servera ~]# gluster volume create firstvol servera:/bricks/brick-a1/brick serverb:/bricks/brick-b1/brick
```

Now start the volume
```
[root@servera ~]# gluster volume start firstvol
volume start: firstvol: success
```

Verify the information
```
[root@servera ~]# gluster volume info firstvol

Volume Name: firstvol
Type: Distribute
Volume ID: aa4082e6-474d-480c-bc77-33a6f6bfc85a
Status: Started
Number of Bricks: 2
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a1/brick
Brick2: serverb:/bricks/brick-b1/brick
Options Reconfigured:
performance.readdir-ahead: on
```

Install the client on another server to test it (`glusterfs-fuse`)
```
[root@workstation ~]# yum -y install glusterfs-fuse
Loaded plugins: langpacks, search-disabled-repos
Resolving Dependencies
--> Running transaction check
---> Package glusterfs-fuse.x8664 0:3.7.1-16.el7 will be installed
--> Finished Dependency Resolution

Dependencies Resolved

=============================================================================================================================================================================
 Package                                      Arch                                 Version                                      Repository                              Size
=============================================================================================================================================================================
Installing:
 glusterfs-fuse                               x86_64                               3.7.1-16.el7                                 rhel_dvd                               107 k

Transaction Summary
=============================================================================================================================================================================
Install  1 Package

Total download size: 107 k
Installed size: 313 k
Downloading packages:
glusterfs-fuse-3.7.1-16.el7.x86_64.rpm                                                                                                                | 107 kB  00:00:00
Running transaction check
Running transaction test
Transaction test succeeded
Running transaction
  Installing : glusterfs-fuse-3.7.1-16.el7.x86_64                                                                                                                        1/1
  Verifying  : glusterfs-fuse-3.7.1-16.el7.x86_64                                                                                                                        1/1

Installed:
  glusterfs-fuse.x86_64 0:3.7.1-16.el7

Complete!
```

Test the mount now
```
[root@workstation ~]# mount.glusterfs servera:/firstvol /mnt
[root@workstation ~]# df -h
Filesystem         Size  Used Avail Use% Mounted on
/dev/vda1           10G  3.1G  7.0G  31% /
devtmpfs           902M     0  902M   0% /dev
tmpfs              920M  152K  920M   1% /dev/shm
tmpfs              920M   17M  904M   2% /run
tmpfs              920M     0  920M   0% /sys/fs/cgroup
tmpfs              184M   32K  184M   1% /run/user/1000
servera:/firstvol  4.0G   66M  4.0G   2% /mnt
```
### Replicated

Replicated volumes has your data on X number of bricks. In short, this "mirrors" your data (synchronisly). This is has fault tolerance becuase if one server were to go away, you'll still have all your data. It's like RAID1 for servers.

To create a replicated volume

```
[root@servera ~]# gluster volume create repvol replica 2 servera:/bricks/brick-a1/brick serverb:/bricks/brick-b1/brick
volume create: repvol: success: please start the volume to access data
```

Verify with
```
[root@servera ~]# gluster volume info
 
Volume Name: repvol
Type: Replicate
Volume ID: 48e6f131-92f7-4362-ace3-86325c59d672
Status: Created
Number of Bricks: 1 x 2 = 2
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a1/brick
Brick2: serverb:/bricks/brick-b1/brick
Options Reconfigured:
performance.readdir-ahead: on
```

### Distributed Replicated

Volume types can be combined if desired. This allows the creation of the volume types `distributed-replicated`. In these types, a number of either replicated or dispersed sets are formed, and the files are distributed across these sets. You can think of this as RAID10 (or RAID 0+1)

To create a distributed replicae volume

```
[root@servera ~]# gluster volume create distreplvol replica 2 servera:/bricks/brick-a3/brick serverb:/bricks/brick-b3/brick serverc:/bricks/brick-c3/brick serverd:/bricks/brick-d3/brick
volume create: distreplvol: success: please start the volume to access data

[root@servera ~]# gluster volume info distreplvol

Volume Name: distreplvol
Type: Distributed-Replicate
Volume ID: 695e26b1-b463-4e94-9297-f929270eeff7
Status: Created
Number of Bricks: 2 x 2 = 4
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a3/brick
Brick2: serverb:/bricks/brick-b3/brick
Brick3: serverc:/bricks/brick-c3/brick
Brick4: serverd:/bricks/brick-d3/brick
Options Reconfigured:

```

Another option is to disperse the data (break the files into chunks)

```
gluster volume create distdispvol disperse-data 4 redundancy 2 \
servera:/bricks/brick-a4/brick \
serverb:/bricks/brick-b4/brick \
serverc:/bricks/brick-c4/brick \
serverd:/bricks/brick-d4/brick \
servera:/bricks/brick-a5/brick \
serverb:/bricks/brick-b5/brick \
serverc:/bricks/brick-c5/brick \
serverd:/bricks/brick-d5/brick \
servera:/bricks/brick-a6/brick \
serverb:/bricks/brick-b6/brick \
serverc:/bricks/brick-c6/brick \
serverd:/bricks/brick-d6/brick froce
```
## Client

The native client is installed by default on Red Hat Gluster Storage servers. On Red Hat Enterprise Linux 7 systems the client is available in the `glusterfs-fuse` package in the `rhel-x86_64-server-7-rh-gluster-3-client` channel, and in the `rhel-x86_64-server-rh-common-7` channel, although the latter may be deprecated in the future.

Red Hat Enterprise Linux 6 clients can install the client using the `glusterfs-fuse` package in the `rhel-x86_64-server-rhsclient-6` channel. Red Hat Enterprise Linux 5 clients should use the `rhel-x86_64-server-rhsclient-5` channel.

After subscribing, install via

```
yum -y install glusterfs-fuse
```

You can easily mount it with `mount.glusterfs`

```
mount.glusterfs servera:custdata /mnt
```

Use  `backup-volfile-servers` to connect to other servers. Use `acl` for POSIX access control

```
mount.glusterfs -o backup-volfile-servers=servera:serverb:serverc:serverd,acl servera:custdata /mnt
```

Using `_netdev` in the `/etc/fstab` is a smart move

```
servera:custdata	/mnt	glusterfs	_netdev,backup-volfile-servers=servera:serverb:serverc:serverd,acl 0 0
```

Timeout is very slow by default, so you may want to shorten the timeout on the gluster server

```
root@servera# gluster volume set custdata network.frame-timeout 60
root@servera# gluster volume set custdata network.ping-timeout 20
```

## NFS

By default, any new Red Hat Gluster Storage volume will be exported over NFSv3, with ACLs enabled. This is for clients that can not run the native client.

NFSv3 exports do not use the NFSv3 server in the Linux kernel. Instead they use a dedicated NFSv3 server written in gluster, that only exports over TCP.

Unlike the native client, clients using NFSv3 will not automatically fail over during a failure.


You can see the built in NFS export with the `volume status` command

```
[root@servera ~]# gluster volume status mediadata
Status of volume: mediadata
Gluster process                             TCP Port  RDMA Port  Online  Pid
------------------------------------------------------------------------------
Brick servera:/bricks/brick-a1/brick        49159     0          Y       3038
Brick serverb:/bricks/brick-b1/brick        49159     0          Y       1917
Brick serverc:/bricks/brick-c1/brick        49159     0          Y       1926
Brick serverd:/bricks/brick-d1/brick        49159     0          Y       1900
NFS Server on localhost                     2049      0          Y       3018
NFS Server on serverc.lab.example.com       2049      0          Y       1935
NFS Server on serverd.lab.example.com       2049      0          Y       1879
NFS Server on serverb.lab.example.com       2049      0          Y       1935

Task Status of Volume mediadata
------------------------------------------------------------------------------
There are no active volume tasks
```

In order to mount this from the client, you need to open `nfs` and (since it's NFSv3) `rpc-bind` on the firewall

```
[root@servera ~]# firewall-cmd --permanent --add-service=nfs --add-service=rpc-bind
success
[root@servera ~]# firewall-cmd --reload
success
```


On the client, mount with v3 options
```
[root@workstation ~]# mount.nfs -o vers=3,proto=tcp servera:/mediadata /mnt/mediadata/
[root@workstation ~]# df -hF nfs
Filesystem          Size  Used Avail Use% Mounted on
servera:/mediadata  8.0G  130M  7.9G   2% /mnt/mediadata
```

## CIFS

You can use gluster to serve CIFS shares (Windows); using glusterfs as the backend storage system.

First, add the appropriate firewall rules

```
[root@servera ~]# firewall-cmd --permanent --add-service=samba
success
[root@servera ~]# firewall-cmd --reload
success
```

Next, install `samba`

```
[root@servera ~]# yum -y install samba
```

Next, start and enable the service

```
[root@servera ~]# systemctl enable smb.service
Created symlink from /etc/systemd/system/multi-user.target.wants/smb.service to /usr/lib/systemd/system/smb.service.

[root@servera ~]# systemctl start smb.service
```

Create a user and map it back to samba

```
[root@servera ~]# adduser smbuser

[root@servera ~]# smbpasswd -a smbuser
New SMB password:
Retype new SMB password:
Added user smbuser.
```

Now export a volume over smb; first you need to set some volume options

```
[root@servera ~]# gluster volume set mediadata stat-prefetch off
volume set: success

[root@servera ~]# gluster volume set mediadata server.allow-insecure on
volume set: success

[root@servera ~]# gluster volume set mediadata storage.batch-fsync-delay-usec 0
volume set: success
```

Configure `glusterd`  to allow Samba to communicate with bricks using insecure ports by editing the `/etc/glusterfs/glusterd.vol` file to include `option rpc-auth-allow-insecure on`

```
[root@servera ~]# cat /etc/glusterfs/glusterd.vol
volume management
    type mgmt/glusterd
    option working-directory /var/lib/glusterd
    option transport-type socket,rdma
    option transport.socket.keepalive-time 10
    option transport.socket.keepalive-interval 2
    option transport.socket.read-fail-log off
    option ping-timeout 0
    option event-threads 1
    ### ADDING THIS LINE
    option rpc-auth-allow-insecure on
    ###
#   option base-port 49152
end-volume
```

Restart `glusterd`

```
[root@servera ~]# systemctl restart glusterd
```

Now restart the volume to export it as a samba share

```
[root@servera ~]# gluster volume stop mediadata
Stopping volume will make its data inaccessible. Do you want to continue? (y/n)
volume stop: mediadata: success

[root@servera ~]#  gluster volume start mediadata
volume start: mediadata: success
```

Now, on the client, test to see if you can access the share via samba

```
[student@workstation ~]$ smbclient -L servera -U smbuser%redhat
Domain=[MYGROUP] OS=[Windows 6.1] Server=[Samba 4.2.4]

	Sharename       Type      Comment
	---------       ----      -------
	gluster-mediadata Disk      For samba share of volume mediadata
	IPC$            IPC       IPC Service (Samba Server Version 4.2.4)
Domain=[MYGROUP] OS=[Windows 6.1] Server=[Samba 4.2.4]

	Server               Comment
	---------            -------

	Workgroup            Master
	---------            -------
```

If you see that, you should be able to mount the share

Create a dir for the share

```
[student@workstation ~]$ sudo mkdir /mnt/smbdata
```

Add it to `/etc/fstab`

```
[student@workstation ~]$ grep cifs /etc/fstab
//servera/gluster-mediadata	/mnt/smbdata	cifs	user=smbuser,pass=redhat	0 0
```

Mount the share

```
[student@workstation ~]$ sudo mount /mnt/smbdata/

[student@workstation ~]$ df -h /mnt/smbdata/
Filesystem                   Size  Used Avail Use% Mounted on
//servera/gluster-mediadata  8.0G  131M  7.9G   2% /mnt/smbdata
```
## Volume Options

You can set specific options to volumes to change the behavior. The list is extensive. You can get a list of them by running...

```
[root@servera ~]# gluster volume set help
```

To set an option; use the `gluster volume set <VOLUME NAME> <OPTION> <VALUE>` syntax

```
[root@servera ~]# gluster volume set mediadata features.read-only on
```

You can reset the option to it's default value with `gluster volume reset <VOLUME NAME> <OPTION>`

```
[root@servera ~]# gluster volume reset mediadata features.read-only
```

To show all options are set on a volume

```
[root@servera ~]# gluster volume get mediadata all
```

To show specific option status

```
[root@servera ~]# gluster volume get mediadata features.read-only
```

To show non-default options are set on a volume

```
[root@servera ~]# gluster volume info mediadata
```
## ACLs And Quotas

You can set access control lists and Quotas with gluster.

* [ACLs](#acls)
* [Quotas](#quotas)

### ACLs

You can enable ACLs by adding `acl` as a mount option. ACLs is handled by the `setfacl` and `getfacl`

```
mount.glusterfs -o backup-volfile-servers=servera:serverb:serverc:serverd,acl servera:custdata /mnt
```

You can also put this in the `/etc/fstab` as well.

Now you should be able to assign ACLs

```
setfacl -m u:lisa:rwX /mnt/springfield-library

```

To make it "default"...

```
setfacl -m d:g:admins:rX /mnt/springfield-library
```

### Quotas

Gluster supports setting quotas on individual directories in a volume. This allows an admin to control where storage is being used. Unlike the file system quotas, this is done on a directory basis, and not on a user/group basis.

To enable quotas on a volume

```
gluster volume quota <VOLUME> enable
```

Note that quotas are not set for the entire volume by default, but for individual directories in that volume. To limit the entire volume, a limit can be placed on `/` on that volume. Also, after enabling a quota, clients using the native client must remount the volume.

The dir must exists first; so make sure it's there or create it

```
user@client$ mkdir /mnt/graphics/raw
user@client$ sudo umount /mnt/graphics
```

On one of the gluster servers; set up the volume to be quota-fied

```
[root@servera ~]# gluster volume quota graphics enable
volume quota : success
```

Now you can set a hard limit of `1GB` and the soft limit (the warning limit) to `50%` by using the `limit-usage /path` option (where `/path` is relative to the volume, not where it will be mounted)

```
[root@servera ~]# gluster volume quota graphics limit-usage /raw 1GB 50%
volume quota : success
```

So above; the limit will be set on `/mnt/graphics/raw` on the client.

Now you can set the `soft-timeout` (number of seconds to wait until after the quota has been reached to throw an error)

```
[root@servera ~]# gluster volume quota graphics soft-timeout 5s
volume quota : success
```

You can do the same for `hard-timeout`

```
[root@servera ~]# gluster volume quota graphics hard-timeout 1s
volume quota : success
```

Configure the graphics volume in such a way that the `df` command reports the amount of space left using quotas, and not the physical available space.

```
[root@servera ~]# gluster volume set graphics quota-deem-statfs on
volume set: success
```

You can remount and test the volume

```
[root@workstation ~]# umount /mnt/graphics
[root@workstation ~]# mount /mnt/graphics
[root@workstation ~]# dd if=/dev/zero of=/mnt/graphics/raw/testfile bs=1M
dd: error writing ’/mnt/graphics/raw/testfile’: Disk quota exceeded
dd: closing output file ’/mnt/graphics/raw/testfile’: Disk quota exceeded
```
## Resizing Volumes

A volume can be extended, without causing any downtime. You do this by adding bricks to it. Enough bricks must be added to match the current layout. Example, to extend a volume with `replica` set to two, two bricks must be added. Ergo, when extending a 4x2 `distributed-replicated` volume, at a minimum; two bricks or a multiple of two must be added. 

To extend a replicated volume (note, this turns my replicate to a distributed replicate)

```
[root@servera ~]# gluster volume add-brick myvol serverc:/bricks/brick-c1/brick serverd:/bricks/brick-d1/brick 
volume add-brick: success
```

Next you'd want to rebalance the data

```
[root@servera ~]# gluster volume rebalance  myvol start
volume rebalance: myvol: success: Rebalance on myvol has been started successfully. Use rebalance status command to check status of the rebalance process.
ID: db170509-c9dd-4f4c-bc7a-9c19408fd6e6
```

You can check the status as well

```
[root@servera ~]# gluster volume rebalance  myvol status
                                    Node Rebalanced-files          size       scanned      failures       skipped               status   run time in secs
                               ---------      -----------   -----------   -----------   -----------   -----------         ------------     --------------
                               localhost               49      980Bytes           100             0             0            completed               3.00
                 serverb.lab.example.com                0        0Bytes             0             0             0            completed               1.00
                 serverc.lab.example.com                0        0Bytes             0             0             0            completed               1.00
                 serverd.lab.example.com                0        0Bytes             0             0             0            completed               1.00
volume rebalance: myvol: success

```

Volumes can be shrunk by removing bricks. The replica count can be adjusted during removal as well.

When shrinking `distributed-replicated` volumes, the number of bricks being removed must be a multiple of the replica count. I.E. to shrink a distributed-replicated volume with a replica count of two, bricks need to be remove in multiples of two.


To shrink a volume, remove the bricks (making sure you maintain the layout or you get dataloss)

```
[root@servera ~]# gluster volume info myvol

Volume Name: myvol
Type: Distributed-Replicate
Volume ID: e6916f89-9e1d-4283-8549-5caa16e695cb
Status: Started
Number of Bricks: 2 x 2 = 4
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a2/brick
Brick2: serverb:/bricks/brick-b2/brick
Brick3: serverc:/bricks/brick-c2/brick
Brick4: serverd:/bricks/brick-d2/brick
Options Reconfigured:
performance.readdir-ahead: on
```

Here I would choose `serverc:/bricks/brick-c2/brick`  and `serverd:/bricks/brick-d2/brick` because it's a replica set

To remove these brick, start the migration process (this migrates any data that needs to)

```
[root@servera ~]# gluster volume remove-brick myvol serverc:/bricks/brick-c2/brick serverd:/bricks/brick-d2/brick start
volume remove-brick start: success
ID: 687c1cd3-a35d-467c-a36b-12f0c0223b74
```

Check the status to make sure it's complete

```
[root@servera ~]# gluster volume remove-brick myvol serverc:/bricks/brick-c2/brick serverd:/bricks/brick-d2/brick status
                                    Node Rebalanced-files          size       scanned      failures       skipped               status   run time in secs
                               ---------      -----------   -----------   -----------   -----------   -----------         ------------     --------------
                 serverc.lab.example.com               51         2.0KB            51             0             0            completed               3.00
                 serverd.lab.example.com                0        0Bytes             0             0             0            completed               0.00
```

Once it's complete; you can commit the change

```
[root@servera ~]# gluster volume remove-brick myvol serverc:/bricks/brick-c2/brick serverd:/bricks/brick-d2/brick commit
Removing brick(s) can result in data loss. Do you want to Continue? (y/n) y
volume remove-brick commit: success
Check the removed bricks to ensure all files are migrated.
If files with data are found on the brick path, copy them via a gluster mount point before re-purposing the removed brick.
```

If you want (like the msg says)...make sure the bricks are empty

```
[root@servera ~]# ssh root@serverc -- ls -l /bricks/brick-c2/brick
root@serverc's password:
total 0

[root@servera ~]# ssh root@serverd -- ls -l /bricks/brick-d2/brick
root@serverd's password:
total 0
```
## Configuring IP Failover

When using Samba to export Gluster volumes, the Samba Clustered Trivial Database (CTDB) can be used to provide IP failover for Samba shares. This can be used to provide HA Samba shares.

When using local user accounts, the local user databases can also be synced between nodes.

First (on *all* gluster nodes); install the needed packages

```
yum -y install samba ctdb
```

Next (on *all* gluster nodes); setup the firewall to allow `samba` services and port `4379` over tcp

```
firewall-cmd --permanent --add-service=samba
firewall-cmd --permanent --add-port=4379/tcp
firewall-cmd --reload
```

Now, on one of the gluster node, stop the `ctdbmeta` volume

```
[root@servera ~]# gluster volume stop ctdbmeta
```

Next (on *all* gluster nodes); edit `/var/lib/glusterd/hooks/1/start/post/S29CTDBsetup.sh` and `/var/lib/glusterd/hooks/1/stop/pre/S29CTDB-teardown.sh`  so that the `META="all"` line becomes `META="ctdbmeta"`

```
grep 'META=' /var/lib/glusterd/hooks/1/start/post/S29CTDBsetup.sh
META="ctdbmeta"

grep 'META=' /var/lib/glusterd/hooks/1/stop/pre/S29CTDB-teardown.sh
META="ctdbmeta"
```

Next (on *all* gluster nodes); edit the `/etc/samba/smb.conf` to include `clustering=yes` inside the `[global]` block

```
grep 'clustering' /etc/samba/smb.conf
clustering=yes
```

Now, on one server, start the `ctdbmeta` volume

```
[root@servera ~]# gluster volume start ctdbmeta
```

Next (on *all* gluster nodes);  create the file `/etc/ctdb/nodes` with the IP addresses of all the gluster servers

```
cat /etc/ctdb/nodes
172.25.250.10
172.25.250.11
```

Next (on *all* gluster nodes);  create the file `/etc/ctdb/public_addresses` with the floating IP and interface you want it on.

```
cat /etc/ctdb/public_addresses
172.25.250.15/24 eth0
```

Start and enable the ctdb service (ALL gluster nodes

```
systemctl enable ctdb
systemctl start ctdb
systemctl status ctdb
```

Now (on one server) create a samba user

```
[root@servera ~]# smbpasswd -a smbuser
```

Set the following options for the volume you want to export (again, on one server)

```
[root@servera ~]# gluster volume set custdata stat-prefetch off
volume set: success
[root@servera ~]# gluster volume set custdata server.allow-insecure on
volume set: success
[root@servera ~]# gluster volume set custdata storage.batch-fsync-delay-usec 0
volume set: success
```

On ALL servers, add `option rpc-auth-allow-insecure on` to the `/etc/glusterfs/glusterd.vol` file, inside the `type mgmt` block, then restart `glusterd.service`

```
grep rpc-auth-allow-insecure /etc/glusterfs/glusterd.vol
    option rpc-auth-allow-insecure on

systemctl restart glusterd
```

Now restart the volume you're exporting (just one server)

```
yes | gluster volume stop custdata
gluster volume start custdata
```

Try mounting it by putting it in your `/etc/fstab`

```
//172.25.250.15/gluster-custdata /mnt/custdata cifs user=smbuser,pass=redhat 0 0
```
## Configuring NFS Ganesha

The built-in NFS server for Gluster Storage supports NFSv3. If NFSv4 or IP failover is needed, administrators should use NFS-Ganesha. Note that NFS-Ganesha cannot be run with the built-in NFSv3 server. NFSv3 should be disabled on all nodes that will be running NFS-Ganesha. 

To get started; install the packages on ALL gluster servers by first stoping the `glusterd` service (here `serverX` stands for all servers)

```
[root@serverX ~]# systemctl stop glusterd
[root@serverX ~]# yum install glusterfs-ganesha
```

Update the firewall rules on ALL servers

```
[root@serverX ~]# firewall-cmd --permanent --add-service=high-availability --add-service=nfs --add-service=rpc-bind --add-service=mountd
[root@serverX ~]# firewall-cmd --reload
```

Now on just one server (here `servera` for example); copy over the sample `/etc/ganesha/ganesha-ha.conf.sample` file over to the right location.

```
[root@servera ~]# cp /etc/ganesha/ganesha-ha.conf.sample /etc/ganesha/ganesha-ha.conf
```

Update `/etc/ganesha/ganesha-ha.conf` with the following

```
HA_NAME="gls-ganesha"
HA_VOL_SERVER="servera"
HA_CLUSTER_NODES="servera.lab.example.com,serverb.lab.example.com"
VIP_servera_lab_example_com="172.25.250.16"
VIP_serverb_lab_example_com="172.25.250.17"
```

Now copy this config to all servers

```
[root@servera ~]# scp /etc/ganesha/ganesha-ha.conf serverX:/etc/ganesha/
```

On ALL servers, enable the proper clustering services

```
[root@serverX ~]# systemctl enable pacemaker pcsd
```

Start `pcsd` on ALL servers

```
[root@serverX ~]# systemctl start pcsd
```

On both your systems, set the password for the `hacluster` user

```
[root@serverX ~]# echo -n redhat | passwd --stdin hacluster
```

Now from one of your servers, authenticate `pcs` communication between all nodes. 

```
[root@servera ~]# pcs cluster auth -u hacluster -p redhat servera.lab.example.com serverb.lab.example.com
```

Create an ssh keypair on one of your servers and copy it over to all your other servers

```
[root@servera ~]# ssh-keygen -f /var/lib/glusterd/nfs/secret.pem -t rsa -N ''
[root@servera ~]# scp /var/lib/glusterd/nfs/secret.pem* serverb:/var/lib/glusterd/nfs/
[root@servera ~]# ssh-copy-id -i /var/lib/glusterd/nfs/secret.pem.pub root@servera
[root@servera ~]# ssh-copy-id -i /var/lib/glusterd/nfs/secret.pem.pub root@serverb
```

Now start `glusterd` on ALL nodes

```
[root@serverX ~]# systemctl start glusterd
```

Next, you'll need to use Gluster's built-in shared clustered storage mechanism (only needs to be done on one server)

```
[root@servera ~]# gluster volume set all cluster.enable-shared-storage enable
```

On ALL servers, Configure nfs-ganesha to use the default port (20048/tcp/20048/UDP) for `mountd` in the `NFS_core_Param` in `/etc/ganesha/ganesha.conf`. IN the end the file should look like this

```
[root@serverX ~]# grep -A9 NFS_Core_Param /etc/ganesha/ganesha.conf
NFS_Core_Param {
        #Use supplied name other tha IP In NSM operations
        NSM_Use_Caller_Name = true;
        #Copy lock states into "/var/lib/nfs/ganesha" dir
        Clustered = false;
        #Use a non-privileged port for RQuota
        Rquota_Port = 4501;
        MNT_Port = 20048;
}
```

Next, on one of the servers, enable `nfs-ganesha` and export `custdata`

```
[root@servera ~]# yes | gluster nfs-ganesha enable
[root@servera ~]# gluster volume set custdata ganesha.enable on
```

From the client, verify that the `custdata` volume is exported over NFS using the Virtual IP

```
[student@workstation ~]$ showmount -e 172.25.250.16
Export list for 172.25.250.16:
/custdata (everyone)
```

Mount it if you wish

```
[root@workstation ~]# mkdir /mnt/nfs
[root@workstation ~]# echo "172.25.250.16:/custdata     /mnt/nfs        nfs     rw,vers=4       0 0" >> /etc/fstab
[root@workstation ~]# mount /mnt/nfs
[root@workstation ~]# df -h /mnt/nfs/
Filesystem               Size  Used Avail Use% Mounted on
172.25.250.16:/custdata  2.0G   33M  2.0G   2% /mnt/nfs
```

## Georeplication

Georeplication can be configured between volumes on the same host, or between a local volume, and a volume on a remote host. This can be connected using a LAN in the same data center, over a WAN, or even over the Internet. Georeplication can also be cascaded. A volume can be synced to more than one slave, and/or each of those slaves can then be synchronized to one or more slaves.

* [Setup](#georeplication-setup)
* [Management](#georeplication-management)

### Georeplication Setup

Before georeplication can be configured, a number of prerequisites must be met.

* The master/slaves must be on the same version of Gluster.
* The slave can not be a peer of any of the nodes on the master system.
* Passwordless SSH access is required between the root account on one of the nodes of the master volume, the node where the geo-replication create command will be run, and the account that will be used for georeplication on the slave node.

Let's get started!

First you need to enable shared storage for the `georeplication` daemon

```
[root@servera ~]# gluster volume set all cluster.enable-shared-storage enable
volume set: success
```

Set up ssh-keys from the `root` account on the master to the `geoaccount` on the slave system.

```
[root@servera ~]# ssh-keygen -f ~/.ssh/id_rsa -N ''
```

Now, copy this id over

```
[root@servera ~]# ssh-copy-id -i ~/.ssh/id_rsa.pub geoaccount@servere
```

On the slave system, create a new dir called `/var/mountbroker-root`. This must be `0711` permission and have SELinux context equal to `/home` 

```
[root@servere ~]# mkdir -m 0711 /var/mountbroker-root
[root@servere ~]# semanage fcontext -a -e /home /var/mountbroker-root
[root@servere ~]# restorecon -vR /var/mountbroker-root
restorecon reset /var/mountbroker-root context unconfined_u:object_r:var_t:s0->unconfined_u:object_r:home_root_t:s0
```

On the slave system, configure the following options; then restart `glusterd`

* Set the `mountbroker-root` dir to `/var/mountbroker-root`
* Set the mountbroker user for `slavevol` to `geoaccount`
* Set `geo-replication-log-group` group to `geogroup`
* Allow RPC connections

```
[root@servere ~]# gluster system:: execute mountbroker opt mountbroker-root /var/mountbroker-root
Command executed successfully.

[root@servere ~]# gluster system:: execute mountbroker user geoaccount slavevol
Command executed successfully.

[root@servere ~]# gluster system:: execute mountbroker opt geo-replication-log-group geogroup
Command executed successfully.

[root@servere ~]# gluster system:: execute mountbroker opt rpc-auth-allow-insecure on
Command executed successfully.

[root@servere ~]# systemctl restart  glusterd
```

Configure, and start, georeplication between the master and the slave, using the `geoaccount` account.

On the master

```
[root@servera ~]# gluster system:: execute gsec_create
Common secret pub file present at /var/lib/glusterd/geo-replication/common_secret.pem.pub

[root@servera ~]# gluster volume geo-replication mastervol geoaccount@servere::slavevol create push-pem
Creating geo-replication session between mastervol & geoaccount@servere::slavevol has been successful
```

Then on the slave

```
[root@servere ~]# /usr/libexec/glusterfs/set_geo_rep_pem_keys.sh geoaccount mastervol slavevol
Successfully copied file.
Command executed successfully.
```

Back on the master

```
[root@servera ~]# gluster volume geo-replication mastervol geoaccount@servere::slavevol config use_meta_volume true
geo-replication config updated successfully

[root@servera ~]# gluster volume geo-replication mastervol geoaccount@servere::slavevol start
Starting geo-replication session between mastervol & geoaccount@servere::slavevol has been successful
```

On the master; verify that geo-rep is running

```
[root@servera ~]# gluster volume geo-replication status

MASTER NODE                MASTER VOL    MASTER BRICK              SLAVE USER    SLAVE                                 SLAVE NODE    STATUS     CRAWL STATUS       LAST_SYNCED
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
servera.lab.example.com    mastervol     /bricks/brick-a1/brick    geoaccount    ssh://geoaccount@servere::slavevol    servere       Active     Changelog Crawl    2018-05-29 14:01:34
serverb.lab.example.com    mastervol     /bricks/brick-b1/brick    geoaccount    ssh://geoaccount@servere::slavevol    servere       Passive    N/A                N/A
```

On the slave; check a brick to see the files being copied over

```
[root@servere ~]# ll /bricks/brick-e1/brick/ | head
total 0
-rw-r--r--. 2 root root 0 May 29 13:21 file00
-rw-r--r--. 2 root root 0 May 29 13:21 file01
-rw-r--r--. 2 root root 0 May 29 13:21 file02
-rw-r--r--. 2 root root 0 May 29 13:21 file03
-rw-r--r--. 2 root root 0 May 29 13:21 file04
-rw-r--r--. 2 root root 0 May 29 13:21 file05
-rw-r--r--. 2 root root 0 May 29 13:21 file06
-rw-r--r--. 2 root root 0 May 29 13:21 file07
-rw-r--r--. 2 root root 0 May 29 13:21 file08
```
### Georeplication Management

Various options of a georeplication configuration can be modified. This includes settings for the location of log files, if deleted files should be deleted on the slave, and more.


The following shows (for an example)

* The changelog changed to be parsed every 5 seconds to see if any new changes need to be synchronized.
* Files deleted on `mastervol` should not be deleted from `slavevol`.
* A new checkpoint should be created, with the current date and time. 

1. To change the changelog; first unmount the share from the client

```
[student@workstation ~]$ sudo umount /mnt/mastervol
```

Then you can set the rollover time on the master

```
[root@servera ~]# gluster volume  set mastervol changelog.rollover-time 5
volume set: success
```

You can now remount the clients

```
[student@workstation ~]$ sudo mount /mnt/mastervol/
[student@workstation ~]$ df -h /mnt/mastervol/
Filesystem          Size  Used Avail Use% Mounted on
servera:/mastervol  2.0G   33M  2.0G   2% /mnt/mastervol
```
2. To keep deleted files on the slave

Set the `ignore-deletes` to `true` on the master

```
[root@servera ~]# gluster volume geo-replication mastervol geoaccount@servere::slavevol config ignore-deletes true
geo-replication config updated successfully
```

Test it by removing a file on the client

```
[student@workstation ~]$ sudo rm -f /mnt/mastervol/importantfile
[student@workstation ~]$ sudo ls -d /mnt/mastervol/importantfile
ls: cannot access /mnt/mastervol/importantfile: No such file or directo
```

On the master, check the status of `LAST_SYNCED`

```
[root@servera ~]# gluster volume geo-replication status

MASTER NODE                MASTER VOL    MASTER BRICK              SLAVE USER    SLAVE                                 SLAVE NODE    STATUS     CRAWL STATUS       LAST_SYNCED
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
servera.lab.example.com    mastervol     /bricks/brick-a1/brick    geoaccount    ssh://geoaccount@servere::slavevol    servere       Active     Changelog Crawl    2018-05-29 14:47:43
serverb.lab.example.com    mastervol     /bricks/brick-b1/brick    geoaccount    ssh://geoaccount@servere::slavevol    servere       Passive    N/A                N/A
```

Once it's synced; verify the file is still there on the slave

```
[root@servere ~]# ls -ld /bricks/brick-e1/brick/importantfile
-rw-r--r--. 2 root root 6 May 29 14:32 /bricks/brick-e1/brick/importantfile
```

3. To create a checkpoint for the geo-rep to start from

You can set it to the keyword `now`, on the master

```
[root@servera ~]# gluster volume geo-replication mastervol geoaccount@servere::slavevol config checkpoint now
geo-replication config updated successfully
```

Next, Verify that a checkpoint has been created (Inspect the `CHECKPOINT TIME` column). 

```
gluster volume geo-replication mastervol geoaccount@servere::slavevol status detail
```

## Basic Troubleshooting

Here are some things you can check when something goes wrong

* [Defective Bricks](#defective-bricks)
* [BitRot](#bitrot)

### Defective Bricks

When a brick has been offline, you will need to "heal" the volume. There are self heal mechanisms you can use. Also you can heal files that have gone into a "split brain" scenario

For example, if `servera` went down; and you booted it back up; you can trigger a heal. First check the status on `serverb`

```
[root@serverb ~]# gluster volume heal replvol info
Brick servera:/bricks/brick-a1/brick
Status: Transport endpoint is not connected

Brick serverb:/bricks/brick-b1/brick
Number of entries: 0
```

After making sure `servera` is back online...

```
[root@servera ~]# gluster volume heal replvol info
Brick servera:/bricks/brick-a1/brick
Number of entries: 0

Brick serverb:/bricks/brick-b1/brick
Number of entries: 0
```

Make sure this says that it has no more files to heal


You can replace a brick in `replvol` on `servera` with another brick. First check the info

```
[root@servera ~]# gluster volume info replvol

Volume Name: replvol
Type: Replicate
Volume ID: c7f2848e-e588-4055-b45e-3de96936c17f
Status: Started
Number of Bricks: 1 x 2 = 2
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a1/brick
Brick2: serverb:/bricks/brick-b1/brick
Options Reconfigured:
performance.readdir-ahead: on
```

Then, replace the brick

```
[root@servera ~]# gluster volume replace-brick replvol servera:/bricks/brick-a1/brick servera:/bricks/brick-a2/brick commit force
```

Verify that the brick has been replaced

```
[root@servera ~]# gluster volume info replvol

Volume Name: replvol
Type: Replicate
Volume ID: c7f2848e-e588-4055-b45e-3de96936c17f
Status: Started
Number of Bricks: 1 x 2 = 2
Transport-type: tcp
Bricks:
Brick1: servera:/bricks/brick-a2/brick
Brick2: serverb:/bricks/brick-b1/brick
Options Reconfigured:
performance.readdir-ahead: on
```
### Bitrot

In gluster; when BitRot Detection is enabled, all files on a volume will be scrubbed at regular intervals, and a checksum will be calculated/verified.


To enable BitRot detection

```
[root@servera ~]# gluster volume bitrot replvol enable
```

Configure to scan all files once an hour

```
[root@servera ~]# gluster volume bitrot replvol scrub-frequency hourly
```

Set to scan the maximum number of files at once

```
[root@servera ~]# gluster volume bitrot replvol scrub-throttle aggressive
```
