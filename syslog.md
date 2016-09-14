# On EL 6

The */etc/rsyslog.conf* file is pretty well commented.

## Server Setup

Edit the // /etc/rsyslog.conf // file and uncomment entries that begin with // $ //

It should look something like this

	
	[root@loghost ~]# sed -n '/# ### begin forwarding rule ###/,$p' /etc/rsyslog.conf | sed -n '1,/^# ### end of the forwarding rule ###/p'
	# ### begin forwarding rule ###
	# The statement between the begin ... end define a SINGLE forwarding
	# rule. They belong together, do NOT split them. If you create multiple
	# forwarding rules, duplicate the whole block!
	# Remote Logging (we use TCP for reliable delivery)
	#
	# An on-disk queue is created for this action. If the remote host is
	# down, messages are spooled to disk and sent when it is up again.
	$WorkDirectory /var/lib/rsyslog # where to place spool files
	$ActionQueueFileName fwdRule1 # unique name prefix for spool files
	$ActionQueueMaxDiskSpace 1g   # 1gb space limit (use as much as possible)
	$ActionQueueSaveOnShutdown on # save messages to disk on shutdown
	$ActionQueueType LinkedList   # run asynchronously
	$ActionResumeRetryCount -1    # infinite retries if host is down
	# remote host is: name/ip:port, e.g. 192.168.0.1:514, port optional
	#*.* @@remote-host:514
	# ### end of the forwarding rule ###


Make sure the following is also in the // /etc/rsyslog.conf // file

	
	# Provides TCP syslog reception
	$ModLoad imtcp
	$InputTCPServerRun 514


Then make sure you open the popper ports

	
	[root@loghost ~]# lokkit --port=514:tcp


Start the service

	
	[root@loghost ~]# chkconfig rsyslog on
	[root@loghost ~]# service rsyslog restart
	Shutting down system logger:                               [  OK  ]
	Starting system logger:                                    [  OK  ]


## Client Setup

On the client side...edit the // /etc/rsyslog.conf // file to tell it to send the logs to a remote server

It should look something like this

	
	[root@broker ~]#  sed -n '/# remote host is/,$p' /etc/rsyslog.conf | sed -n '1,/^# ### end of the forwarding rule ###/p'
	# remote host is: name/ip:port, e.g. 192.168.0.1:514, port optional
	#*.* @@remote-host:514

	*.* @@loghost.example.net:514
	# ### end of the forwarding rule ###
	


You can use IP addresses instead of dns name.

Then restart // rsyslog //

	
	service rsyslog restart


# On EL 7

he */etc/rsyslog.conf* file is pretty well commented.


## Custom Messages

You configure rsyslog to write specific messages to a new log file. 

First configure rsyslog to write specific messages to a new log file. 

	
	[root@host ~]# echo "*.debug /var/log/messages-debug" >/etc/rsyslog.d/debug.conf


Restart the rsyslog service on serverX.

	
	[root@host ~]# systemctl restart rsyslog



Monitor the // /var/log/messages-debug // with the tail command

    [root@host ~]# tail -f /var/log/messages-debug

On a separate terminal window, use the logger command to generate a debug message.

    [root@host ~]# logger -p user.debug "Debug Message Test"

Switch back to the terminal still running the // tail -f /var/log/messages-debug // command and verify the message sent with the logger command shows up.

	
	    [root@serverX ~]# tail -f /var/log/messages-debug
	    ...
	    Feb 13 10:37:44 localhost root: Debug Message Test

