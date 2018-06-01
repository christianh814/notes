# Creating an XFS File System

To create an XFS file system, use the `mkfs.xfs /dev/device` command. In general, the default options are optimal for common use.
When using `mkfs.xfs` on a block device containing an existing file system, use the `-f` option to force an overwrite of that file system.

```
mkfs.xfs -f /dev/sda
```

For striped block devices (for example, RAID5 arrays), the stripe geometry can be specified at the time of file system creation. Using proper stripe geometry greatly enhances the performance of an XFS filesystem.
When creating filesystems on LVM or MD volumes, mkfs.xfs chooses an optimal geometry. This may also be true on some hardware RAIDs that export geometry information to the operating system.

To specify stripe geometry, use the following mkfs.xfs sub-options: 

`su=value` ~ Specifies a stripe unit or RAID chunk size. The value must be specified in bytes, with an optional k, m, or g suffix. 
`sw=value` ~ Specifies the number of data disks in a RAID device, or the number of stripe units in the stripe. 

The following example specifies a chunk size of 64k on a RAID device containing 4 stripe units:

```
mkfs.xfs -d su=64k,sw=4 /dev/sda
```

# Mounting an XFS File System

An XFS file system can be mounted with no extra options, for example:

```
mount /dev/sda /mnt
```

XFS also supports several mount options to influence behavior. 

XFS allocates inodes to reflect their on-disk location by default. However, because some 32-bit userspace applications are not compatible with inode numbers greater than 232, XFS will allocate all inodes in disk locations which result in 32-bit inode numbers. This can lead to decreased performance on very large filesystems (that is, larger than 2 terabytes), because inodes are skewed to the beginning of the block device, while data is skewed towards the end.
To address this, use the inode64 mount option. 

This option configures XFS to allocate inodes and data across the entire file system, which can improve performance:

```
mount -o inode64 /dev/sda /mnt
```

By default, XFS uses write barriers to ensure file system integrity even when power is lost to a device with write caches enabled. For devices without write caches, or with battery-backed write caches, disable the barriers by using the nobarrier option:

```
 mount -o nobarrier /dev/sda /mnt
```

# XFS Quota Management

The XFS quota subsystem manages limits on disk space (blocks) and file (inode) usage. XFS quotas control or report on usage of these items on a user, group, or directory or project level. Also, note that while user, group, and directory or project quotas are enabled independently, group and project quotas are mutually exclusive.

When managing on a per-directory or per-project basis, XFS manages the disk usage of directory hierarchies associated with a specific project. In doing so, XFS recognizes cross-organizational "group" boundaries between projects. This provides a level of control that is broader than what is available when managing quotas for users or groups.

XFS quotas are enabled at mount time, with specific mount options. Each mount option can also be specified as noenforce; this will allow usage reporting without enforcing any limits. Valid quota mount options are:

  * **uquota/uqnoenforce** - User quotas
  * **gquota/gqnoenforce** - Group quotas
  * **pquota/pqnoenforce** - Project quota

Once quotas are enabled, the `xfs_quota` tool can be used to set limits and report on disk usage. By default, `xfs_quota` is run interactively, and in basic mode. Basic mode sub-commands simply report usage, and are available to all users.

Basic `xfs_quota` sub-commands include:

  * `quota username/userID` ~ Show usage and limits for the given username or numeric userID
  * `df` ~ Shows free and used counts for blocks and inodes.

In contrast, `xfs_quota` also has an expert mode. The sub-commands of this mode allow actual configuration of limits, and are available only to users with elevated privileges. To use expert mode sub-commands interactively, run `xfs_quota` -x. 

Expert mode sub-commands include:

  * `report /path` ~ Reports quota information for a specific file system.
  * `limit` ~ Modify quota limits.

For a complete list of sub-commands for either basic or expert mode, use the sub-command `help`. All sub-commands can also be run directly from a command line using the `-c` option, with `-x` for expert sub-commands.

To display a sample quota report for `/home` (on `/dev/blockdevice`), use the command  `xfs_quota -cx 'report -h' /home `. 

This will display output similar to the following:

