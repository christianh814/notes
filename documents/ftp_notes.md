# FTP Notes

FTP Notes in no paticular order

* [Chrooted FTP](#chrooted-ftp)
* [Linux TFPT Set Up](#linux-tfpt-set-up)

## Chrooted FTP

Install vsFTP

``
root@host# yum -y install vsftpd
root@host# chkconfig vsftpd on
``

Add this to the `/etc/vsftpd/vsftpd.conf`

```
chroot_local_user=YES
```

Also; make sure the following is set as well

```
anonymous_enable=NO
local_enable=YES
write_enable=YES
```

Make sure the homedir is set to where you want them to be chrooted to (Make sure the shell is set to `/bin/true`) ...

```
ftpuser:x:401:401::/var/www/html/./:/bin/true
```

If `/bin/true` isn't in the `/etc/shells` you can add it with the following...

```
root@host# echo "/bin/true" >> /etc/shells
```

Restart the service

```
root@host# service vsftpd restart
```

You'll probably need to check SELinux settings (aka `getsebool -a` or `semanage fcontext`)

## Linux TFPT Set Up

This is how you set up tftp.

Make sure that you have tftp installed ... you can install it via yum

```
root@host# yum -y install tftp
```

To change where you want the root tftp directory edit the `/etc/xinetd.d/tftp` file.

```
service tftp
{
	socket_type             = dgram
	protocol                = udp
	wait                    = yes
	user                    = root
	server                  = /usr/sbin/in.tftpd
	server_args             = -s /tftpboot
	disable                 = no
	per_source              = 11
	cps                     = 100 2
	flags                   = IPv4
}
```

Remember - in TFTP the file must ALREADY exist with wide open permissions. You can't upload "new" files; the file must already exist.

```
root@host# touch /tftpboot/myfile.txt
root@host# chmod 777 /tftpboot/myfile.txt
```

-30-
