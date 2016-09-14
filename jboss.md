# Install Proper JDK

Install the proper java developmnet toolkit

	
	[root@jboss ~]# yum -y install java-1.7.0-openjdk-devel


Issue the following command to confirm that the proper version of the JDK is on your classpath

	
	[root@jboss ~]# java -version
	java version "1.7.0_65"
	OpenJDK Runtime Environment (rhel-2.5.1.2.el6_5-x86_64 u65-b17)
	OpenJDK 64-Bit Server VM (build 24.65-b04, mixed mode)


For our installation, we are not defining a explicit JAVA_HOME for JBoss AS 7. The default works in this situation, because we don’t have multiple java versions installed. For most production environments with multiple versions of Java, it is recommended to set the JAVA_HOME in the standalone.conf or domain.conf files.

# Download JBOSS

The next step is to download the appropriate version of JBoss AS 7. We will download the .zip version of JBoss AS 7, and install it using the unzip utility.

	
	[root@jboss ~]# curl -s -O http://download.jboss.org/jbossas/7.1/jboss-as-7.1.1.Final/jboss-as-7.1.1.Final.zip 
	[root@jboss ~]# ll -d jboss-as-7.1.1.Final.zip 
	-rw-r--r--. 1 root root 133255203 Jul 28 18:57 jboss-as-7.1.1.Final.zip


Next, we issue the following unzip command to finally install jboss in the // /usr/local/jboss // directory:

	
	[root@jboss ~]# unzip jboss-as-7.1.1.Final.zip -d /usr/local
	[root@jboss ~]# cd /usr/local/
	[root@jboss local]# ln -s jboss-as-7.1.1.Final jboss
	[root@jboss local]# ll -d jboss
	lrwxrwxrwx. 1 root root 20 Jul 28 19:00 jboss -> jboss-as-7.1.1.Final


# Create a JBOSS user

Now that JBoss AS 7, is installed, we need to make sure that we create a user with the appropriate privileges. It is never a good idea to run JBoss as root for various reasons.
Create the new user:

We create a new user called jboss by issuing the following command:

	
	[root@jboss local]# useradd jboss


We need to assign the appropriate ownership to the installation directory for the newly created jboss user by issuing the command:

	
	[root@jboss ~]# chown -R jboss.jboss /usr/local/jboss
	[root@jboss ~]# ll -d /usr/local/jboss
	lrwxrwxrwx. 1 jboss jboss 20 Jul 28 19:00 /usr/local/jboss -> jboss-as-7.1.1.Final
	


We switch to the jBoss user, so that this new installation can be administered properly. It is not recommended to administer JBoss as root.

	
	[root@jboss ~]# su - jboss
	[jboss@jboss ~]$ id
	uid=500(jboss) gid=500(jboss) groups=500(jboss) context=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023


# Add JBOSS Managment User

Now, lets change directories to the JBoss bin directory. This dorectory contains the necessary scripts to start, stop and manage your JBoss installation.

The final step before we start JBoss, is to add a management user. This is an internal JBoss management user, necessary to access the new JBoss management console.

You should see the following message on the console after executing the command:

	
	[jboss@jboss ~]$ cd /usr/local/jboss/bin/
	[jboss@jboss bin]$ ./add-user.sh 
	
	What type of user do you wish to add? 
	 a) Management User (mgmt-users.properties) 
	 b) Application User (application-users.properties)
	(a): a
	
	Enter the details of the new user to add.
	Realm (ManagementRealm) : 
	Username : jboss
	Password : 
	Re-enter Password : 
	About to add user 'jboss' for realm 'ManagementRealm'
	Is this correct yes/no? yes
	Added user 'jboss' to file '/usr/local/jboss-as-7.1.1.Final/standalone/configuration/mgmt-users.properties'
	Added user 'jboss' to file '/usr/local/jboss-as-7.1.1.Final/domain/configuration/mgmt-users.properties'
	


We select the default value for the Realm (ManagementRealm), by hitting enter, and select “jboss” as our username. By default, we supply “jb0ss” as our password, of course, you can provide any password you prefer here.

# Starting JBOSS

## Manually

A standalone instance of JBoss 7 can be starting by executing

	
	[jboss@jboss bin]$ ./standalone.sh -Djboss.bind.address=0.0.0.0 -Djboss.bind.address.management=0.0.0.0 &
	


By default, JBoss 7 will only bind to localhost. This does not allow any remote access to your jboss server. For our amazon aws installation, we define the jboss.bind.address property as *0.0.0.0* and jboss.bin.address.management property to *0.0.0.0* as well. This allows us to access the remote JBoss amazon instance over the internet. We could have also defined the hostname of the ami or the ip address. However, unless an elastic ip is used, this value can change. This is why we opted for *0.0.0.0*.

To shut it down we use this command

	
	[jboss@jboss bin]$ ./jboss-cli.sh --connect command=:shutdown

## INIT

To Use the init script first cd into the directoy

	
	[jboss@jboss init.d]$ cd /usr/local/jboss/bin/init.d


copy the init script 

	
	[root@jboss ~]# cd /usr/local/jboss/bin/init.d/
	[root@jboss init.d]# cp jboss-as-standalone.sh /etc/init.d/jboss  
	[root@jboss init.d]# chmod 755 /etc/init.d/jboss  
	[root@jboss init.d]# chkconfig --add jboss 
	[root@jboss init.d]# chkconfig --level 234 jboss on 


Make sure you edit the init script so that the paths are correct

	
	# Load JBoss AS init.d configuration.
	if [ -z "$JBOSS_CONF" ]; then
	  ##JBOSS_CONF="/etc/jboss-as/jboss-as.conf"
	  JBOSS_CONF="/usr/local/jboss/bin/init.d/jboss-as.conf"
	fi
	
	[ -r "$JBOSS_CONF" ] && . "${JBOSS_CONF}"
	
	# Set defaults.
	
	if [ -z "$JBOSS_HOME" ]; then
	  JBOSS_HOME=/usr/local/jboss

By default, JBoss 7.1.1 is bound to the loopback IP of 127.0.0.1, so if we want to make it available on the web, we need to change this.

Locate standalone.xml

	
	/usr/local/jboss/standalone/configuration/standalone.xml



Open standalone.xml in vi or a text editor and look for the public interfaces node as shown below.

	
	`<interface name="public">`
	`<inet-address value="${jboss.bind.address:127.0.0.1}"/>`
	`</interface>`


To make JBoss publicly accessible, change 127.0.0.1 to either 0.0.0.0 to allow access on all interfaces or to your public IP.

Now start the service

	
	service jboss start

