# TCP WRAPPERS

TCP Wrappers is another way to secure your system. These are controlled through
the /etc/hosts.allow /etc/hosts.deny files

NOTE: It works like the cron.{allow,deny} files.
 1.  If it's in the /etc/hosts.allow file; it allows action
 2.  If it's not; it checks the /etc/hosts.deny file
 3.  If it's in niether of the; the default is to allow access

To find out if the application uses tcpwrappers run the following ldd command
        root@host# ldd /usr/sbin/sshd | grep -i wrap
          libwrap.so.0 => /lib64/libwrap.so.0 (0x00007f8763874)

Basic config is this (NOTE: tcp wrappers doesn't use CIDR notation)
        root@host# vi /etc/hosts.allow
          sshd: 192.168.1.0/255.255.255.0

Comma delimited values are accepted as well
        root@host# vi /etc/hosts.allow
          sshd: 192.168.1.0/255.255.255.0,192.168.2.0/255.255.255.0
          in.tftpd, vsftpd: 192.168.100.0/255.255.255.0

