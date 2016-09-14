# Create LVM Snapshot

 

If you are using LVM (Logical Volume Manager) for Linux; you can leverage it's snapshot ability. Snapshots are a good way of having a preserved copy of a filesystem; perfect for backups

First you have to scan to see which volume you want to take a snapshot of (in this case; the root filesystem)

	[root@host ~]# lvs
    	LV      VG   Attr     LSize Pool Origin Data%  Move Log Copy%  Convert
    	lv_root vg0  -wi-ao-- 6.50g                                           
    	lv_swap vg0  -wi-ao-- 1.00g                                           

You can see that the logical volume (LV) is named lv_root and "lives" in VG (volume group) that's named vg0. Run a scan to see if you have enough room to create a snapshot.

	[root@host ~]# vgs
    	VG   #PV #LV #SN Attr   VSize VFree
    	vg0    1   2   0 wz--n- 9.50g 2.00g

You can see that we have about 2GB free for a snapshot (MORE than enough); on vg0. Now you can create a snapshot with enough room (keeping in mind that LVM uses COW for the snapshot) for the snapshot to live.

	[root@host ~]# lvcreate --snapshot /dev/vg0/lv_root --name lv_root_snapshot --size 15M
    	Rounding up size to full physical extent 32.00 MiB
    	Logical volume "lv_root_snapshot" created

The basic things you'll need are: the name of the LV, what you want to name the snapshot, and the size. Now that you've got your snapshot...take a look at it.

	[root@host ~]# lvs
    	LV               VG   Attr     LSize  Pool Origin  Data%  Move Log Copy%  Convert
    	lv_root          vg0  owi-aos-  6.50g                                            
    	lv_root_snapshot vg0  swi-a-s- 32.00m      lv_root   0.31                        
    	lv_swap          vg0  -wi-ao--  1.00g                                            

Now you can mount it like any 'ol filesystem (making sure you make it read only)

	[root@host ~]# mount -o rw /dev/vg0/lv_root_snapshot /mnt

Now that it's mounted...you can see it in your df. You can also perform backups on that filesystem

	[root@host ~]# df -hF ext3
	Filesystem                        Size  Used Avail Use% Mounted on
	/dev/mapper/vg0-lv_root           6.5G  4.1G  2.4G  64% /
	/dev/sda1                         497M   83M  390M  18% /boot
	/dev/mapper/vg0-lv_root_snapshot  6.5G  4.1G  2.4G  64% /mnt

To remove the snapshot; unmount the filesystem

	[root@host ~]# umount /mnt

Then you can remove the snaphost for LVM

	[root@host ~]# lvremove vg0/lv_root_snapshot
	Do you really want to remove active logical volume lv_root_snapshot? [y/n]: y
    	Logical volume "lv_root_snapshot" successfully removed

Check to see if the snapshot is gone

	[root@host ~]# lvs
    	LV      VG   Attr     LSize Pool Origin Data%  Move Log Copy%  Convert
    	lv_root vg0  -wi-ao-- 6.50g                                           
    	lv_swap vg0  -wi-ao-- 1.00g

# LVM - roll back

Before the upgrade the snapshot can be created. The snapshot should be big enough to cover all changes done by the upgrade. Furthermore, since /boot is on a separate filesystem, we create a copy of it first and include it in the snapshot.

	
	root@host# mkdir -p /var/backup  
	root@host# cp -r /boot /var/backup  
	root@host# lvcreate -s -L 2G -n lvrootsnap /dev/sysvg/lvroot


Perform the changes to the system now, for example initiated by yum update. You can also reboot the system and test if the upgrade was successful.

If the upgrade is considered successful then remove the snapshot and delete the backup directory

	
	root@host# lvremove /dev/sysvg/lvrootsnap 
	root@host# rm -r /var/backup


