Firewall rules needed for communication (not sure if NEEDED...need to test this without firewall rules)
	lokkit -f -v --enabled --service=ssh --port=389:tcp --port=389:udp --port=636:tcp --port=636:udp --port=88:tcp --port=88:udp --port=749:tcp --port=749:udp

		*NOTE: This might "look like" it failed...but if you run an "iptables --list" you'll see that the rules got implamented

# Via NSLCD


Make sure "legacy" mode is on for authentication
	grep LEGACY /etc/sysconfig/authconfig
	     FORCELEGACY=yes


Create /etc/nslcd.conf
	uri ldap://sm-flor-xdc19.wdw.disney.com
	base ou=pam,dc=wdw,dc=disney,dc=com
	##binddn CN=9hernanc,OU=Users,OU=PAM,DC=wdw,DC=disney,DC=com
	binddn CN=Hernandez\, Christian,OU=Users,OU=PAM,DC=wdw,DC=disney,DC=com
	bindpw Buzz#999
	# Mappings for Active Directory
	pagesize 1000
	referrals off
	filter passwd (&(objectClass=user)(!(objectClass=computer))(uidNumber=*)(unixHomeDirectory=*))
	map    passwd uid              sAMAccountName
	map    passwd homeDirectory    unixHomeDirectory
	map    passwd gecos            displayName
	filter shadow (&(objectClass=user)(!(objectClass=computer))(uidNumber=*)(unixHomeDirectory=*))
	map    shadow uid              sAMAccountName
	map    shadow shadowLastChange pwdLastSet
	filter group  (objectClass=group)
	map    group  uniqueMember     member
	ssl no
	# tls_cacertdir /etc/openldap/cacerts

NOTE: If setting up pam_ldap - use the /etc/ldap.conf file instead

	
	host gc.swna.wdpr.disney.com:3268
	base ou=pam,dc=wdw,dc=disney,dc=com
	binddn CN=Hernandez\, Christian,OU=Users,OU=PAM,DC=wdw,DC=disney,DC=com
	bindpw Buzz#999
	timelimit 120
	bind_timelimit 120
	idle_timelimit 3600
	
	nss_initgroups_ignoreusers root,ldap,named,avahi,haldaemon,dbus,radvd,tomcat,radiusd,news,mailman,nscd,gdm
	
	referrals no
	nss_schema rfc2307bis 
	
	nss_base_passwd ou=pam,dc=wdw,dc=disney,dc=com?sub?&(objectCategory=user)(uidnumber=*)
	nss_base_shadow ou=pam,dc=wdw,dc=disney,dc=com?sub?&(objectCategory=user)(uidnumber=*)
	nss_base_group  ou=pam,dc=wdw,dc=disney,dc=com?sub?&(objectCategory=group)(gidnumber=*)
	nss_map_objectclass posixAccount user
	nss_map_objectclass shadowAccount user
	nss_map_objectclass posixGroup group
	nss_map_attribute homeDirectory unixHomeDirectory
	nss_map_attribute shadowLastChange pwdLastSet
	nss_map_attribute gecos displayName
	nss_map_attribute uniqueMember member
	pam_member_attribute member
	pam_login_attribute sAMAccountName
	
	pam_password ad
	ssl no
	


Run authconfig-tui
	User Information Section: Select "Use LDAP
	Authentication Section: Select "User Kerberos"
	
	LDAP Settings:
		Leave "TLS" unchecked
		Server: ldap://sm-flor-xdc19.wdw.disney.com
		Base DN: ou=pam,dc=wdw,dc=disney,dc=com

	Kerberos Settings:
		Realm: WDW.DISNEY.COM
		KDC: 153.6.4.126
		Admin Server: 153.6.4.126

Make sure homedirs are created
		 authconfig  --enablemkhomedir --update

Enable ACLs for access restrictions
		authconfig --enablepamaccess --update

Create custom PAM rules for logining in
		/etc/pam.d/login
			#%PAM-1.0
			auth [user_unknown=ignore success=ok ignore=ignore default=bad] pam_securetty.so
			auth       include      system-auth
			account    required     pam_nologin.so
			account    include      system-auth
			password   include      system-auth
			# pam_selinux.so close should be the first session rule
			session    required     pam_selinux.so close
			session    required     pam_loginuid.so
			session    optional     pam_console.so
			# pam_selinux.so open should only be followed by sessions to be executed in the user context
			session    required     pam_selinux.so open
			session    required     pam_namespace.so
			session    optional     pam_keyinit.so force revoke
			session    include      system-auth
			-session   optional     pam_ck_connector.so
			# Custom Configs
			auth		sufficient	pam_ldap.so
			account		sufficient	pam_ldap.so
			password	sufficient	pam_ldap.so
			session		sufficient	pam_ldap.so
			#
			##EOF

Create custom PAM rules for SSH
		/etc/pam.d/sshd
			#%PAM-1.0
			auth       include      system-auth
			account    required     pam_nologin.so
			account    include      system-auth
			password   include      system-auth
			session    optional     pam_keyinit.so force revoke
			session    include      system-auth
			session    required     pam_loginuid.so
			auth       requisite    pam_deny.so
			# Custom Configs :: Christian Hernandez 2013-MAR-28
			auth            sufficient      pam_ldap.so
			account         sufficient      pam_ldap.so
			password        sufficient      pam_ldap.so
			session         sufficient      pam_ldap.so
			#
			##EOF

