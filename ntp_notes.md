# EL6

Below is EL6 specific (might work on EL7 but ymmv)

## Server Set Up

This is pretty straight forward. The default package is enough and the conf file is well commented.

Basic Steps

*  Install ntp

*  Edit // /etc/ntp.conf //

*  Open ports 123 tcp/udp


Install NTP package

	
	root@host# yum -y install ntp


In the  // /etc/ntp.conf // file (around line 18) uncomment the following line and modify to your network range

	
	restrict 192.168.1.0 mask 255.255.255.0 nomodify notrap


Make sure the service starts on boot

	
	root@host# service ntpd start
	root@host# chkconfig ntpd on


Open ports 123 on TCP/UDP using...

	
	root@host# system-config-firewall

## Client Set Up

Simple client set up using Pacific Time

	
	root@host# yum -y install tzdata ntp
	root@host# rm /etc/localtime
	root@host# ln -s /usr/share/zoneinfo/US/Pacific /etc/localtime
	root@host# echo "server timehost.4over.com" >> /etc/ntp.conf
	root@host# chkconfig ntpd on
	root@host# service ntpd restart
	root@host# ntpq -p


# EL7

Below is EL7 specific

## Client Set Up

The *timedatectl* command shows an overview of the current time-related system settings, including current time, time zone, and NTP synchronization settings of the system.

	
	[student@host ~]$ timedatectl
	      Local time: Thu 2014-02-13 02:16:15 EST
	  Universal time: Thu 2014-02-13 07:16:15 UTC
	        RTC time: Thu 2014-02-13 07:16:15
	        Timezone: America/New_York (EST, -0500)
	     NTP enabled: yes
	NTP synchronized: no
	 RTC in local TZ: no
	      DST active: no
	 Last DST change: DST ended at
	                  Sun 2013-11-03 01:59:59 EDT
	                  Sun 2013-11-03 01:00:00 EST
	 Next DST change: DST begins (the clock jumps one hour forward) at
	                  Sun 2014-03-09 01:59:59 EST
	                  Sun 2014-03-09 03:00:00 EDT

A database with known time zones is available and can be listed with:

	
	[root@host ~]# timedatectl list-timezones
	Africa/Abidjan
	Africa/Accra
	Africa/Addis_Ababa
	Africa/Algiers
	Africa/Asmara
	Africa/Bamako
	...


The system setting for the current time zone can be adjusted as user root:

	
	[root@host ~]# timedatectl set-timezone America/Los_Angeles
	[root@host ~]# timedatectl
	      Local time: Thu 2014-02-13 00:23:54 MST
	  Universal time: Thu 2014-02-13 07:23:54 UTC
	        RTC time: Thu 2014-02-13 07:23:53
	        Timezone: America/Phoenix (MST, -0700)
	     NTP enabled: yes
	NTP synchronized: no
	 RTC in local TZ: no
	      DST active: n/a

To change the current time and date settings with the *timedatectl* command, the set-time option is available. The time is specified in the "YYYY-MM-DD hh:mm:ss" format, where either date or time can be omitted. To change the time to 09:00:00, run:

	
	[root@host ~]# timedatectl set-time 9:00:00
	[root@host ~]# timedatectl
	      Local time: Thu 2014-02-13 09:00:27 MST
	  Universal time: Thu 2014-02-13 16:00:27 UTC
	        RTC time: Thu 2014-02-13 16:00:28
	        Timezone: America/Phoenix (MST, -0700)
	     NTP enabled: yes
	NTP synchronized: no
	 RTC in local TZ: no
	      DST active: n/a

The set-ntp option enables or disables NTP synchronization for automatic time adjustment. The option requires either a true or false argument to turn it on or off. To turn on NTP synchronization, run:

`[root@host ~]# timedatectl set-ntp true`
