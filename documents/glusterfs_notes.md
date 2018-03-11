# GlusterFS Notes

These are my GlusterFS notes in no paticular order. More infrmation can be found on the [official documentation page](https://access.redhat.com/documentation/en-us/red_hat_gluster_storage/)

* [Create Bricks with XFS](#create-bricks-with-xfs)

## Create Bricks With XFS

Created "bricks" with LVM (leaving room for snapshots) and fromat with [xfs](MYXFSDOCS)

__Notes:__
  * Install XFS filesystem (with EPEL Repo if needed)
  * GlusterFS uses extended attributes on files, you must increase the inode size to 512 bytes (default is 256 bytes)

This example uses `/dev/sdb`...your device might, obviously, be different.
```
pvcreate /dev/sdb
vgcreate gfs_brick /dev/sdb
lvcreate -L +nG -n lv_brick gfs_brick
mkfs.xfs -i size=512 /dev/gfs_brick/lv_brick
mount.xfs  /dev/gfs_brick/lv_brick /exp
echo "/dev/gfs_brick/lv_brick	/exp	xfs	defaults	0 0" >> /etc/fstab
```

Make sure that the service is up and it starts on boot

```
service glusterd start
chkconfig glusterd on
```

