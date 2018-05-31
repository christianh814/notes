# Encrypted Volumes

We will be using an LVM volume but this can be done with "raw" discs/partitions. you will just skip the LVM part and replace the disk with `/dev/sdX` ...furthermore these notes can be used as "create LVM volume" notes as well (rezise notes will be elseware)

First initialize the physical device

```
root@host# pvcreate /dev/sdb
root@host# pvs
  PV         VG   Fmt  Attr PSize  PFree
  /dev/sda2  vg1  lvm2 a-   14.51g    0 
  /dev/sdb        lvm2 a-    1.00g 1.00g
```

Now create a "Volume Group" and give it a name 

```
root@host# vgcreate vg2 /dev/sdb 
root@host# vgs
  Volume group "vg2" successfully created
  VG   #PV #LV #SN Attr   VSize    VFree   
  vg1    1   2   0 wz--n-   14.51g       0 
  vg2    1   0   0 wz--n- 1020.00m 1020.00m
```

You can now create a Logical Volume within the volume group

```
root@host# lvcreate -l 100%VG -n lockbox vg2
  Logical volume "lockbox" created
root@host#  lvs
  LV       VG   Attr   LSize    Origin Snap%  Move Log Copy%  Convert
  lvm_root vg1  -wi-ao   10.54g                                      
  lvm_swap vg1  -wi-ao    3.97g                                      
  lockbox  vg2  -wi-a- 1020.00m
```

Now here you would have formated the drive and mounted it on a directory; but since we want to encrypt the filesystem, we have to perform some steps before the format.

First you fill the "raw" disc blocks with random data (this will fill up the "raw" disc until it's full)

```
root@host# # cat /dev/urandom > /dev/vg2/lockbox 
cat: write error: No space left on device
```

Next you must use the "cryptsetup" command to format with the encryption layer. It will ask you to enter a passphrase **REMEMBER THIS PASSPHRASE**

```
root@host# # cryptsetup luksFormat /dev/vg2/lockbox 
WARNING!
========
This will overwrite data on /dev/vg2/lockbox irrevocably.

Are you sure? (Type uppercase yes): YES 
Enter LUKS passphrase: 
Verify passphrase:
```

Now you will use the `luksOpen` subcommand to "map" the encrypted partiotion to a readable (i.e. unecrypted) layer. It will prompt you for the passphrase. 

```
root@host# cryptsetup luksOpen /dev/vg2/lockbox lockbox
Enter passphrase for /dev/vg2/lockbox: 

root@host# ll /dev/mapper/lockbox 
lrwxrwxrwx. 1 root root 7 Aug 30 21:19 /dev/mapper/lockbox -> ../dm-3
```

Now that we have the unencrypted layer mapped; we can now format this layer as if it were a "regular" partion.

```
root@host# mkfs.ext4 /dev/mapper/lockbox 
mke2fs 1.41.12 (17-May-2010)
Filesystem label=
OS type: Linux
Block size=4096 (log=2)
Fragment size=4096 (log=2)
Stride=0 blocks, Stripe width=0 blocks
65152 inodes, 260608 blocks
13030 blocks (5.00%) reserved for the super user
First data block=0
Maximum filesystem blocks=268435456
8 block groups
32768 blocks per group, 32768 fragments per group
8144 inodes per group
Superblock backups stored on blocks: 
        32768, 98304, 163840, 229376

Writing inode tables: done                            
Creating journal (4096 blocks): done
Writing superblocks and filesystem accounting information: done

This filesystem will be automatically checked every 35 mounts or
180 days, whichever comes first.  Use tune2fs -c or -i to override.
```

You can go ahead and mount it as if it were a regular partion

```
root@host# mount /dev/mapper/lockbox /mnt/lockbox
```

If you want the volume to mount on boot you must enter this in the crypttab file.
The format of the `/etc/crypttab` file is: `MAPNAME        RAWDEV          none`

```
root@host# cat /etc/crypttab
lockbox         /dev/vg2/lockbox        none
```


Now enter it in the fstab file as any other partition 

```
root@host# vi /etc/fstab
        /dev/mapper/lockbox     /mnt/lockbox    ext4    defaults        0 0
```

**NOTE** This will prompt you for a password on boot