```
User quota on /home (/dev/blockdevice)
                        Blocks              
User ID      Used   Soft   Hard Warn/Grace   
---------- --------------------------------- 
root            0      0      0  00 [------]
testuser   103.4G      0      0  00 [------]
```

To set a soft and hard inode count limit of 500 and 700 respectively for user `john` (whose home directory is `/home/john`) use the following command:

```
xfs_quota -x -c 'limit isoft=500 ihard=700 /home/john'
```

By default, the `limit` sub-command recognizes targets as users. When configuring the limits for a group, use the `-g` option (as in the previous example). Similarly, use `-p` for projects.
Soft and hard block limits can also be configured using `bsoft` or `bhard` instead of `isoft` or `ihard`.

Before configuring limits for project-controlled directories, add them first to `/etc/projects`. Project names can be added to `/etc/projectid` to map project IDs to project names. 

Once a project is added to `/etc/projects`, initialize its project directory using the following command:

```
xfs_quota -c 'project -s projectname'
```

Quotas for projects with initialized directories can then be configured, with:

```
xfs_quota -x -c 'limit -p bsoft=1000m bhard=1200m projectname'
```

Generic quota configuration tools (`quota`, `repquota`, and `edquota` for example) may also be used to manipulate XFS quotas. However, these tools cannot be used with XFS project quotas.

# Increasing the Size of an XFS File System

An XFS file system may be grown while mounted using the `xfs_growfs` command:

