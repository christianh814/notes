# Require Password

To passwrod protect page you must first create a password DB file

	
	htpasswd [-c] /data/tools_4over/etc/passwd username

The // -c // means "create" the passwd file, // /data/tools_4over/etc/passwd // represents the passwd file for that web instance and “username” is the username you’re creating/updating. You will be prompted for a password, twice.

To assign the password to a particular directory, see the following example:

	
	# 28-APR-07  Add password protection to goals page; lo/pw = xxxxx/yyyyy
	`<Directory /var/apache/htdocs/intranet/news/goals>`
	  AuthType Basic
	  AuthName "4over Priviledged Information"
	  AuthUserFile /var/apache/etc/goals_passwd
	  Require valid-user
	`</Directory>`



To allow a particular directory visibility, even though it would otherwise require authorization ...

	
	# Special permissions -- allow access w/o validation ...
	`<Directory /data/djwalters/www/images>`
	  Order allow,deny
	  Allow from all
	  Options Indexes 
	  Satisfy any
	`</Directory>`


Sample Using LDAP instead of "flat file"

	
	`<VirtualHost *:80>`
	 `<Directory /data/dokuwiki/www/>`
	    Options FollowSymLinks
	    AllowOverride None
	    Order deny,allow
	    Allow from all
	    Deny from all
	  `</Directory>`
	  `<Directory /data/dokuwiki/www/>`
	    AuthType Basic
	    AuthName "4over Priviledged Information"
	    AuthBasicProvider ldap
	    AuthzLDAPAuthoritative on
	    AuthLDAPURL ldap://ldap-gln1.4over.com:389/dc=4over,dc=com?uid
	    AuthLDAPBindDN "cn=directory manager"
	    AuthLDAPBindPassword minus273
	    Require valid-user
	  `</Directory>`
	    ServerAdmin webmaster@4over.com
	    DocumentRoot /data/dokuwiki/www/
	    ServerName dokuwiki.4over.com
	    ErrorLog logs/dokuwiki.4over.com-error_log
	    CustomLog logs/dokuwiki.4over.com-access_log common
	    ScriptAlias /cgi-bin/ "/data/dokuwiki/cgi-bin/"
	`</VirtualHost>`