NOW. If the upgrade was not successful merge the snapshot into the original volume: The snapshot has to be marked for merging into the original volume. At this state the root volume is in use by the system. The merging can only be performed when the volume is not in use, so the following command will only mark the snapshot to be merged upon the next activation of the root volume. We will also remove the possibly changed contents of the /boot directory and place the old contents there. For reinstallation of grub change vda in the example to the device of the systems harddisk

	
	root@host# lvconvert --merge /dev/sysvg/lvrootsnap
	root@host# rm -rf /boot/*
	root@host# cp -r /var/backup/boot/* /boot && rm -r /var/backup
	root@host# grub-install /dev/vda


# LVM Import Volume

If you move disks from one server to another; you can import the volume groups and logical volume.

First you need to "see" the new disks

	
	for hostfile in ls -1d /sys/class/scsi_host/host*
	do
	  echo "- - -" > ${hostfile}/scan
	done


Now that you can see the disks; scan for volume groups

	
	[root@r ~]# pvscan 
	  PV /dev/sda2   VG VolGroup00   lvm2 [74.88 GB / 0    free]
	  Total: 1 [74.88 GB] / in use: 1 [74.88 GB] / in no VG: 0 [0   ]
	[root@r ~]# vgscan 
	  Reading all physical volumes.  This may take a while...
	  Found volume group "VolGroup00" using metadata type lvm2
	  Found volume group "store1" using metadata type lvm2


You should be able to see your volume group now

	
	[root@r ~]# vgs
	  VG         #PV #LV #SN Attr   VSize   VFree
	  VolGroup00   1   2   0 wz--n- 74.88G     0 
	  store1       1   2   0 wz--n- 974.88G    0 

    
Now import the volume group

	
	[root@r ~]# vgimport store1


Next import the logical volume

	
	[root@r ~]# vgchange -ay store1


Now you should be able to mount the volume.

	
	[root@r ~]# mount -o ro /dev/store1/store1 /mnt


Remember (if necessary ) the update the /etc/fstab

# Create Swapfile

Sometimes you need create a swapfile in order to quickly add more swap space.

First create a 2G file

	
	[root@host ~]# dd if=/dev/urandom of=/var/swapfile bs=512k count=4096


Now format it for swap

	
	[root@host ~]# mkswap -f /var/swapfile 
	Setting up swapspace version 1, size = 2097148 KiB
	no label, UUID=2dde6cc4-f905-4c41-9ccc-1a36e2171d18


Now add this to the FSTAB file

	
	/var/swapfile	swap	swap	sw,pri=1	0 0


Now activate it with the "swapon" command.

	
	[root@host ~]# free -m
	             total       used       free     shared    buffers     cached
	Mem:           996        927         69          0          1         28
	-/+ buffers/cache:        897         98
	Swap:         2015       1471        544
	[root@host ~]# swapon -a
	[root@host ~]# free -m
	             total       used       free     shared    buffers     cached
	Mem:           996        935         61          0          1         26
	-/+ buffers/cache:        906         89
	Swap:         4063       1469       2594


# Re-Read Partition Table

The command // partprobe // was commonly used in RHEL 5 to inform the OS of partition table changes on the disk. In RHEL 6, it will only trigger the OS to update the partitions on a disk that none of its partitions are in use (e.g. mounted). If any partition on a disk is in use, partprobe will not trigger the OS to update partitions in the system because it is considered unsafe in some situations.

If a new partition was added and none of the existing partitions were modified, consider using the // partx // command to update the system partition table. Do note that the *partx * command does not do much checking between the new and the existing partition table in the system and assumes the user knows what they are are doing. So it can corrupt the data on disk if the existing partitions are modified or the partition table is not set correctly. So use at one's own risk.

For example, a partition #1 is an existing partition and a new partition #2 is already added in // /dev/sdb // by // fdisk//. Here we use // partx -v -a /dev/sdb // to add the new partition to the system:

`<code>`# ls /dev/sdb*  
/dev/sdb  /dev/sdb1`</code>`

List the partition table of disk:

	
	# partx -l /dev/sdb
	# 1:        63-   505007 (   504945 sectors,    258 MB)  
	# 2:    505008-  1010015 (   505008 sectors,    258 MB)  
	# 3:         0-       -1 (        0 sectors,      0 MB)  
	# 4:         0-       -1 (        0 sectors,      0 MB)  