(NOTE: If you're using LVM make sure you have [[lvm_notes#add_size_from_disk|grown the underlying logical volume]])

```
xfs_growfs /mnt
```

The `-D size` option grows the file system to the specified size (expressed in file system blocks). 

Without the `-D size` option, `xfs_growfs` will grow the file system to the maximum size supported by the device.

Before growing an XFS file system with `-D size`, ensure that the underlying block device is of an appropriate size to hold the file system later. Use the appropriate resizing methods for the affected block device.

**__Note__** While XFS file systems can be grown while mounted, their size cannot be reduced at all.
# Repairing an XFS File System

To repair an XFS file system, use `xfs_repair`:

```
xfs_repair /dev/sda
```

The `xfs_repair` utility is highly scalable and is designed to repair even very large file systems with many inodes efficiently. Unlike other Linux file systems, `xfs_repair` does not run at boot time, even when an XFS file system was not cleanly unmounted. In the event of an unclean unmount, `xfs_repair` simply replays the log at mount time, ensuring a consistent file system.

**__//WARNING://__**

The `xfs_repair` utility cannot repair an XFS file system with a dirty log. To clear the log, mount and unmount the XFS file system. If the log is corrupt and cannot be replayed, use the `-L` option ("force log zeroing") to clear the log. Be aware that this may result in further corruption or data loss.

```
xfs_repair -L /dev/sda
```

# Suspending an XFS File System

To suspend or resume write activity to a file system, use xfs_freeze. Suspending write activity allows hardware-based device snapshots to be used to capture the file system in a consistent state.

To suspend (that is, freeze) an XFS file system, use:
```
xfs_freeze -f /mnt
```

To unfreeze an XFS file system, use:
```
xfs_freeze -u /mnt
```

When taking an LVM snapshot, it is not necessary to use xfs_freeze to suspend the file system first. Rather, the LVM management tools will automatically suspend the XFS file system before taking the snapshot.

# Backup and Restoration of XFS File Systems

XFS file system backup and restoration involves two utilities: `xfsdump` and `xfsrestore`

To backup or dump an XFS file system, use the `xfsdump` utility. It supports backups to tape drives or regular file images, and also allows multiple dumps to be written to the same tape. The `xfsdump` utility also allows a dump to span multiple tapes, although only one dump can be written to a regular file. In addition, `xfsdump` supports incremental backups, and can exclude files from a backup using size, subtree, or inode flags to filter them.

In order to support incremental backups, `xfsdump` uses dump levels to determine a base dump to which a specific dump is relative. The `-l` option specifies a dump level (0-9). To perform a full backup, perform a level 0 dump on the file system (that is, /path/to/filesystem), as in:

```
xfsdump -l 0 -f /dev/device /path/to/filesystem
```

The `-f` option specifies a destination for a backup. For example, the `/dev/st0` destination is normally used for tape drives. An `xfsdump` destination can be a tape drive, regular file, or remote tape device.

In contrast, an incremental backup will only dump files that changed since the last level 0 dump. A level 1 dump is the first incremental dump after a full dump; the next incremental dump would be level 2, and so on, to a maximum of level 9. So, to perform a level 1 dump to a tape drive:

```
xfsdump -l 1 -f /dev/st0 /path/to/filesystem
```

Conversely, the `xfsrestore` utility restores file systems from dumps produced by xfsdump. The `xfsrestore` utility has two modes: a default simple mode, and a cumulative mode. Specific dumps are identified by session ID or session label. As such, restoring a dump requires its corresponding session ID or label. To display the session ID and labels of all dumps (both full and incremental), use the `-I` option:
```
xfsrestore -I
```

The simple mode allows users to restore an entire file system from a level 0 dump. After identifying a level 0 dump's session ID (that is, session-ID), restore it fully to `/path/to/destination` using:
```
xfsrestore -f /dev/st0 -S session-ID /path/to/destination
```

Note:
The -f option specifies the location of the dump, while the -S or -L option specifies which specific dump to restore. The -S option is used to specify a session ID, while the -L option is used for session labels. The -I option displays both session labels and IDs for each dump.

The `xfsrestore` utility also allows specific files from a dump to be extracted, added, or deleted. To use `xfsrestore` interactively, use the `-i` option, as in:
```
xfsrestore -f /dev/st0 -i
```

The interactive dialogue will begin after `xfsrestore` finishes reading the specified device. Available commands in this dialogue include `cd`, `ls`, `add`, `delete`, and `extract` for a complete list of commands, use `help`

EXAMPLE:

I used the following to "autodump" my filesystem

```
xfsdump -L root_fs -M root_fs -l 0 -f /tmp/root_xfs.dmp /
```

  * `-L` ~ Specifies a label for the dump session.  It can be any arbitrary string up to 255 characters long.
  * `-M` ~ Specifies a label for the first media object (for example, tape cartridge) written on the corresponding destination during the session.  It can be any arbitrary string up to 255 characters long.  Multiple media object labels can be specified, one for each destination.
  * `-l` ~ Level of dump (0 means "full backup")
  * `-f` ~ File to dump to....in this case a file 

# MISC Notes

Other utilities for managing XFS file systems:

**__xfs_fsr__**
Used to defragment mounted XFS file systems. When invoked with no arguments, xfs_fsr defragments all regular files in all mounted XFS file systems. This utility also allows users to suspend a defragmentation at a specified time and resume from where it left off later.
In addition, xfs_fsr also allows the defragmentation of only one file, as in xfs_fsr /path/to/file. Red Hat advises against periodically defragmenting an entire file system, as this is normally not warranted.

**__xfs_bmap__**
Prints the map of disk blocks used by files in an XFS filesystem. This map lists each extent used by a specified file, as well as regions in the file with no corresponding blocks (that is, holes).

**__xfs_info__**
Prints XFS file system information.

**__xfs_admin__**
Changes the parameters of an XFS file system. The xfs_admin utility can only modify parameters of unmounted devices or file systems.

**__xfs_copy__**
Copies the contents of an entire XFS file system to one or more targets in parallel.

The following utilities are also useful in debugging and analyzing XFS file systems:

**__xfs_metadump__**
Copies XFS file system metadata to a file. The xfs_metadump utility should only be used to copy unmounted, read-only, or frozen/suspended file systems; otherwise, generated dumps could be corrupted or inconsistent.

**__xfs_mdrestore__**
Restores an XFS metadump image (generated using xfs_metadump) to a file system image.

**__xfs_db__**
Debugs an XFS file system.

**__Auto Partition__**

Automatically create a partition with `parted`

```
parted --script /dev/vdb mklabel gpt mkpart primary 1MiB 100%
```