***__NOTE:__*** You can use BOTH *flat file* and *ldap*...info [here](openshift_enterprise_2.x#ldap)

__Sample Using Kerberos__

This is using IPA...may be different if using "straight" kerberos

First create a policy on the IPA server

	
	root@ipaserver# ipa service-add HTTP/montools2.4over.com


Next on the web server grab the keytab file (make sure mod_auth_kerb is installed)

	
	root@web# yum -y install mod_auth_kerb mod_authz_ldap
	root@web# ipa-getkeytab -s ipa1.gln.4over.com -p HTTP/montools2.4over.com -k /etc/httpd/conf/krb5.keytab


Now the config looks something like this

	
	`<VirtualHost *:80>`
	 `<Directory /var/www/html/montools2>`
	    Options FollowSymLinks
	    AllowOverride None
	    Order deny,allow
	    Allow from all
	    Deny from all
	  `</Directory>`
	  `<Directory /var/www/html/montools2>`
	    Require SSL
	    AuthType Kerberos
	    AuthName "4over Priviledged Information"
	    KrbServiceName HTTP
	    Krb5KeyTab /etc/httpd/conf/krb5.keytab
	    KrbAuthRealms 4OVER.COM
	    Require valid-user
	  `</Directory>`
	    ServerAdmin webmaster@4over.com
	    DocumentRoot /var/www/html/montools2
	    ServerName montools2.4over.com
	    ErrorLog logs/montools2.4over.com-error_log
	    CustomLog logs/montools2.4over.com-access_log common
	    #ScriptAlias /cgi-bin/ "/data/dokuwiki/cgi-bin/"
	`</VirtualHost>`


For more information on IPA have a look [here](ipa_notes)
# Rewrite Rules

Sample re-write to another site

	
	    ###################################################################
	    #   	 	REWRITE TO GODADDY		              #
	    ###################################################################
	    RewriteEngine on
	    RewriteCond	%{HTTP_HOST}	^webmail.fundink.com$	[NC]
	    RewriteRule	^(.*)$	https://login.secureserver.net/$1	[R=301,L]
	    ###################################################################


Rewrite "anything" to site (basically if you want to redirect www.site.com to myblog.blogsite.com for example)

	
	    ###################################################################
	    #   	 	REWRITE TO GODADDY		              #
	    ###################################################################
	    RewriteEngine on
	    RewriteCond %{HTTP_HOST}             ^(.*)$                   [NC]
	    RewriteRule ^(.*)$  http://myblog.blogsite.com/$1       [R=301,L]
	    ###################################################################


Another (better?) way?

	
	    RewriteEngine On
	    RewriteCond %{HTTPS} off
	    RewriteRule (.*) https://%{HTTP_HOST}%{REQUEST_URI}

# Vitual Hosts

Sample “Standard” Virtual Host Entry

    `<VirtualHost 192.168.11.146:80 69.237.62.146:80>`
    ServerAdmin donw@4over.com
    DocumentRoot /data/pro_4over/www
    ServerName pro.4over.com
    ErrorLog /var/apache/logs/pro.4over.com-error_log
    CustomLog /var/apache/logs/pro.4over.com-access_log common
    ScriptAlias /cgi-bin/ "/data/pro_4over/cgi-bin/"
    `</VirtualHost>`

Sample Virtual Host Entry With Directory Access

    `<VirtualHost 192.168.11.146:80 69.237.62.146:80>`
    `<Directory /data/pro_4over/www>`
    Options FollowSymLinks
    AllowOverride None
    Order deny,allow
    Allow from all
    Deny from all
    `</Directory>`
    ServerAdmin donw@4over.com
    DocumentRoot /data/pro_4over/www
    ServerName pro.4over.com
    ErrorLog /var/apache/logs/pro.4over.com-error_log
    CustomLog /var/apache/logs/pro.4over.com-access_log common
    ScriptAlias /cgi-bin/ "/data/pro_4over/cgi-bin/"
    `</VirtualHost>`

Sample Virtual Host Entry With Directory Access (CGI-BIN)

With the new version of Apache (2.2 and up) You need a “`<Directory>`” entry. Example below (in bold)

    `<VirtualHost 192.168.11.146:80 69.237.62.146:80>`
    `<Directory /data/pro_4over/www>`
    Options FollowSymLinks
    AllowOverride None
    Order deny,allow
    Allow from all
    Deny from all
    `</Directory>`
    ServerAdmin donw@4over.com
    DocumentRoot /data/pro_4over/www
    ServerName pro.4over.com
    ServerAlias www.pro.4over.com
    ErrorLog /var/apache/logs/pro.4over.com-error_log
    CustomLog /var/apache/logs/pro.4over.com-access_log common
    `<Directory /data/pro_4over/cgi-bin/>`
    AllowOverride None
    Options None
    Order allow,deny
    Allow from all
   `</Directory>`
    ScriptAlias /cgi-bin/ "/data/pro_4over/cgi-bin/"
    `</VirtualHost>`

# Big5 Apache Logging (X Forwarding)


In the httpd.conf file add in the following entries (changes in bold)

	<IfModule log_config_module>
	  LogFormat "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\"" combined
	  LogFormat "%h %l %u %t \"%r\" %>s %b" common
	  LogFormat "%v %{X-Forwarded-For}i %l %u %t \"%r\" %>s %b" X-Forwarded-For
	<IfModule logio_module>
	  LogFormat "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\" %I %O" combinedio
	</IfModule>
	  CustomLog "/web/logs/access_log" common
	  CustomLog "/web/logs/acces_log-xforwarded.log" X-Forwarded-For
	</IfModule>

After making the changes restart apache

	
	root@host# service httpd restart


# SSL


After adding an SSL certificate to apache; it will prompt you for the “ssl passphrase” when you start apache. In order to make it not prompt you for the passphrase (necessary for it to startup on boot). You must edit the httpd.conf file and the httpd-ssl.conf file

Create a script that performs an “echo” of the passphrase. Remember to make it executable (chmod +x script.sh).

	#!/bin/sh
	echo "passphrase"

Edit the httpd.conf file (httpd.conf) and make sure that this include statement is there:

	# Secure (SSL/TLS) connections
	Include conf/extra/httpd-ssl.conf

In the httpd-ssl.conf file (/usr/local/apache2/conf/extra/httpd-ssl.conf); and add where the script is:

	#   Pass Phrase Dialog:
	#   Configure the pass phrase gathering process.
	#   The filtering dialog program (`builtin' is a internal
	#   terminal dialog) has to provide the pass phrase on stdout.
	SSLPassPhraseDialog  exec:/usr/local/apache2/ssl/sslpw.sh

You might also need a vhost file...something similar to this

	
	#
	# Bidpit
	#
	SSLPassPhraseDialog  exec:/etc/httpd/ssl/sslpw.sh
	LoadModule ssl_module modules/mod_ssl.so
	Listen 443
	`<VirtualHost *:443>`
	 `<Directory /data/bidpit>`
	    Options FollowSymLinks
	    AllowOverride All
	    Order deny,allow
	    Allow from all
	    Deny from all
	  `</Directory>`
	    SSLEngine on
	    SSLCipherSuite ALL:!ADH:!EXPORT56:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP:+eNULL
	    SSLCertificateFile /etc/httpd/ssl/bidpit.com.crt
	    SSLCertificateKeyFile /etc/httpd/ssl/wildbidpit.key
	    SSLCertificateChainFile /etc/httpd/ssl/gd_bundle.crt
	    `<FilesMatch "\.(cgi|shtml|phtml|php)$">`
	        SSLOptions +StdEnvVars
	    `</FilesMatch>`
	    `<Directory "/usr/local/apache2/cgi-bin">`
	        SSLOptions +StdEnvVars
	    `</Directory>`
	    BrowserMatch ".*MSIE.*" \
	             nokeepalive ssl-unclean-shutdown \
	             downgrade-1.0 force-response-1.0
	    CustomLog "/var/log/httpd/bidpit.com-ssl_request_log" \
	              "%t %h %{SSL_PROTOCOL}x %{SSL_CIPHER}x \"%r\" %b"
	    ServerAdmin webmaster@bidpit.com
	    DocumentRoot /data/bidpit
	    ServerName test.bidpit.com
	    ServerAlias beta.bidpit.com
	    ErrorLog /var/log/httpd/bidpit.com-ssl_error_log
	    DirectoryIndex index.php
	`</VirtualHost>`
	#
	#-30-


`</code>`

# Misc

Use python to start a simple HTTP server on port 8000. Note that it runs as the user you are logged in as and serves up the pwd where you ran the command.
`python -m SimpleHTTPServer`
