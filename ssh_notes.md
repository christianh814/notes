# Password-less SSH Login


Sometimes administrative tasks require to connect to a machine via ssh without being asked for a password, this can be accomplished by doing tthe follwoing.

On the machine that you want to ssh from create the key

    root@host1# ssh-keygen -b 1024 -t rsa

The id_pub.rsa file will be in the ~/.ssh directory. Now scp the id_pub.rsa key to the machine you want to ssh to

    root@host1# ssh-copy-id -i ~/.ssh/id_rsa.pub host2

Now login into host2 and add the key to the authorized_keys file

    root@host1# ssh host2
    Password:
    root@host2# cd ~/.ssh
    root@host2# cat id_rsa.pub >> authorized_keys

Now you should be able to ssh into host2 into host1 and it wont promp’t you for a password. If you are having troubles make sure that authorized_keys is set to 664 permissions. 

# SSH Tunnels

You can establish "tunnels" via SSH.

Bind your local port 8080 to remote server's 80

	
	root@host# ssh host2 -L 8080:host2:80


Then on a web browser you can go to...

	
	http://localhost:8080


And see the site.

This is not restricted to the remote host. For example If you can "see" a Linux host in a network but you can't see a Windows machine that you want to RDP into (that's in the same network); you can establish a tunnel.

	
	root@host# ssh host2 -L 3389:windowshost:3389


Then on your RDP Client type in

	
	localhost:3389


You can use ssh tunnels to "bounce" around a network.

__Gateways__

You can use SSH as a gateway. I did the following at 4over to set up a consisten SSH tunnel gateway that runs in the background

	
	root@host# ssh -i /usr/local/lib/cci.key -CfNg 216.133.247.10 -l over4 -L 1433:192.180.1.7:1433


Basically I created an SSH key first then used port fowarding with special options.

*** -i /usr/local/lib/cci.key ** * The identity file to use \\
*** -C ** * Use Compression\\
*** -f ** * Requests ssh to go to background just before command execution\\
*** -N ** * Do not execute a remote command.  This is useful for just forwarding ports\\
*** -g ** * Allows other machines to connect to this machine (gateway mode)\\
*** -l over4 *** login account to use\\
*** -L 1433:192.180.1.7:1433 *** Use local port 1433 and bind it to 192.180.1.7:1433\\

For a "reverse tunnel" (i.e. tunnel back into the server I'm ssh-ing from)

	
	root@host# ssh -CfN cci-bridge.4over.com -l root -R 2222:localhost:22

# User Security

You can allow who can ssh into the system within the /etc/ssh/sshd_config file.
        root@host# vi /etc/ssh/sshd_config
          AllowUsers mike chrish

You can also restrict from where they ssh from
        root@host# vi /etc/ssh/sshd_config
          AllowUsers mike@192.168.2.122 chrish@192.168.2.54

Other options include (which are pretty self explanitory)
        AllowGroups
        DenyUsers
        DenyGroups

SSH logs to: /var/log/secure

# X11 Forwarding

First install on the server you want Forwading on

	
	root@host# yum -y install xorg-x11-xauth xorg-x11-utils xorg-x11-fonts-* 



The next step is to check the configuration of the SSH service running on the server. By default, the SSH
server on Red Hat Enterprise Linux 5 and Red Hat Enterprise Linux 6 has variable X11 forwarding enabled
in file // /etc/ssh/sshd_config // through variable // X11Forwarding yes//. Ensure that this has not been
changed. If it has changed, set the variable to yes and restart the sshd.

	
	root@host# service sshd restart


Now you can ssh into the server and run X11 Apps

	
	root@host# ssh -X host2


# Misc Commands

Keep SSH session on if you keep getting kicked out

	
	ssh -o ServerAliveInterval=15 -o ServerAliveCountMax=4 $HOST


Explanation:

*  // ServerAliveInterval // ~ Sets a timeout interval in seconds after which if no data has been received from the server, ssh will send a message through the encrypted channel to request a response from the server.

*  // ServerAliveCountMax // ~ Sets the number of server alive messages which may be sent without ssh receiving any messages back from the server.  If this threshold is reached while server alive messages are being sent, ssh will disconnect from the server, terminating the session

Using SSH Shell (ssh-agent)

	
	$ exec /usr/bin/ssh-agent $SHELL
	
	$ ssh-add -l
	The agent has no identities.
	
	$ ssh-add ~/.ssh/github
	Identity added: /home/chernand/.ssh/github (/home/chernand/.ssh/github)
	
	$ ssh-add -l
	