Read disk and try to add all partitions to the system:

	
	# partx -v -a /dev/sdb                                         
	device /dev/sdb: start 0 size 2097152  
	gpt: 0 slices  
	dos: 4 slices  
	# 1:        63-   505007 (   504945 sectors,    258 MB)  
	# 2:    505008-  1010015 (   505008 sectors,    258 MB)  
	# 3:         0-       -1 (        0 sectors,      0 MB)  
	# 4:         0-       -1 (        0 sectors,      0 MB)  
	BLKPG: Device or resource busy
	error adding partition 1


These last 2 lines are normal in this case because partition 1 is already added in the system before partition 2 is added

Check that we have device nodes for // /dev/sdb // itself and the partitions on it:

	
	# ls /dev/sdb*  
	/dev/sdb  /dev/sdb1  /dev/sdb2`</code>`
	
	#  CLVM
	
	CLVM stands for "Clusterd LVM" - meaning if you have an iscsi disk that two servers can see...you should be able to "see" the lvm volumes on both nodes.
	
	Make sure you have an iscsi disk available on both nodes. Check [[iscsi_notes|here]] if you don't know how
	
	First make sure you have the service installed and running (on both nodes).
	
	`<code>`
	root@host# yum -y install lvm2-cluster
	root@host# chkconfig clvmd on
	root@host# service clvmd start


Next make sure you have the cluster locking type to "3" in the // /etc/lvm/lvm.conf // file

	
	locking_type = 3


One way to do this (other than VI the file) is use the **lvmconf** command

	
	lvmconf --enable-cluster 


Once configured and // clvmd // is on all nodes - you should be able to "see" all you LVM configurations (try *lvs* and *vgs*) on all nodes.

NOTE: When creating your clusterd LVM device...create the volume group with the proper switches

	
	vgcreate --clusterd y vg0 /dev/sda


# Add Size From Disk

When you've "grown" a disk (from either VMWare or RHEV); you need to make the OS "see" the extended size...

You can do this by running either...

	
	/usr/bin/rescan-scsi-bus.sh


Or...

	
	for hostfile in $(ls -1d /sys/class/scsi_host/host*)
	do
	  echo "- - -" > ${hostfile}/scan
	done


Finally...

	
	disk=/dev/vdb
	echo "1" > /sys/block/${disk##*/}/device/rescan


Verify that your disk is the size you expect

	
	fdisk -cul /dev/vdb


Now Make sure you resize using pvresize...

	
	root@host# pvresize /dev/vdb 
	  Physical volume "/dev/vdb" changed
	  1 physical volume(s) resized / 0 physical volume(s) not resized
	root@host# pvs
	  PV         VG          Fmt  Attr PSize  PFree 
	  /dev/vda2  vg_gluster1 lvm2 a--  19.51g     0 
	  /dev/vdb   vg1         lvm2 a--  50.00g 40.00g


Now grow the logical volume (checking with vgs and lvs every step of the way)

	
	root@host# lvs
	  LV      VG          Attr      LSize  Pool Origin Data%  Move Log Cpy%Sync Convert
	  brick   vg1         -wi-ao--- 10.00g                                             
	  lv_root vg_gluster1 -wi-ao--- 15.57g                                             
	  lv_swap vg_gluster1 -wi-ao---  3.94g                                             
	root@host# lvextend -l +100%FREE vg1/brick
	  Extending logical volume brick to 50.00 GiB
	  Logical volume brick successfully resized
	root@host# lvs
	  LV      VG          Attr      LSize  Pool Origin Data%  Move Log Cpy%Sync Convert
	  brick   vg1         -wi-ao--- 50.00g                                             
	  lv_root vg_gluster1 -wi-ao--- 15.57g                                             
	  lv_swap vg_gluster1 -wi-ao---  3.94g


Now you can grow the filesystem

	
	root@host# resize2fs /dev/vg1/brick


For filesytems that are XFS (as opposed to EXT) consult the following to grow the filesystem [XFS GrowFS Notes](xfs#increasing_the_size_of_an_xfs_file_system)
