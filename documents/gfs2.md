# GFS2 Configuration

Make sure you have [CLVM](https://github.com/christianh814/notes/blob/master/documents/lvm_notes.md#clvm) set up before configuring gfs2

We are assuming 2 nodes like we have for the [Cluster Suite](./cluster_suite.md)

First make sure that the GFS2 kernel module is loaded on all nodes

```
root@host# modprobe gfs2
```

Now take you CLVM volume...

```
root@host# lvs
  LV      VG      Attr      LSize  Pool Origin Data%  Move Log Cpy%Sync Convert
  lv_root vg0     -wi-ao--- 17.54g                                             
  lv_swap vg0     -wi-ao---  1.97g                                             
  lv_clvm vg_clvm -wi-ao---  9.99g 
```

... and format it with gfs2

```
root@host# mkfs.gfs2 -j 2 -p lock_dlm -t mycluster:myvol /dev/vg_clvm/lv_clvm
 
This will destroy any data on /dev/vg_clvm/lv_clvm.
It appears to contain: symbolic link to `../dm-2'

Are you sure you want to proceed? [y/n] y

Device:                    /dev/vg_clvm/lv_clvm
Blocksize:                 4096
Device Size                9.99 GB (2619392 blocks)
Filesystem Size:           9.99 GB (2619389 blocks)
Journals:                  2
Resource Groups:           40
Locking Protocol:          "lock_dlm"
Lock Table:                "mycluster:myvol"
UUID:                      76cea934-1690-a898-a61c-e11d59ef37e2
```

Here are some options you need to know

* `-j` The journal count (i.e. there are 2 nodes in this cluster…hence the journal count should be 2)
* `-p` Name of the locking protocol
* `-t` Name of lock table
* The `mycluster` is the name of the cluster you get from [[cluster_suite#misc_notes|"clustat"]]. What's after the colon is any name you want (**must** be unique to gfs2 though)

Enter this to `/etc/fstab`

```
/dev/vg_clvm/lv_clvm    /gfs2   gfs2     defaults 0 0
```

Make sure the service starts on boot

```
root@host# service gfs2 restart ; chkconfig gfs2 on
```
