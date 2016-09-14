# GRUB 2

First See your grub entries

	
	root@host# grep menuentry /boot/grub2/grub.cfg


Make note of what "position" the entry is. GRUB2 is 0 based. So if it's the 3rd entry; it'd be counted as 2.

Now open the // /etc/default/grub // file and make a change...

	
	GRUB_DEFAULT=2


Also while I was here I edited the file to turn off apic (//apic=off//) in the // GRUB_CMDLINE_LINUX= // line.

Now after you made the changes; you need to recreate the grub file.

	
	root@host# grub2-mkconfig -o /boot/grub2/grub.cfg


__Misc Notes__

I ran into an issue after doing a "yum -y upgrade" on Fedora 18 where grub didn't "see" my Windows partition (and subsequently grub2-mkconfig gave errors when trying to rebuild the grub.conf file). It seems that the script // /etc/grub.d/30_os-prober // wasn't able to detect Windows when running the // grub2-mkconfig -o /boot/grub2/grub.cfg // command.

The errors looks something like this (it spit out quite a few of these lines).

	
	ERROR: ddf1: seeking device "/dev/dm-6" to 18446744073709421056


I found something online (in an Ubuntu forum out of all places) that solved my issue.

First create a script called // /etc/grub.d/31_windows-probe // that looked like this

	
	[root@fedora18 ~]# cat /etc/grub.d/31_windows-probe
	#!/bin/sh -e
	#
	cat << EOF
	menuentry "Microsoft Windows 7" {
	  set root=(hd0,1)
	  chainloader +1
	}
	EOF
	#
	##-30-


Then I made it executable

	
	[root@fedora18 ~]# chmod +x /etc/grub.d/31_windows-probe


Now I was able to run the // grub2-mkconfig -o /boot/grub2/grub.cfg // and was able to see Windows 7 in the grub menu