# SSSD Notes

SSSD Config required some changes (this is a supplement to the above)...


The /etc/sysconfig/authconfig file needed the following changes
	USESSSD=yes
	USESSSDAUTH=yes
	FORCELEGACY=no
	CACHECREDENTIALS=yes

In the end...it looked something like this...
	[root@DI-FLOR-LBD394d ~]# cat /etc/sysconfig/authconfig 
	  IPADOMAINJOINED=no
	  USEMKHOMEDIR=yes
	  USEPAMACCESS=yes
	  CACHECREDENTIALS=yes
	  USESSSDAUTH=yes
	  USESHADOW=yes
	  USEWINBIND=no
	  USESSSD=yes
	  USEDB=no
	  FORCELEGACY=no
	  USEFPRINTD=no
	  FORCESMARTCARD=no
	  USELDAPAUTH=no
	  USEPASSWDQC=no
	  IPAV2NONTP=no
	  USELOCAUTHORIZE=yes
	  USECRACKLIB=yes
	  USEIPAV2=no
	  USEWINBINDAUTH=no
	  USESMARTCARD=no
	  USELDAP=yes
	  USENIS=no
	  USEKERBEROS=yes
	  USESYSNETAUTH=no
	  PASSWDALGORITHM=sha512
	  USEHESIOD=no

NOTE: Maybe I could have gotten this by running...not sure didn't test
	root@host# authconfig --disableforcelegacy --enablesssd --enablesssdauth --update

Then I edited the /etc/sssd/sssd.conf file to look something like this
	[root@DI-FLOR-LBD394d ~]# cat /etc/sssd/sssd.conf
	  [sssd]
	  config_file_version = 2
	  reconnection_retries = 3
	  sbus_timeout = 30
	  services = nss, pam
	  domains = default
	  
	  [nss]
	  filter_groups = root
	  filter_users = root
	  reconnection_retries = 3
	  
	  [pam]
	  reconnection_retries = 3
	  
	  [domain/default]
	  id_provider = ldap
	  chpass_provider = krb5
	  
	  ldap_uri = ldap://sm-flor-xdc19.wdw.disney.com
	  ldap_search_base = ou=pam,dc=wdw,dc=disney,dc=com
	  
	  ldap_id_use_start_tls = False
	  
	  ldap_default_bind_dn = CN=Hernandez\, Christian,OU=Users,OU=PAM,DC=wdw,DC=disney,DC=com
	  ldap_default_authtok_type = password
	  ldap_default_authtok = Buzz#999
	  
	  ldap_schema = rfc2307bis
	  ldap_force_upper_case_realm = True
	  ldap_user_object_class = person
	  ldap_group_object_class = group
	  ldap_user_principal = userPrincipalName
	  ldap_user_fullname = displayName
	  ldap_user_name = sAMAccountName
	  ldap_user_object_class = user
	  ldap_user_home_directory = unixHomeDirectory
	  ldap_user_shell = loginShell
	  ldap_user_principal = userPrincipalName
	  ldap_force_upper_case_realm = True
	  
	  auth_provider = krb5
	  krb5_server = 153.6.4.126
	  krb5_realm = WDW.DISNEY.COM
	  cache_credentials = True
	  krb5_kpasswd = 153.6.4.126
	  ldap_user_gecos = displayName
	  debug_level = 0
	  ldap_tls_cacertdir = /etc/openldap/cacerts

Run authconfig-tui
	User Information Section: Select "Use LDAP
	Authentication Section: Select "User Kerberos"
	
	LDAP Settings:
		Leave "TLS" unchecked
		Server: ldap://sm-flor-xdc19.wdw.disney.com
		Base DN: ou=pam,dc=wdw,dc=disney,dc=com

	Kerberos Settings:
		Realm: WDW.DISNEY.COM
		KDC: 153.6.4.126
		Admin Server: 153.6.4.126

You can get the same thing with this command...

__RHEL 5__

	
	root@host# authconfig --enablesssd --enablesssdauth --enablelocauthorize --enableldap --ldapserver="ldap://sm-flor-xdc19.wdw.disney.com" --ldapbasedn="ou=pam,dc=wdw,dc=disney,dc=com" --disableldapssl --enablekrb5 --krb5realm=WDW.DISNEY.COM --krb5kdc=153.6.4.126 --krb5adminserver=153.6.4.126 --enablemkhomedir --enablepamaccess --update


__RHEL 6__

	
	root@host# authconfig --disableforcelegacy --enablesssd --enablesssdauth --enablelocauthorize --enableldap --ldapserver="ldap://sm-flor-xdc19.wdw.disney.com" --ldapbasedn="ou=pam,dc=wdw,dc=disney,dc=com" --disableldaptls --enablekrb5 --krb5realm=WDW.DISNEY.COM --krb5kdc=153.6.4.126 --krb5adminserver=153.6.4.126 --enablemkhomedir --enablepamaccess --update

