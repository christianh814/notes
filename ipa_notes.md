Documentation found [here](https///docs.fedoraproject.org/en-US/) \\
FreeIPA Site [here](http://freeipa.org)

# Server Notes

Things assumed
 

*  Server Name: ipa.example.com

*  IP  Address: 172.16.1.10/24

*  iptables disabled

*  selinux disabled

You need to ensure that your // /etc/hosts // file is configured correctly. A misconfigured file can prevent the IPA command-line tools from functioning correctly and can prevent the IPA web interface from connecting to the IPA server.

Configure the // /etc/hosts // file to list the FQDN for the IPA server before any aliases. Also ensure that the hostname is not part of the localhost entry. 

The following is an example of a valid hosts file

	
	127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
	::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
	172.16.1.10	ipa.example.com	ipa


The IPA packages are now a part of the CentOS repos...if you're setting up for RHEL you may need to enable the proper repos.

Install what you need along with IPA (The package *bind* is optional if you are also making IPA your DNS server - recommended but not required)

	
	root@server# yum -y install ipa-server bind bind-dyndb-ldap


The IPA setup script creates a server instance, which includes configuring all of the required services for the IPA domain:


*  The network time daemon (ntpd)

*  A 389 Directory Server instance

*  A Kerberos key distribution center (KDC)

*  Apache (httpd)

*  An updated SELinux targeted policy

*  The Active Directory WinSync plug-in

*  A certificate authority

*  Optional. A domain name service (DNS) server 

The IPA setup process can be minimal, where the administrator only supplies some required information, or it can be very specific, with user-defined settings for many parts of the IPA services. The configuration is passed using arguments with the *ipa-install-server* script. 

Invoke the script. If DNS is properly set up - the defaults are okay to use (the // --setup-dns // option is if you want to set up IPA as a DNS server)

	
	root@server# ipa-server-install --setup-dns 


The output is something like this

	
	
	The log file for this installation can be found in /var/log/ipaserver-install.log
	==============================================================================
	This program will set up the IPA Server.
	
	This includes:

	  * Configure a stand-alone CA (dogtag) for certificate management
	  * Configure the Network Time Daemon (ntpd)
	  * Create and configure an instance of Directory Server
	  * Create and configure a Kerberos Key Distribution Center (KDC)
	  * Configure Apache (httpd)
	  * Configure DNS (bind)
	
	To accept the default shown in brackets, press the Enter key.
	
	Existing BIND configuration detected, overwrite? [no]: yes
	Enter the fully qualified domain name of the computer
	on which you're setting up server software. Using the form
	`<hostname>`.`<domainname>`
	Example: master.example.com.
	
	
	Server host name [ipa.example.com]: 
	
	Warning: skipping DNS resolution of host ipa.example.com
	The domain name has been calculated based on the host name.
	
	Please confirm the domain name [example.com]: 
	
	The kerberos protocol requires a Realm name to be defined.
	This is typically the domain name converted to uppercase.
	
	Please provide a realm name [EXAMPLE.COM]: 
	Certain directory server operations require an administrative user.
	This user is referred to as the Directory Manager and has full access
	to the Directory for system management tasks and will be added to the
	instance of directory server created for IPA.
	The password must be at least 8 characters long.
	
	Directory Manager password: 
	Password (confirm): 
	
	The IPA server requires an administrative user, named 'admin'.
	This user is a regular system account used for IPA server administration.
	
	IPA admin password: 
	Password (confirm): 
	
	Do you want to configure DNS forwarders? [yes]: 
	Enter the IP address of DNS forwarder to use, or press Enter to finish.
	Enter IP address for a DNS forwarder: 10.0.2.3
	DNS forwarder 10.0.2.3 added
	Enter IP address for a DNS forwarder: 
	Do you want to configure the reverse zone? [yes]: 
	Please specify the reverse zone name [1.16.172.in-addr.arpa.]: 
	Using reverse zone 1.16.172.in-addr.arpa.
	
	The IPA Master Server will be configured with:
	Hostname:      ipa.example.com
	IP address:    172.16.1.10
	Domain name:   example.com
	Realm name:    EXAMPLE.COM
	
	BIND DNS server will be configured to serve IPA domain with:
	Forwarders:    10.0.2.3
	Reverse zone:  1.16.172.in-addr.arpa.
	
	Continue to configure the system with these values? [no]: yes
	
	The following operations may take some minutes to complete.
	Please wait until the prompt is returned.
	
	Configuring ntpd...done configuring ntpd.
	Configuring directory server for the CA: Estimated time 30 seconds...done configuring pkids.
	Configuring certificate server: Estimated time 3 minutes 30 seconds...done configuring pki-cad.
	Configuring directory server: Estimated time 1 minute...done configuring dirsrv.
	Configuring Kerberos KDC: Estimated time 30 seconds...done configuring krb5kdc.
	Configuring kadmin...done configuring kadmin.
	Configuring ipa_memcached...done configuring ipa_memcached.
	Configuring the web interface: Estimated time 1 minute...done configuring httpd.
	Applying LDAP updates
	Restarting the directory server
	Restarting the KDC
	Configuring named...done configuring named.
	
	Global DNS configuration in LDAP server is empty
	You can use 'dnsconfig-mod' command to set global DNS options that
	would override settings in local named.conf files
	
	Restarting the web server
	==============================================================================
	Setup complete
	
	Next steps:
		1. You must make sure these network ports are open:
			TCP Ports:

			  * 80, 443: HTTP/HTTPS
			  * 389, 636: LDAP/LDAPS
			  * 88, 464: kerberos
			  * 53: bind
			UDP Ports:

			  * 88, 464: kerberos
			  * 53: bind
			  * 123: ntp
	
		2. You can now obtain a kerberos ticket using the command: 'kinit admin'
		   This ticket will allow you to use the IPA tools (e.g., ipa user-add)
		   and the web user interface.
	
	Be sure to back up the CA certificate stored in /root/cacert.p12
	This file is required to create replicas. The password for this
	file is the Directory Manager password
	


After installation reboot the server

Login and test the IPA server

	
	root@server# kinit admin
	Password for admin@EXAMPLE.COM:


Then do a basic search

	
	root@server# ipa user-find admin
	  --------------
	  1 user matched
	  --------------
	  User login: admin
	  Last name: Administrator
	  Home directory: /home/admin
	  Login shell: /bin/bash
	  Account disabled: False
	  Member of groups: admins
	  ----------------------------
	  Number of entries returned 1
	  ----------------------------


Note that kerberos server has very specific DNS requirements, if you have a DNS server already on your network add the SRV records of the kerberos, ntp and ldap server to that. A sample zone file will be created in your // /tmp // directory after the // ipa-server-install // , do a copy paste of all the SRV record from this file to your zone file. Make sure you add Forward and Reverse DNS names for ALL servers and services in the IPA domain.

I should look something like this

	
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;;              THIS IS NECISSARY FOR THE IPA SERVER            ;;
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;
	; ldap servers
	_ldap._tcp              IN SRV 0 100 389        ipa1
	
	;kerberos realm
	_kerberos               IN TXT 4OVER.COM
	
	; kerberos servers
	_kerberos._tcp          IN SRV 0 100 88         ipa1
	_kerberos._udp          IN SRV 0 100 88         ipa1
	_kerberos-master._tcp   IN SRV 0 100 88         ipa1
	_kerberos-master._udp   IN SRV 0 100 88         ipa1
	_kpasswd._tcp           IN SRV 0 100 464        ipa1
	_kpasswd._udp           IN SRV 0 100 464        ipa1
	
	;ntp server
	_ntp._udp               IN SRV 0 100 123        ipa1.gln.4over.com.
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;
	;-30-

# Server Replicas

Replicas are functionally the same as FreeIPA servers, so they have the same installation requirements and packages. The version of the OS and IPA must be the same on both the server and replica.

Install the server packages first (Do **//__not__//** run the // ipa-server-install // script on the replica)

	
	root@replica# yum -y install ipa-server


On the master server, create a replica information file. This contains realm and configuration information taken from the master server which will be used to configure the replica server.

	
	root@server# ipa-replica-prepare ipa1.la3.4over.com
	Directory Manager (existing master) password: 
	
	Preparing replica for ipa1.la3.4over.com from ipa1.gln.4over.com
	Creating SSL certificate for the Directory Server
	Creating SSL certificate for the dogtag Directory Server
	Creating SSL certificate for the Web Server
	Exporting RA certificate
	Copying additional files
	Finalizing configuration
	Packaging replica information into /var/lib/ipa/replica-info-ipa1.la3.4over.com.gpg


If you are serving DNS on another server make sure you have the ***SAME*** entries for each replica in the SRV records. Should look like this

	
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;;              THIS IS NECISSARY FOR THE IPA SERVER            ;;
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;
	; ldap servers
	_ldap._tcp              IN SRV 0 100 389        server1
	_ldap._tcp              IN SRV 0 100 389        server2
	
	;kerberos realm
	_kerberos               IN TXT 4OVER.COM
	
	; kerberos servers
	_kerberos._tcp          IN SRV 0 100 88         server1
	_kerberos._udp          IN SRV 0 100 88         server1
	_kerberos-master._tcp   IN SRV 0 100 88         server1
	_kerberos-master._udp   IN SRV 0 100 88         server1
	_kpasswd._tcp           IN SRV 0 100 464        server1
	_kpasswd._udp           IN SRV 0 100 464        server1
	
	; kerberos servers la3
	_kerberos._tcp          IN SRV 0 100 88         server2
	_kerberos._udp          IN SRV 0 100 88         server2
	_kerberos-master._tcp   IN SRV 0 100 88         server2
	_kerberos-master._udp   IN SRV 0 100 88         server2
	_kpasswd._tcp           IN SRV 0 100 464        server2
	_kpasswd._udp           IN SRV 0 100 464        server2
	
	;ntp server
	_ntp._udp               IN SRV 0 100 123        ipa1.gln.4over.com.
	;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
	;
	;-30-


NOTE: Using the // --ip-address // option automatically creates DNS entries for the replica, including the A and PTR records for the replica to the DNS (if you are using IPA as a DNS server).

	
	root@server# ipa-replica-prepare ipa1.la3.4over.com --ip-address


Each replica information file is created in the // /var/lib/ipa/ // directory as a GPG-encrypted file. Each file is named specifically for the replica server for which it is intended, such as // replica-info-ipareplica.example.com.gpg //.

Copy this replica information file to the replica server

	
	root@server# scp /var/lib/ipa/replica-info-ipareplica.example.com.gpg root@ipareplica:/var/lib/ipa/


Now on the replica server, run the replica installation script, referencing the replication information file. There are other options for setting up DNS, much like the server installation script. Additionally, there is an option to configure a CA for the replica; while CA's are installed by default for servers, they are optional for replicas. You should probably set up a "like for like" copy.

	
	root@replica# ipa-replica-install --setup-ca --no-ntp /var/lib/ipa/replica-info-ipa1.la3.4over.com.gpg 


If you set up with DNS then you need to specify

	
	root@replica# ipa-replica-install --setup-ca --setup-dns /var/lib/ipa/replica-info-ipa1.la3.4over.com.gpg 


GOTCHAS: Make sure the "admin" user can ssh into the servers in the // /etc/security/access.conf // file

Now you should be able to see all your IPA servers

	
	root@server# ipa-replica-manage list
	ipa1.la3.4over.com: master
	ipa1.gln.4over.com: master 

__Failover__

Failover is done by DNS...you add both of your servers as a SRV record. This is done at the application level. The clients/servers querey DNS for it's next available server.

Load balancing is done automatically by servers, replicas, and clients. The configuration used for load balancing can be altered by changing the priority and the weight given to a server or replica.

If doing DNS on another server...the entries look something like this

	
	; ldap servers
	_ldap._tcp              IN SRV 0 100 389        server1
	_ldap._tcp              IN SRV 1 100 389        server2


In the above example...server1 is weighted; this gives server1 a higher priority than server 2, meaning it will be contacted first

If you are using IPA as a DNS server you can create an entry using the cmd line tools

	
	root@server# ipa dnsrecord-add server.example.com _ldap._tcp --srv-rec="0 100 389 server1.example.com." 
	root@server# ipa dnsrecord-add server.example.com _ldap._tcp --srv-rec="1 100 389 server2.example.com."


__Forcing Replication__

Replication between servers and replicas occurs on a schedule. Although replication is frequent, there can be times when it is necessary to initiate the replication operation manually. For example, if a server is being taken offline for maintenance, it is necessary to flush all of the queued replication changes out of its changelog before taking it down.
To initiate a replication update manually, use the force-sync command. The server which receives the update is the local server; the server which sends the updates is specified in the // --from//  option.

	
	root@host# ipa-replica-manage force-sync --from srv1.example.com


__Reinitializing FreeIPA Servers__

When a replica is first created, the database of the master server is copied, completely, over to the replica database. This process is called initialization. If a server/replica is offline for a long period of time or there is some kind of corruption in its database, then the server can be re-initialized, with a fresh and updated set of data.
This is done using the re-initialize command. The target server being initialized is the local host. The server or replica from which to pull the data to initialize the local database is specified in the --from option:

	
	root@host# ipa-replica-manage re-initialize --from srv1.example.com


__Forrest Configuration__

When you set up multiple multi-master replica...you have to set up subsequent replication agreements manually.

A ~> B \\
B ~> A \\
A ~> C \\
C ~> A \\

BUT we need to set up

B ~> C \\
C ~> B \\

At 4over I added a 3rd replication server. So my GLN list looked like this

	
	root@GLN# ipa-replica-manage list ipa1.gln.4over.com
	ipa1.da2.4over.com: replica
	ipa1.la3.4over.com: replica


BUT When I looked at LA3...I only saw GLN not the new DA2 server...

	
	root@LA3# ipa-replica-manage list ipa1.la3.4over.com
	ipa1.gln.4over.com: replica


We have to set up the agreement between LA3 and DA2 now (example below done with encryption keys)...

First you'll see only one agreement per server...

	
	root@LA3# # ipa-replica-manage list ipa1.la3.4over.com
	ipa1.gln.4over.com: replica
	
	root@DA2# # ipa-replica-manage list ipa1.da2.4over.com
	ipa1.gln.4over.com: replica



Add the two servers with the 'connect' command (syntax is // I want to have a replication agreement between C and B // )

	
	root@LA3# # ipa-replica-manage connect --cacert=/etc/ipa/ca.crt ipa1.la3.4over.com ipa1.da2.4over.com
	Added CA certificate /etc/ipa/ca.crt to certificate database for ipa1.la3.4over.com
	Connected 'ipa1.la3.4over.com' to 'ipa1.da2.4over.com'


Now all servers should come up in the list

	
	root@LA3# # ipa-replica-manage list ipa1.la3.4over.com
	ipa1.da2.4over.com: replica
	ipa1.gln.4over.com: replica
	
	root@DA2# # ipa-replica-manage list ipa1.da2.4over.com
	ipa1.gln.4over.com: replica
	ipa1.la3.4over.com: replica

# Client Notes

__Install Client__

First install the client Binaries

	
	root@client# yum -y install ipa-client


Then run command (The --enable-dns-updates is only if you're using your IPA server as a DNS server)

	
	root@client# ipa-client-install --enable-dns-updates --mkhomedir
	Discovery was successful!
	Hostname: server1.example.com
	Realm: EXAMPLE.COM
	DNS Domain: example.com
	IPA Server: ipa.example.com
	BaseDN: dc=example,dc=com
	
	
	Continue to configure the system with these values? [no]: yes
	User authorized to enroll computers: admin
	Synchronizing time with KDC...
	Password for admin@EXAMPLE.COM: 
	
	Enrolled in IPA realm EXAMPLE.COM
	Created /etc/ipa/default.conf
	Domain example.com is already configured in existing SSSD config, creating a new one.
	The old /etc/sssd/sssd.conf is backed up and will be restored during uninstall.
	Configured /etc/sssd/sssd.conf
	Configured /etc/krb5.conf for IPA realm EXAMPLE.COM
	DNS server record set to: server1.example.com -> 172.16.1.11
	SSSD enabled
	NTP enabled
	Client configuration complete.


__Manual Adding__

First you must add Forward and Reverse DNS on the IPA server ( notes [below](https///dokuwiki-chernandez1982.rhcloud.com/doku.php?id=ipa_notes&#misc_notes) )

Next on the IPA server create a host entry for the new server

	
	root@server# ipa host-add server2.example.com --password secret


Run the command on the client

	
	root@client# hostname > /tmp/hostname.txt
	root@client# ipa-client-install --domain=EXAMPLE.COM --enable-dns-updates --mkhomedir --password=secret --realm=EXAMPLE.COM --server=ipa.example.com --unattended


__Remove Client__

Run the uninstall command (reboot required)

	
	root@client# ipa-client-install --uninstall



Then on the server remove the old host from the IPS DNS domain. While this is optional, it cleans up the old IPA entries associated with the system and allows it to be re-enrolled cleanly at a later time. 

	
	root@server# ipa host-del server1.example.com


__Home Dirs__

While PAM modules can be used to create home directories for users automatically, this may not be desirable behavior in every environment. In that case, home directories can be manually added to the IPA server from separate locations using NFS shares and automount.

Create a new location for the user directory maps

	
	root@server# ipa automountlocation-add userdirs
	Location: userdirs


Add a direct map to the new location's auto.direct file. In this example, the mount point is // /share //

	
	root@server# ipa automountkey-add userdirs auto.direct --key=/share --info="-ro,soft, ipaserver.example.com:/home/share"
	
	Key: /share
	Mount information: -ro,soft, ipaserver.example.com:/home/share


There are two ways to enable the pam_oddjob_mkhomedir (or pam_mkhomedir) module:

 1.  The // --mkhomedir // option can be used with the // ipa-client-install // command. While this is possible for clients, this option is not available to servers when they are set up.
 2.  The // pam_oddjob_mkhomedir // module can be enabled using the system's // authconfig // command. For example:

`authconfig --enablemkhomedir --update`

This option can be used for both server and client machines post-installation. 
# MISC Notes

__Migrate An Existing LDAP Server to IPA__

A migration tool **ipa migrate-ds** is provided - Documentation is [here](https///access.redhat.com/knowledge/docs/en-US/Red_Hat_Enterprise_Linux/6/html/Identity_Management_Guide/Migrating_from_a_Directory_Server_to_IPA.html)

First you have to enable migration.

	
	root@server# ipa config-mod --enable-migration=true


Command I used at 4over

	
	root@host# ipa migrate-ds --with-compat --base-dn="dc=4over,dc=com"  --user-container="dc=4over,dc=com" --group-container="dc=4over,dc=com" ldap://ldap-la3.4over.com 


Then you can prompt your users to use the following link to reset their passwords

	
	https://ipa.example.com/ipa/migration/


__Add DNS Record__

You can add a DNS record from the Web GUI or from the command line with

	
	root@server@ ipa dnsrecord-add example.com www --a-rec 10.64.14.165


PTR Record

	
	root@server# ipa dnsrecord-add 1.16.172.in-addr.arpa. 99 --ptr-rec=server2.example.com.


CNAME Record

	
	root@server# ipa dnsrecord-add example.org nfs --cname-rec=ipa.example.org.


Add Reverse Zone

	
	root@server# ipa dnszone-add 206.65.10.in-addr.arpa.

OR (adding Reverse Zone)

	
	root@server# ipa dnszone-add 10.65.206.0/24

__Remove DNS Record__

	
	root@server# ipa dnsrecord-del example.com www --a-rec 10.64.14.213


__Delegate Subdomain__

	
	root@host# ipa dnsrecord-add example.com ns1.cloud  --a-rec=192.168.3.170
	root@host# ipa dnsrecord-add example.com cloud  --ns-rec=ns1.subdomain.example.com


__named Daemon Fails to Start__

If a IPA server is configured to manage DNS and is set up successfully, but the named service fails to start, this can indicate that there is a package conflict. Check the // /var/log/messages // file for error messages related to the named service and the // ldap.so // library

	
	ipaserver named[6886]: failed to dynamically load driver 'ldap.so': libldap-2.4.so.2: cannot open shared object file: No such file or directory


This usually means that the bind-chroot package is installed and is preventing the named service from starting. To resolve this issue, remove the bind-chroot package and then restart the FreeIPA server. 

	
	root@server# yum -y remove bind-chroot


__Adding User__

	
	root@server# ipa user-add chrish --first=Christian --last=Hernandez --gecos="Christian Hernandez"  --email=christianh@example.com --homedir=/home/christianh --password --uid=637 --gidnumber=637 
	Password: 
	Enter Password again to verify: 
	-------------------
	Added user "chrish"
	-------------------
	  User login: chrish
	  First name: Christian
	  Last name: Hernandez
	  Full name: Christian Hernandez
	  Display name: Christian Hernandez
	  Initials: CH
	  Home directory: /home/christianh
	  GECOS field: Christian Hernandez
	  Login shell: /bin/sh
	  Kerberos principal: chrish@EXAMPLE.COM
	  Email address: christianh@example.com
	  UID: 637
	  GID: 637
	  Password: True
	  Kerberos keys available: True
	  

    
__Delete User__

Deleting user is a little easier

	
	root@server#  ipa user-del chrish


__Changing Password__

 Changing a password — your own or another user's — is done using the user-mod command, as with other user account changes.

	
	root@server# kinit admin
	root@server# ipa user-mod chrish --password


__Add Groups__

	
	root@server# ipa group-add webbie --gid=504 --desc="Web Masters"


__Add Users To Groups__

Members are added to a group using the group-add-member command. This command can add both users as group members and other groups as group members.
The syntax of the group-add-member command requires only the group name and a comma-separated list of users to add

	
	root@server# ipa group-add-member webbie --users=chrish,donw


__Services__

On the IPA server

	
	root@server# ipa service-add HTTP/montools2.4over.com


Remove a service

	
	root@server# ipa service-disable HTTP/montools2.4over.com
	
	root@server# ipa service-del HTTP/montools2.4over.com


List services

	
	root@server# ipa service-find --pkey-only


List specific service

	
	root@server# ipa service-show HTTP/montools2.4over.com


__Key Tabs__

Getting a keytab is easiest on the client. 

	
	root@client# kinit admin
	root@client# ipa-getkeytab -s ipa1.gln.4over.com -p HTTP/montools2.4over.com -k /etc/httpd/conf/krb5.keytab


__Defaults Change__

By default the shell for users is set to // /bin/sh // change this to bash by running...

	
	root@server# ipa config-mod --defaultshell=/bin/bash


At 4over, I also changed the default home dir

	
	root@server# ipa config-mod --homedirectory=/rhome


Now you can see the defaults change...

	
	root@server# ipa config-show --all
	dn: cn=ipaconfig,cn=etc,dc=4over,dc=com
	  Maximum username length: 32
	  Home directory base: /rhome
	  Default shell: /bin/bash
	  Default users group: ipausers
	  Default e-mail domain: 4over.com
	  Search time limit: 2
	  Search size limit: 100
	  User search fields: uid,givenname,sn,telephonenumber,ou,title
	  Group search fields: cn,description
	  Enable migration mode: TRUE
	  Certificate Subject base: O=4OVER.COM
	  Default group objectclasses: top, groupofnames, nestedgroup, ipausergroup, ipaobject
	  Default user objectclasses: top, person, organizationalperson, inetorgperson, inetuser, posixaccount, krbprincipalaux, krbticketpolicyaux, ipaobject, ipasshuser
	  Password Expiration Notification (days): 4
	  Password plugin features: AllowNThash
	  SELinux user map order: guest_u:s0$xguest_u:s0$user_u:s0-s0:c0.c1023$staff_u:s0-s0:c0.c1023$unconfined_u:s0-s0:c0.c1023
	  Default SELinux user: guest_u:s0
	  cn: ipaConfig
	  objectclass: nsContainer, top, ipaGuiConfig, ipaConfigObject


__Password Policy__

List current policy

	
	root@server# ipa pwpolicy-show


List policy by policy group name

	
	root@server# ipa pwpolicy-show policygroupName


List policy by user

	
	root@server# ipa pwpolicy-show --user=chrish


The default password age policy sets it for 90 days...I changed it at 4over to something that won't annoy anyone (Syntax it's in days)

	
	root@server# ipa pwpolicy-mod --maxlife=1830


Changing the password expiration time in the password policy does not affect the expiration date for a user, until the user password is changed. If the password expiration date needs to be changed immediately, it can be changed by editing the user entry.

BUT - For some reason the following command doesn't work (or is not yet supported)

	
	root@server# ipa user-mod chrish --setattr=krbPasswordExpiration=20121231011529Z


I had to directly modify LDAP...I used an LDIF file for this...

Fist the LDIF file

	
	dn: uid=cvs,cn=users,cn=accounts,dc=4over,dc=com
	changetype: modify
	replace: krbpasswordexpiration
	krbpasswordexpiration: 20180101000059Z
	-


The // krbpasswordexpiration// format is YYYYMMDDHHMMSSZ 

Now I can "read" this file into the // ldapmodify // command to make modifications. (below is the command I used at 4over)

	
	root@server# ldapmodify -x -h ipa1.gln.4over.com -p 389 -D 'cn=Directory Manager' -w `<password>` -vv -f filename.ldiff


__Red Hat IPA server w/ Fedora 18 Client__

I needed to "manually" add the server...the auto config of adding hosts didn't work. I also had to disable SSH configuration

	
	root@host# ipa-client-install --no-dns-sshfp --no-ssh


__NSLCD Config__

If you are going to configure IPA without the client app (i.e. Straight LDAP with NSLCD). You have to make sure that your config file looks something like this.

	
	[chrish@chrish.sbx.4over.com ~]$ sudo cat /etc/nslcd.conf 
	timelimit 10
	bind_timelimit 10
	idle_timelimit 10
	uid nslcd
	gid ldap
	uri ldap://ipa1.gln.4over.com/ ldap://ipa1.la3.4over.com.ipa1.da2.4over.com/
	base cn=accounts,dc=4over,dc=com
	map group uniqueMember member
	ssl no
	tls_cacertdir /etc/openldap/cacerts


This should fix the problem of the group mappings not working correctly.

__Allow AXFR__

To allow localhost to do a zone transfer...

	
	ipa dnszone-mod example.com --allow-transfer='localhost;'


Then you can do this...

	
	dig @localhost example.com axfr


This allows you to see the "raw" dns zonefile


