# GlusterFS Notes

These are my GlusterFS notes in no paticular order. More infrmation can be found on the [official documentation page](https://access.redhat.com/documentation/en-us/red_hat_gluster_storage/)

* [Create Bricks with XFS](#create-bricks-with-xfs)
* [Probing Peers](#probing-peers)
* [Creating Volumes](#creating-volumes)
* [Mounting](#mounting)

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

