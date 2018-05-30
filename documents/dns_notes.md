# DNS Notes

DNS notes in no paticular order

* [Master DNS](#master-dns)
* [Slave Server](#slave-server)
* [Caching Only Server](#caching-only-server)
* [Forwarding DNS Server](#forwarding-dns-server)
* [Delegating a Sub Domain](#delegating-a-sub-domain)
* [Dynamic DNS Setup](#dynamic-dns-setup)
* [DNS The Easy Way](#dns-the-easy-way)
* [Misc DNS Notes](#misc-dns-notes)

## Master DNS

Make sure you have the proper packages installed

```
root@master1# yum -y install bind-chroot
```

With the "chroot" package; the default path for named is `/var/named/chroot`

Go into `/var/named/chroot` and make a backup of the etc directory...

```
root@master1# cd /var/named/chroot
root@master1# mv etc ,etc
```

Now download the [following tarball](./var_named_chroot_etc.tar) into this directory

```
root@master1# pwd
/varnamed/chroot
root@master1# curl http://dokuwiki-chernandez1982.rhcloud.com/lib/exe/fetch.php?media=var_named_chroot_etc.tar > var_named_chroot.tar
root@master1# tar -xf var_named_chroot.tar
```

Change the ownersip to the `named` user...

```
root@master1# chown -R named:named etc
```

If you're using SELinux make sure you change the context too.. (which is why you made a backup of the etc dir ;) )

```
root@master1# chcon -R --reference=,etc etc
root@master1# semanage fcontext -a -t etc_t "etc(/.*)?"
root@master1# restorecon -R -v etc
```

Make sure you have the most updated `root.hints` file

```
root@master1# pwd
/var/named/chroot/etc
root@master1# curl http://www.internic.net/domain/named.root > root.hints
```

You should be able to (if your permissions are set correctly AND SELinux context are right (if applicable) ) start named

```
root@master1# service named start
```

If you see an error that there is a directory missing...

```
ls: cannot access /var/named/chroot/etc/named: No such file or directory
```

Just create it and set the permissions right...(SELinux if applicable)

```
root@master1# cd /var/named/chroot/etc/
root@master1# mkdir -m 775 named
root@master1# chown named:named named
root@master1# cd /var/named/chroot/
root@master1# semanage fcontext -a -t etc_t "etc(/.*)?"
root@master1# restorecon -R -v etc
```

Now in the tarball you just downloaded; it contains a "basic" setup that gets your up and running. Change the `named.conf`, the `zonemaster/$domainname`, and the `dbmaster/$iprange` file to match the domain you are serving.

Below is sample `named.conf` file

```
# This tells bind where to find the named.conf file and other files.
# the "notify yes" tells bind to notify the "slaves" of any changes.
#
# "options" Directives:
# allow-query    (default is allow all queries)
# allow-transfer (default is allow all transfer requests)
# allow-update   (default is to DENY all update requests)
#

options {
	directory	"/etc";
	dump-file	"/etc/data/cache_dump.db";
	statistics-file	"/etc/data/named_stats.txt";
	notify	yes;
};

# Below are different zones and reverse lookups (in-addr.arpa). The first
# is the root.hints file that tells us where the root servers are so that
# we may quirey them. Then the zone and reverses for our local host.
zone "." in {
	type hint;
	file "root.hints";
};

zone "localhost" in {
	type	master;
	file	"localhost.zone";
};

zone "0.0.127.in-addr.arpa" in {
	type	master;
	file	"127.0.0.db";
};

# Below is our zone and reverse for our domain.
zone "example.org" {
	type master;
	file	"zonemaster/example.org";
};

zone "149.168.192.in-addr.arpa" {
	type	master;
	file	"dbmaster/192.168.149.db";
};

#
#EOF

```

Below is an example of a `zonemaster/$domainname` file

```
$TTL 1W
@	IN	SOA	ns1.example.org.	root (
			2007030200	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.example.org.
	IN	NS	ns2.example.org.
	IN	MX 10	smtp.example.org.
;
; 
	IN	A	192.168.149.10
www	IN	CNAME	example.org.
ns1	IN	A	192.168.149.15
ns2	IN	A	192.168.149.16
smtp	IN	A	192.168.149.5
host	IN	A	192.168.149.25
server	IN	CNAME	host
;
;EOF
```

Example of a `dbmaster/$iprange` file

```
$TTL 1W
@	IN	SOA	ns1.example.org.	root (
			2007030200	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.example.org.
	IN	NS	ns2.example.org.
;
; 
1	IN	PTR	NULL-ENTRY.example.org.
2	IN	PTR	NULL-ENTRY.example.org.
3	IN	PTR	NULL-ENTRY.example.org.
4	IN	PTR	NULL-ENTRY.example.org.
5	IN	PTR	smtp.example.org.
6	IN	PTR	NULL-ENTRY.example.org.
7	IN	PTR	NULL-ENTRY.example.org.
8	IN	PTR	NULL-ENTRY.example.org.
9	IN	PTR	NULL-ENTRY.example.org.
10	IN	PTR	example.org.
11	IN	PTR	NULL-ENTRY.example.org.
12	IN	PTR	NULL-ENTRY.example.org.
13	IN	PTR	NULL-ENTRY.example.org.
14	IN	PTR	NULL-ENTRY.example.org.
15	IN	PTR	ns1.example.org.
16	IN	PTR	ns2.example.org.
17	IN	PTR	NULL-ENTRY.example.org.
18	IN	PTR	NULL-ENTRY.example.org.
19	IN	PTR	NULL-ENTRY.example.org.
20	IN	PTR	NULL-ENTRY.example.org.
21	IN	PTR	NULL-ENTRY.example.org.
22	IN	PTR	NULL-ENTRY.example.org.
23	IN	PTR	NULL-ENTRY.example.org.
24	IN	PTR	NULL-ENTRY.example.org.
25	IN	PTR	host.example.org.
;
;EOF
```

Again, this should all be clear in the tarball file you downloaded (goes without saying...but make sure you change he serial number ;) ).

## Slave Server

To set up the salve server - you basically perform the exact same steps as the Master Configrations. But, you do the following differently.

Since you're a slave. Take out and/or modify the zones to say that you're a slave

```
##################################################################
# This is an example of what a "slave" would look like...this example
# is what the slave would look like on another machine that is slaving
# off of this DNS server. (the 0.0.0.0 part is where the master DNS
# server's IP address would go. The "options" section; along with the 
# "localhost" "0.0.127.in-addr.arpa" and "." stay the same on a 
# slave machine. Only the zones and it's reverses change...look below.
#
zone "example.org" {
	type	slave;
	file	"slave/example.org";
 	masters	{0.0.0.0;};
};


zone "149.168.192.in-addr.arpa" {
	type	slave;
	file	"slave/192.168.149.db";
	masters	{0.0.0.0;};
};

#EOF
```

Then remove the zone files for the masters (since you're a slave)

```
root@slave1# cd /var/named/chroot/etc
root@slave1# rm -f zonemaster/example.org
root@slave1# rm -f dbmaster/192.168.149.db
```

Now you can start named

```
root@slave1# service named start
```

## Caching Only Server

Now if you are not serving a domain but want to serve DNS, you can install a caching only DNS server

* [Doing it with BIND](#bind)
* [Doing it with Unbound](#unbound)

### BIND

Don't install the chroot version

```
root@host# yum -y install bind
```

Now in the `/etc/named.conf` file look for the following lines

```
listen-on port 53 { 127.0.0.1; };
allow-query     { localhost; };
```

And change them to your IP (for the listen-on part) and network's IP range ) for the "allow-query"

```
listen-on port 53 { 127.0.0.1; 172.16.1.102; };
allow-query     { localhost; 172.16.1.0/24; };
```

In the end the file should look like this

```
//
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
	listen-on port 53 { 127.0.0.1; 172.16.1.103; };
	listen-on-v6 port 53 { ::1; };
	directory 	"/var/named";
	dump-file 	"/var/named/data/cache_dump.db";
        statistics-file "/var/named/data/named_stats.txt";
        memstatistics-file "/var/named/data/named_mem_stats.txt";
	allow-query     { localhost; 172.16.1.0/24; };
	recursion yes;

	dnssec-enable yes;
	dnssec-validation yes;
	dnssec-lookaside auto;

	/* Path to ISC DLV key */
	bindkeys-file "/etc/named.iscdlv.key";

	managed-keys-directory "/var/named/dynamic";
};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

zone "." IN {
	type hint;
	file "named.ca";
};

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";

//EOF
```

### Unbound

As root, install the unbound package.

```
[root@sandbox ~]# yum -y install unbound
```

Start and enable the `unbound` service

```
[root@sandbox ~]# systemctl start unbound.service
[root@sandbox ~]# systemctl enable unbound.service
ln -s '/usr/lib/systemd/system/unbound.service' '/etc/systemd/system/multi-user.target.wants/unbound.service'
```

Configure the network interface to listen on.

By default, unbound only listens on the localhost network interface. To make unbound available to remote clients as a caching nameserver, use the interface option in the server clause of `/etc/unbound/unbound.conf` to specify the network interface(s) to listen on. A value of `0.0.0.0` will configure unbound to listen on all network interfaces.

```
[root@sandbox ~]# grep -w interface /etc/unbound/unbound.conf  | egrep -v '#|automatic'
	interface: 0.0.0.0
```

Configure client access.

By default, unbound refuses recursive queries from all clients. In the `server` clause of `/etc/unbound/unbound.conf` use the `access-control` option to specify which clients are allowed to make recursive queries.

```
[root@sandbox ~]# grep access-control /etc/unbound/unbound.conf  | grep -v '#'
	access-control: 172.16.0.0/16 allow
```

Configure forwarding

In `/etc/unbound/unbound.conf`, create a `forward-zone` clause to specify which DNS server(s) to forward queries to. DNS servers can be specified by host name using the `forward-host` option, or by IP address using the forward-addr option. For a caching nameserver, forward all queries by specifying a `forward-zone` of `"."`.

```
[root@sandbox ~]# egrep -A4 '^forward' /etc/unbound/unbound.conf 
forward-zone:
 	name: "."
 	forward-addr: 8.8.8.8
 	forward-addr: 8.8.4.4
```

Bypass DNSSEC validation

By default, unbound is enabled to perform DNSSEC validation to verify all DNS responses received. The domain-insecure option in the server clause of `/etc/unbound/unbound.conf` can be used to specify a domain for which DNSSEC validation should be skipped. This is often desirable when dealing with an unsigned internal domain that would otherwise fail trust chain validation.

```
[root@sandbox ~]# grep domain-insecure /etc/unbound/unbound.conf 
	domain-insecure: "example.net"
```

Check The Conf File

Check the `/etc/unbound/unbound.conf` configuration file for syntax errors.

```
[root@sandbox ~]# unbound-checkconf 
unbound-checkconf: no errors in /etc/unbound/unbound.conf
```

Post Config Steps

Restart the service

```
[root@sandbox ~]# systemctl restart unbound.service
[root@sandbox ~]# systemctl status unbound.service
unbound.service - Unbound recursive Domain Name Server
   Loaded: loaded (/usr/lib/systemd/system/unbound.service; enabled)
   Active: active (running) since Tue 2014-12-23 08:46:21 PST; 10s ago
  Process: 3309 ExecStartPre=/usr/sbin/unbound-checkconf (code=exited, status=0/SUCCESS)
  Process: 3306 ExecStartPre=/sbin/runuser --shell /bin/sh -c /usr/sbin/unbound-anchor -a /var/lib/unbound/root.key -c /etc/unbound/icannbundle.pem unbound (code=exited, status=0/SUCCESS)
 Main PID: 3312 (unbound)
   CGroup: /system.slice/unbound.service
           └─3312 /usr/sbin/unbound -d

Dec 23 08:46:21 sandbox.example.net systemd[1]: Starting Unbound recursive Domain Name Server...
Dec 23 08:46:21 sandbox.example.net runuser[3306]: pam_unix(runuser:session): session opened for user unbound by (uid=0)
Dec 23 08:46:21 sandbox.example.net unbound-checkconf[3309]: unbound-checkconf: no errors in /etc/unbound/unbound.conf
Dec 23 08:46:21 sandbox.example.net systemd[1]: Started Unbound recursive Domain Name Server.
Dec 23 08:46:21 sandbox.example.net unbound[3312]: Dec 23 08:46:21 unbound[3312:0] warning: increased limit(open files) from 1024 to 8266
Dec 23 08:46:21 sandbox.example.net unbound[3312]: [3312:0] notice: init module 0: validator
Dec 23 08:46:21 sandbox.example.net unbound[3312]: [3312:0] notice: init module 1: iterator
Dec 23 08:46:21 sandbox.example.net unbound[3312]: [3312:0] info: start of service (unbound 1.4.20).
```

Add firewall rules

```
[root@sandbox ~]# firewall-cmd --permanent --add-service=dns
success
[root@sandbox ~]# firewall-cmd --reload
success
```

Test on client

```
[user@client ~]$ dnsip=172.16.1.2
[user@client ~]$ dig @${dnsip} www.google.com 

; <<>> DiG 9.9.4-RedHat-9.9.4-14.el7_0.1 <<>> @172.16.1.2 www.google.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 364
;; flags: qr rd ra; QUERY: 1, ANSWER: 6, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;www.google.com.			IN	A

;; ANSWER SECTION:
www.google.com.		89	IN	A	74.125.28.106
www.google.com.		89	IN	A	74.125.28.104
www.google.com.		89	IN	A	74.125.28.103
www.google.com.		89	IN	A	74.125.28.105
www.google.com.		89	IN	A	74.125.28.99
www.google.com.		89	IN	A	74.125.28.147

;; Query time: 245 msec
;; SERVER: 172.16.1.2#53(172.16.1.2)
;; WHEN: Tue Dec 23 08:52:52 PST 2014
;; MSG SIZE  rcvd: 139

```

Dumping Cache

Administrators of caching nameservers need to dump out cache data when troubleshooting DNS issues, such as those resulting from stale resource records. With an unbound DNS server, the cache can be dumped by running the `unbound-control` utility in conjunction with the `dump_cache` subcommand.

```
[root@sandbox ~]# unbound-control dump_cache > /tmp/unbound.cache
[root@sandbox ~]# less /tmp/unbound.cache
```

To load cache from file...

```
[root@sandbox ~]# unbound-control load_cache < /tmp/unbound.cache 
ok
```

Flushing Cache

Administrators of caching nameservers also need to purge outdated resource records from cache from time to time.

```
[root@sandbox ~]# unbound-control flush www.example.com
ok
```

If all resource records belonging to a domain need to be purged from the cache of an unbound DNS server, unbound-control can be executed with the `flush_zone` subcommand.

```
[root@sandbox ~]# unbound-control flush_zone example.com
ok removed 3 rrsets, 1 messages and 0 key entries
```

## Forwarding DNS Server

The Steps for a forwarding DNS server is the same as the "caching" DNS server; except you add the following two lines in the "options" directive.

```
forward only;
forwarders { 192.168.2.225; 192.168.2.223; };
```

Replace the IP with the IP you're forwarding your DNS query to.

The file should look something like this.

```
//
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
	listen-on port 53 { 127.0.0.1; 172.16.1.103; };
	listen-on-v6 port 53 { ::1; };
	directory 	"/var/named";
        forward only;
        forwarders { 192.168.2.225; 192.168.2.223; };
	dump-file 	"/var/named/data/cache_dump.db";
        statistics-file "/var/named/data/named_stats.txt";
        memstatistics-file "/var/named/data/named_mem_stats.txt";
	allow-query     { localhost; 172.16.1.0/24; };
	recursion yes;

	dnssec-enable yes;
	dnssec-validation yes;
	dnssec-lookaside auto;

	/* Path to ISC DLV key */
	bindkeys-file "/etc/named.iscdlv.key";

	managed-keys-directory "/var/named/dynamic";
};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

zone "." IN {
	type hint;
	file "named.ca";
};

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";

//EOF
```

If the server you are forwarding to doesn't do "secure" transactions; you can disable it (and only use "insecure" communications) by changing the following line

```
dnssec-validation yes;
```

To this

```
dnssec-validation no;
```

## Delegating a Sub Domain

If you are serving `example.com` but want to "delegate" the DNS services for `lax.example.com` - you specify that in the zonefile (in the above example that's the `zonemaster/$domainname` file).

REMEMBER: BIND is "persinckity" about whitespace. Make sure you use the right number of tabs

```
; Delegation below
ns1.lax          IN     A     10.1.2.3
ns2.lax          IN     A     10.1.2.4
lax              IN     NS    ns1.lax
                 IN     NS    ns2.lax
;
;                 
```

## Dynamic DNS Setup

You can setup bind to edit the zone files "automagically" with a key. These are QnD notes

Create the key that will be used to update named

```
domain=test.tld
keyfile=/var/named/$domain.key
rm -vf /var/named/K$domain*
cd /var/named
dnssec-keygen -a HMAC-SHA256 -b 256 -n USER -r /dev/urandom $domain
KEY="$(grep Key: /var/named/K$domain*.private | cut -d ' ' -f 2)"
rndc-confgen -a -r /dev/urandom
```

Once this is done you can verify by echoing the values. Also let’s make sure that SELinux and permissions are right.
```
echo $domain
echo $keyfile
echo $KEY
restorecon -v /etc/rndc.* /etc/named.*
chown -v root:named /etc/rndc.key
chmod -v 640 /etc/rndc.key
restorecon -vR /var/named
```

Now create the zonefile
```
rm -rvf /var/named/dynamic
mkdir -vp /var/named/dynamic
cat <<EOF > /var/named/dynamic/${domain}.db
\$ORIGIN .
\$TTL 1	; 1 seconds (for testing only)
${domain}               IN SOA  ns1.${domain}. hostmaster.${domain}. (
                                2011112904 ; serial
                                60         ; refresh (1 minute)
                                15         ; retry (15 seconds)
                                1800       ; expire (30 minutes)
                                10         ; minimum (10 seconds)
                                )
                        NS      ns1.${domain}.
                        MX      10 mail.${domain}.
\$ORIGIN ${domain}.
ns1                     A       172.16.1.220
EOF
```
Now create the keyfile
```
cat <<EOF > /var/named/$domain.key
key $domain {
 algorithm HMAC-SHA256;
 secret "${KEY}";
};
EOF
```
Let’s clean up the directory again making sure permissions and SELinux is right.
```
chgrp named -R /var/named
chown named -R /var/named/dynamic
restorecon -vR /var/named
```
Verify your files
```
cat /var/named/$domain.key
cat /var/named/dynamic/${domain}.db
```
Create the // /etc/named.conf // file for dynamic DNS. 
```
cat <<EOF > /etc/named.conf
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
  listen-on port 53 { any; };
  directory 	"/var/named";
  dump-file 	"/var/named/data/cache_dump.db";
        statistics-file "/var/named/data/named_stats.txt";
        memstatistics-file "/var/named/data/named_mem_stats.txt";
  allow-query     { any; };
  recursion no;

  /* Path to ISC DLV key */
  bindkeys-file "/etc/named.iscdlv.key";

};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

// use the default rndc key
include "/etc/rndc.key";

controls {
  inet 172.16.1.220 port 953
  allow { 172.16.1.0/24; } keys { "rndc-key"; };
};

include "/etc/named.rfc1912.zones";

include "$domain.key";

zone "$domain" IN {
  type master;
  file "dynamic/$domain.db";
  allow-update { key $domain; };
};
EOF

```

Open up port 953 for tcp/udp on both nodes

```
firewall-cmd --permanent --add-port=953/tcp --add-port=953/udp
firewall-cmd --reload
```

Set the permissions and SELinux context to the correct values:

```
chown -v root:named /etc/named.conf
restorecon /etc/named.conf
```

Start the service
```
systemctl start named
systemctl enable named
```

Test making an entry
```
scp $keyfile user@otherhost:/tmp/
ssh user@otherhost
dig @172.16.1.220 ns1.test.tld +short
nsupdate -k /tmp/test.tld.key 
> server 172.16.1.220
> update add foobar.test.tld 180 A 172.16.1.220
> send
> [Control+D]
dig @172.16.1.220 foobar.test.tld +short
```

## DNS The Easy Way

Easiest way is to **NOT** use `bind-chroot` and to use just `bind` in the stead. It's less secure but it's okay for home networks ...I tested this on Fedora 27

```
dnf -y install bind bind-utils
```

Edit the `etc/named.conf` file to look similar to what's below. I made only 3 changes
  * `listen-on` was changed to read `any`
  * `allow-query` was changed to read `any`
  * The entries for my zonefiles (noted by comment)

```
//
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
	listen-on port 53 { any; };
	listen-on-v6 port 53 { ::1; };
	directory 	"/var/named";
	dump-file 	"/var/named/data/cache_dump.db";
	statistics-file "/var/named/data/named_stats.txt";
	memstatistics-file "/var/named/data/named_mem_stats.txt";
	allow-query     { any; };

	/* 
	 - If you are building an AUTHORITATIVE DNS server, do NOT enable recursion.
	 - If you are building a RECURSIVE (caching) DNS server, you need to enable 
	   recursion. 
	 - If your recursive DNS server has a public IP address, you MUST enable access 
	   control to limit queries to your legitimate users. Failing to do so will
	   cause your server to become part of large scale DNS amplification 
	   attacks. Implementing BCP38 within your network would greatly
	   reduce such attack surface 
	*/
	recursion yes;

	dnssec-enable yes;
	dnssec-validation yes;

	managed-keys-directory "/var/named/dynamic";

	pid-file "/run/named/named.pid";
	session-keyfile "/run/named/session.key";

	/* https://fedoraproject.org/wiki/Changes/CryptoPolicy */
	include "/etc/crypto-policies/back-ends/bind.config";
};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

zone "." IN {
	type hint;
	file "named.ca";
};

########### Add what's between these comments ###########
zone "cloud.chx" IN {
	type master;
	file	"cloud.chx";
};

zone "1.16.172.in-addr.arpa" IN {
	type	master;
	file	"172.16.1.db";
};
########################################################

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
```

Place your files in `var/named`

```
# cat /var/named/cloud.chx 
$TTL 1W
@	IN	SOA	ns1.cloud.chx.	root (
			2018030100	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.cloud.chx.
	IN	MX 10	smtp.cloud.chx.
;
; 
	IN	A	172.16.1.254
www	IN	CNAME	cloud.chx.
ns1	IN	A	172.16.1.254
smtp	IN	A	172.16.1.254
;
;EOF
```

For reverse...

```
# cat /var/named/172.16.1.db 
$TTL 1W
@	IN	SOA	ns1.cloud.chx.	root (
			2018030101	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.cloud.chx.
;
; 
1	IN	PTR	NULL-ENTRY.cloud.chx.
2	IN	PTR	NULL-ENTRY.cloud.chx.
3	IN	PTR	NULL-ENTRY.cloud.chx.
4	IN	PTR	NULL-ENTRY.cloud.chx.
5	IN	PTR	smtp.cloud.chx.
6	IN	PTR	NULL-ENTRY.cloud.chx.
7	IN	PTR	NULL-ENTRY.cloud.chx.
8	IN	PTR	NULL-ENTRY.cloud.chx.
9	IN	PTR	NULL-ENTRY.cloud.chx.
10	IN	PTR	cloud.chx.
11	IN	PTR	NULL-ENTRY.cloud.chx.
12	IN	PTR	NULL-ENTRY.cloud.chx.
13	IN	PTR	NULL-ENTRY.cloud.chx.
14	IN	PTR	NULL-ENTRY.cloud.chx.
15	IN	PTR	ns1.cloud.chx.
16	IN	PTR	ns2.cloud.chx.
17	IN	PTR	NULL-ENTRY.cloud.chx.
18	IN	PTR	NULL-ENTRY.cloud.chx.
19	IN	PTR	NULL-ENTRY.cloud.chx.
20	IN	PTR	NULL-ENTRY.cloud.chx.
21	IN	PTR	NULL-ENTRY.cloud.chx.
22	IN	PTR	NULL-ENTRY.cloud.chx.
23	IN	PTR	NULL-ENTRY.cloud.chx.
24	IN	PTR	NULL-ENTRY.cloud.chx.
25	IN	PTR	host.cloud.chx.
;
254	IN	PTR	ns1.cloud.chx.
;
;EOF
```

Make sure SELinux is happy
```
restorecon -vR /var/named
```

Open up the firewall if that's your thing
```
firewall-cmd --permanent --add-port=53/tcp --add-port=53/udp
firewall-cmd --add-port=53/tcp --add-port=53/udp
```

Enable it on boot

```
# systemctl enable named
Created symlink /etc/systemd/system/multi-user.target.wants/named.service → /usr/lib/systemd/system/named.service.

# systemctl start named
```

Test

```
# dig @localhost ns1.cloud.chx +short
172.16.1.254

# dig @localhost -x 172.16.1.254 +short
ns1.cloud.chx.

```

## Misc DNS Notes

Zone Transfer

If the DNS server allows it...you can view (and therefore download) their entire zone file.

```
root@host# dig @server domain axfr
```
