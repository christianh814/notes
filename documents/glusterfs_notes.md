# GlusterFS Notes

These are my GlusterFS notes in no paticular order. More infrmation can be found on the [official documentation page](https://access.redhat.com/documentation/en-us/red_hat_gluster_storage/)

* [Create Bricks with XFS](#create-bricks-with-xfs)
* [Probing Peers](#probing-peers)
* [Creating Volumes](#creating-volumes)
* [Mounting](#mounting)
* [GEO Replication](#geo-replication)
* [Object Storage](#object-storage)
* [Quick And Dirty Howto](#quick-and-dirty-howto)
* [Misc Commands](#misc-commands)

## Create Bricks With XFS

Created "bricks" with LVM (leaving room for snapshots) and fromat with [xfs](MYXFSDOCS)

__Notes:__
  * Install XFS filesystem (with EPEL Repo if needed)
  * GlusterFS uses extended attributes on files, you must increase the inode size to 512 bytes (default is 256 bytes)
  * Do this on all GlusterFS servers

This example uses `/dev/sdb`...your device (and size) might, obviously, be different.
```
pvcreate /dev/sdb
vgcreate gfs_brick /dev/sdb
lvcreate -L +10G -n lv_brick gfs_brick
mkfs.xfs -i size=512 /dev/gfs_brick/lv_brick
mount.xfs  /dev/gfs_brick/lv_brick /exp
echo "/dev/gfs_brick/lv_brick	/exp	xfs	defaults	0 0" >> /etc/fstab
```

Make sure that the service is up and it starts on boot

```
service glusterd start
chkconfig glusterd on
```

## Probing Peers


Make sure each server knows about each other... (NOTE: Don't probe yourself)

```
[root@gluster01 ~]# gluster peer probe gluster02.example.com
[root@gluster01 ~]# gluster peer probe gluster03.example.com
```

## Creating Volumes

Create volume - Distributed Replicated i.e. RAID 10

```
glusterN# gluster volume create export replica 2 transport tcp gluster1:/exp gluster2:/exp gluster3:/exp gluster4:/exp
```

Create volume - Replicated i.e. RAID 1

```
glusterN# gluster volume create demo1 replica 2 transport tcp gluster1:/demo1 gluster2:/demo1
```

Start the volume

```
glusterN# gluster volume start export
```

Setup Fast Failover

```
glusterN# gluster volume set export network.frame-timeout 60
glusterN# gluster volume set export network.ping-timeout 20
```

Now you can test Client

```
[root@gluster-client ~]# time df -h
Filesystem            Size  Used Avail Use% Mounted on
/dev/mapper/vg_glusterclient-lv_root
                      	5.5G  759M  4.5G  15% /
tmpfs                 499M     0  499M   0% /dev/shm
/dev/sda1             485M   31M  429M   7% /boot
gluster1:/export       16G   65M   16G   1% /export

real	0m0.003s
user	0m0.000s
sys	0m0.000s
```

SHUTDOWN gluster1 and gluster4

```
[root@gluster-client ~]# time df -h
Filesystem            Size  Used Avail Use% Mounted on
/dev/mapper/vg_glusterclient-lv_root
                      	5.5G  759M  4.5G  15% /
tmpfs                 499M     0  499M   0% /dev/shm
/dev/sda1             485M   31M  429M   7% /boot
gluster1:/export       16G   65M   16G   1% /export

real	0m22.867s
user	0m0.000s
sys	0m0.000s
```

## Mounting

Mounting with native client...

You can specify the following options when using the `mount -t glusterfs` command. Note that you need to separate all options with commas.

* `backup-volfile-servers=server-name` - name of the backup volfile server to mount the client. If this option is added while mounting fuse client, when the first volfile server fails, then the server specified in backup-volfile-servers option is used as volfile server to mount the client.

* `fetch-attempts=3` - number of attempts to fetch volume files while mounting a volume. This option is useful when you mount a server with multiple IP addresses or when round-robin DNS is configured for the server name.

```
mount -t glusterfs -o backup-volfile-servers=gluster3:gluster2,log-level=WARNING,log-file=/var/log/gluster.log,acl gluster1:/export /export
```

Common Options:

  * `log-level` - Logs only specified level or higher severity messages in the log-file.
  * `log-file` - Logs the messages in the specified file.
  * `ro` - Mounts the file system as read only.
  * `acl` - Enables POSIX Access Control List on mount.
  * `selinux` - Enables handling of SELinux xattrs through the mount point.
  * `background-qlen=n` - Enables FUSE to handle n number of requests to be queued before subsequent requests are denied. Default value of n is 64.
  * `enable-ino32` - this option enables file system to present 32-bit inodes instead of 64- bit inodes.

## GEO Replication

GEO Replication is how you'd copy to a remote location.

First Make sure you have NTP set up on both gluster systems

On the "master" and "slave" system...create SSH keys and copy them over.

```
glusterMaster# ssh-keygen -t rsa
glusterMaster# ssh-copy-id -i ~/.ssh/id_rsa.pub glusterN
glusterSlave# ssh-keygen -t rsa
glusterSlave# ssh-copy-id -i ~/.ssh/id_rsa.pub glusterN
```

Now (on the master) set up geo-replication

```
glusterMaster# gluster volume geo-replication demo1 root@gluster2::demo2 start
glusterMaster# gluster volume geo-replication demo1 root@gluster2::demo2 status
```

## Object Storage

GlusterFS can do Object Storage with the `swift` plugin/api; but it's actually easier just to use [minio](https://github.com/minio/minio) with Gluster as the backend storage. I actually have a how-to for OpenShift that's really easy. It can be found [here](https://github.com/christianh814/minio-openshift)

## Quick And Dirty Howto

This is taken from the [RHKB](https://access.redhat.com/articles/2352211) and it's meant as a quick reference rather than an extensive document.

### Install Gluster Nodes

Subscribe to the proper channels

```
subscription-manager repos --enable=rhel-7-server-rpms
subscription-manager repos --enable=rh-gluster-3-for-rhel-7-server-rpms
```

Install RHSS packages

```
yum -y install redhat-storage-server
```

### Configuring Gluster Nodes

1. Edit `/etc/hosts`

Configure each node with an IP address and a hostname. This part can be skipped if you are using DNS.

```
192.168.XX.11  rhs31-a (Primary)
192.168.XX.12  rhs31-b
192.168.XX.13  rhs31-c
192.168.XX.14  rhs31-d
192.168.XX.15 rhs-client
```

2. Configure passwordless SSH between the nodes

From the primary node:
```
sudo su -
ssh-keygen
ssh-copy-id rhs31-b
ssh-copy-id rhs31-c
ssh-copy-id rhs31-d
```

3. Verify that your trusted storage pool is empty

From the primary node:

```
# gluster peer status
No peers present
```

4. Add your nodes to the trusted storage pool

From the primary node:

```
# gluster peer probe rhs31-b
peer probe: success.
# gluster peer probe rhs31-c
peer probe: success.
# gluster peer probe rhs31-d
peer probe: success.
```

Verify that the nodes were added.

```
# gluster peer status
Number of Peers: 3

Hostname: rhs31-b
Uuid: 847623c6-9f39-4ea1-a5a7-f4cb89f47de5
State: Peer in Cluster (Connected)

Hostname: rhs31-c
Uuid: 254fdc9a-901c-428e-90ad-589dee2a4820
State: Peer in Cluster (Connected)

Hostname: rhs31-d
Uuid: 3f2ad138-61eb-4d4b-be65-8343749ac20b
State: Peer in Cluster (Connected)
```

5a. Configure a distributed-replicated volume on four nodes (rhs31-a through rhs31-d)

Run the following commands on all four nodes. When you see `brickN`, use `brick5` for rhs31-a, `brick6` for rhs31-b, `brick7` for rhgs31-c, and `brick8` for rhgs31-d. If you haven't yet configured partitions and volume groups, refer to Appendix A2.

Create a logical volume in the thin-pool volume group.

```
# lvcreate -V 2G -T bricks/thin-pool -n brickN
Logical volume "brickN" created.
```

Create the file system on your new brick.

```
# mkfs.xfs -i size=512 /dev/bricks/brickN
```

Create the mount point.

```
# mkdir -p /bricks/brickN
```

Edit `/etc/fstab` so that your new file system mounts on boot.

```
/dev/bricks/brickN  /bricks/brickN    xfs   rw,noatime,inode64,nouuid   1   2
```

Mount the new file system, and ensure it mounted correctly.

```
# mount -a
# df -h /bricks/brickN
Filesystem                 Size  Used Avail Use% Mounted on
/dev/mapper/bricks-brickN  2.0G   33M  2.0G   2% /bricks/brickN
```

Create a `brick` directory on the mounted file system.

```
# mkdir /bricks/brickN/brick
```

5b. Create the distributed-replicated volume

```
# gluster volume create dis-rep replica 2 \
rhs31-a:/bricks/brick5/brick              \
rhs31-b:/bricks/brick6/brick              \
rhs31-c:/bricks/brick7/brick              \
rhs31-d:/bricks/brick8/brick
volume create: dis-rep: success: please start the volume to access data

# gluster volume list
dis-rep
dist-vol
rep-vol

# gluster volume start dis-rep
volume start: dis-rep: success

# gluster volume status dis-rep
Status of volume: dis-rep
Gluster process                             TCP Port  RDMA Port  
```

### Configuring a Gluster Client

On your client machine, install the required packages.

```
# yum install glusterfs-fuse -y
```

This installs `glusterfs-libs`, `glusterfs`, `glusterfs-client-xlators`, and `glusterfs-fuse`. Mount the gluster file system.

You can either do this manually like so:

```
# mount -t glusterfs -o acl rhs31-a:/dist-vol /gmount/fuse
```

Or configure the gluster file system to mount on boot by adding an entry similar to the following to the /etc/fstab file:

```
rhs31-a:/dist-vol   /gmount/fuse    glusterfs   _netdev   0   0
```


### Taking a Gluster Volume Snapshot

You can take a snapshot of a gluster volume with the gluster snapshot create command:

```
# gluster snapshot create snap-dis-rep dis-rep
snapshot create: success: Snap snap-dis-rep_GMT-2016.05.09-06.31.16 created successfully
```

To find out what snapshots have already been made, use gluster snapshot list:
```
# gluster snapshot list
snap-dis-rep_GMT-2016.05.09-06.31.16
snap-dis-rep_GMT-2016.05.09-06.33.18
snap-dis-rep_GMT-2016.05.09-06.34.14
```

To delete a previously taken snapshot, use gluster snapshot delete:
```
# gluster snapshot delete snap-dis-rep_GMT-2016.05.09-06.31.16
Deleting snap will erase all the information about the snap. Do you still want to continue? (y/n) y
snapshot delete: snap-dis-rep_GMT-2016.05.09-06.31.16: snap removed successfully
```

## Misc Commands

Grow XFS filesystem (with LVM)

```
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
```

When you get an error when trying to re-use bricks that you have removed from a volume that you're no loger using...

```
{path} or a prefix of it is already part of a volume
```

You can fix this by re-setting attributes

For the directory (or any parent directories) that was formerly part of a volume, simply (ON EVERY PEER):

```
root@gluserN# brick_path=/path/to/brick
root@gluserN# setfattr -x trusted.glusterfs.volume-id $brick_path
root@gluserN# setfattr -x trusted.gfid $brick_path
root@gluserN# rm -rf $brick_path/.glusterfs
```

You should be able to add the bricks to your new volume now

I've encountered an error where the volume won't mount. The log gives a `fuse not found error`

```
[2013-01-04 15:27:15.232943] E [mount.c:596:gf_fuse_mount] 0-glusterfs-fuse: cannot open /dev/fuse (No such file or directory)
```

You must make sure the fuse module is loaded.

```
root@host# modprobe fuse
root@host# mount /share/store
root@host# df -h | grep gluster
glusterfs#gluster1.la3.4over.com:/workflow2
```

-30-
