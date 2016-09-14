# EL6 VNC

Following is for EL6 VNC Configuration

## VNC Server Configuration

Install the applicable packages
        root@host# yum -y install vino vinagre tigervnc tiger-vnc server

The configuration file (/etc/sysconfig/vncserver) is pretty straight forward. It's
well commented and you have to worry about 2 options...
        root@host# service vncserver stop
        root@host# vi /etc/sysconfig/vncserver
          VNCSERVERS="2:myusername"
          VNCSERVERARGS[2]="-geometry 800x600"

The options are pretty straight forward the "2:myusername" is "port 5902 with the
username myusername". The Arguments just specify the dementions

Next you will need to add the ports in the firewall. Usually 5900-5909; but since we
are only using 5902; only 5900-5902 is needed. 
        root@host# system-config-firewall

          * Click "Other Ports"
          * Add range
          * 5900-5902 TCP (maybe UDP as well?) 

Make sure VNC starts on boot (don't start it just yet though) 
        root@host# chkconfig vncserver on
        
Now login as the user and run the command...(it will ask you for a password)
        root@host# su - myusername
        myusername@host$ vncserver :2
        
        Password:

You can set custom preferences if you'd like
        myusername@host@$ vino-preferences

## VNC Client Configuration

Now from the remote system you can connect (make sure you have the client vnc
system installed)
        root@host# yum -y install tigervnc
        root@host# vncserver server1.example.com:2

NOTE: IP Addresses work as well

# EL7 VNC

Following is for EL7 VNC Configuration

## VNC Server Conf

Frist, install *tigervnc*

	
	root@host# yum -y install tigervnc-server tigervnc  xorg-x11-fonts-Type1 


Install X Windows if not already installed

	
	root@host#  yum -y groupinstall "Server with GUI"




Copy a sample configuration file to the // /etc/systemd/system/ // directory and rename it to the next available vnc port (in this case "1" means "5901")

	
	root@host# cp /usr/lib/systemd/system/vncserver@.service "/etc/systemd/system/vncserver@:1.service"


Edit the file so that the *`<USER>`* is replaced by the username (I used sed...you can use what you want)

	
	root@host# user=username
	root@host# sed -i "s/`<USER>`/${user}/g" /etc/systemd/system/vncserver@\:1.service
	root@host# egrep -v '^#|^$' /etc/systemd/system/vncserver@\:1.service 
	[Unit]
	Description=Remote desktop service (VNC)
	After=syslog.target network.target
	[Service]
	Type=forking
	ExecStartPre=/bin/sh -c '/usr/bin/vncserver -kill %i > /dev/null 2>&1 || :'
	ExecStart=/sbin/runuser -l christian -c "/usr/bin/vncserver %i"
	PIDFile=/home/christian/.vnc/%H%i.pid
	ExecStop=/bin/sh -c '/usr/bin/vncserver -kill %i > /dev/null 2>&1 || :'
	[Install]
	WantedBy=multi-user.target


Set the password for the user logging in

	
	root@host# su - ${user} -c vncpasswd


NOTE: You can start "classic" session with //  gnome-session --session gnome-classic //

Add the proper firewall rules

	
	root@host# firewallzone=$(firewall-cmd --get-active-zone|head -1)
	root@host# firewall-cmd --permanent --zone=${fireallzone} --add-port=5901/tcp
	root@host# firewall-cmd  --reload
	root@host# firewall-cmd --list-all
	public (default, active)
	  interfaces: eth0
	  sources: 
	  services: dhcpv6-client ssh
	  ports: 5901/tcp
	  masquerade: no
	  forward-ports: 
	  icmp-blocks


Now start and enable the service (since we created a new systemd service file; we need to reload the daemon)

	
	root@host# systemctl daemon-reload
	root@host# systemctl enable vncserver@:1.service 
	ln -s '/etc/systemd/system/vncserver@:1.service' '/etc/systemd/system/multi-user.target.wants/vncserver@:1.service'
	root@host# systemctl start vncserver@:1.service

## VNC Client Conf


Now from the remote system you can connect (make sure you have the client vnc
system installed)
        root@host# yum -y install tigervnc
        root@host# vncserver server1.example.com:2

NOTE: IP Addresses work as well
