# GlusterFS Notes

These are glusterfs notes in no paticular order. Old 3.1 notes can be found [here](glusterfs_notes.md)

* [Installation](#installation)
* [Building a Trusted Storage Pool](#building-a-trusted-storage-pool)
* [Creating Bricks](#creating-bricks)
* [Building a Volume](#building-a-volume)

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
