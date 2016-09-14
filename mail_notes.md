# Installation

Install with yum

	
	root@host# yum -y install postfix


Make sure it starts on boot and that there are no conflicts

	
	root@host# chkconfig postfix on
	root@host# chkconfig sendmail off

Let's stop conflicts while we are at it

	
	root@host# service sendmail stop


Make sure that what you want is the default

	
	root@host# alternatives --set mta /usr/sbin/sendmail.postfix


# Telnet Port 25

To test an email server you can telnet into port 25 and “ask it” if it accepts mail from a domain. (commands in bold for easy reading)

	gorilla:/<16># telnet smtp.4over.com 25
	Trying 192.168.11.150...
	Connected to smtp.4over.com.
	Escape character is '^]'.
	220 exchange-gln.4over.local Microsoft ESMTP MAIL Service, Version: 6.0.3790.1830 ready at  Thu, 15 May 2008 16:33:05 -0700
	HELO 4over.com
	250 exchange-gln.4over.local Hello [192.168.11.42]
	MAIL from:sysman@4over.com
	250 2.1.0 sysman@4over.com....Sender OK
	RCPT to:chrish@4over.com
	250 2.1.5 chrish@4over.com 
	quit
	221 2.0.0 exchange-gln.4over.local Service closing transmission channel
	Connection to smtp.4over.com closed by foreign host.

You can use this to trouble shoot issues with email


## Postfix Configuration

__Basic How To__ 

Edit the // /etc/postfix/main.cf // file with a text editor, such as vi.


*  Uncomment the *mydomain* line by removing the hash mark (#), and replace domain.tld with the domain the mail server is servicing, such as example.com.

*  Uncomment the *myorigin = $mydomain * line.

*  Uncomment the *myhostname* line, and replace *host.domain.tld* with the hostname for the machine.

*  Uncomment the *mydestination = $myhostname, localhost.$mydomain * line.

*  Uncomment the *mynetworks* line, and replace *168.100.189.0/28* with a valid network setting for hosts that can connect to the server.

*  Uncomment the *inet_interfaces = all* line.

*  Edit the // /etc/postfix/access // file and add IP and OK to things you want to let relay (example: 192.168 OK) run // postmap /etc/postfix/access // when done


__More Detailed how to__

Postfix has a method to update it's database (files that has no .cf extention) that you must update with
the "postmap" command (similar to the makemap command). The format of these files are "key      value"  
So lets say you made a change to the "generic" file you will run the following to update the database
        root@host# postmap generic 

Edit the aliases file to allow/deny networks to use your service 
        root@host# vi /etc/postfix/access
          192.168.122   OK
          192.168.100   REJECT  

Make sure that the database reads the new configuration
        root@host# postmap /etc/postfix/access

Postfix is setup to accept localmail; you need to edit the main.cf file. By default, Postfix does not accept network connections from any host other than the local host. (the file is HEAVLY comented; so I am only specifying what I changed)
        root@host# vi /etc/postfix/main.cf
          myhostname = server1.example.com
          mydomain = example.com
          myorigin = $mydomain
          inet_interfaces = all
          mydestination = $mydomain, $myhostname, localhost.$mydomain, localhost (<~~ added "$mydomain")
          mynetworks = 192.168.122.0/24, 127.0.0.0/8
        
NOTE: It's pretty easy to forget "one minor" thing...the things you want to change was found by searcing for ^#my to find what you need (as well as ^#inet_)


Once changes are made the "postconf" and "postfix check" commands are helpful to test the configurations
        root@host# postconf -n
        root@host# postfix check


Make sure the service starts on boot 
        root@host# chkconfig postfix on

Postfix authentication directives can be found under /usr/share/doc/postfix-VERSION/README-postfix-SASL-RedHat.txt
It's good "how to"...so good in fact that you will use it on the test. So just remember where it is :)

To set up a relay you have to uncoment/edit the "relayhost" directive
        root@host# vi /etc/postfix/main.cf
          relayhost = oustsider1.example.org
          
Check your config

	
	root@host# service postfix check


Restart the service

	
	root@host# service postfix restart


Once these steps are complete, the host accepts outside emails for delivery.

Postfix has a large assortment of configuration options. One of the best ways to learn how to configure Postfix is to read the comments within /etc/postfix/main.cf. Additional resources including information about LDAP and SpamAssassin integration are available online at http://www.postfix.org/. 

## Client Configuration


For clients; you just need to setup the "relay host"

	
	root@host# postconf -e 'relayhost = smtp2.4over.com'


Check your configuration

	
	root@host# service postfix check


Restart the service

	
	root@host# service postfix restart


## Misc Commands

Print default settings

	
	root@host# postconf -d 


Print parameter settings that are not left at their built-in default value, because they are explicitly specified in // main.cf //

	
	root@host# # postconf  -n
	alias_database = hash:/etc/aliases
	alias_maps = hash:/etc/aliases
	broken_sasl_auth_clients = yes
	command_directory = /usr/sbin
	config_directory = /etc/postfix
	daemon_directory = /usr/libexec/postfix
	data_directory = /var/lib/postfix
	debug_peer_level = 2
	html_directory = no
	inet_interfaces = all
	inet_protocols = all
	local_header_rewrite_clients = static:all
	mail_owner = postfix
	mailq_path = /usr/bin/mailq.postfix
	manpage_directory = /usr/share/man
	masquerade_domains = $mydomain
	message_size_limit = 20480000
	mydomain = 4over.com
	myhostname = smtp2.4over.com
	mynetworks = 192.168.0.0/16, 10.0.0.0/8, 127.0.0.0/8
	myorigin = $mydomain
	newaliases_path = /usr/bin/newaliases.postfix
	queue_directory = /var/spool/postfix
	readme_directory = /usr/share/doc/postfix-2.6.6/README_FILES
	relayhost = bizmail.sbcglobal.net
	sample_directory = /usr/share/doc/postfix-2.6.6/samples
	sendmail_path = /usr/sbin/sendmail.postfix
	setgid_group = postdrop
	smtp_data_xfer_timeout = 600s
	smtp_sasl_security_options = 
	smtpd_client_connection_count_limit = 150
	smtpd_recipient_restrictions = permit_sasl_authenticated, permit_mynetworks, reject_unauth_destination
	smtpd_sasl_auth_enable = yes
	smtpd_sasl_security_options = noanonymous
	unknown_local_recipient_reject_code = 550



