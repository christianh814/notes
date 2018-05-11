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

### Quotas
