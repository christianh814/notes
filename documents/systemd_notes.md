# SystemD Notes

This is my `systemd` notes, in no paticular order

* [SystemV to SystemD Cheatsheet](#systemv-to-systemd-cheatsheet)
* [Runlevel And Targets](#runlevel-and-targets)
* [Changing runlevels](#changing-runlevels)
* [Creating a Systemd script](#creating-a-systemd-script)
* [RC Local](#rc-local)
* [Pager](#pager)

## SystemV to SystemD Cheatsheet

| `sysvinit` Command | `systemd` Command | Notes |
| -------------------- | ------------------- | ------- |
| service frobozz start | systemctl start frobozz.service | Used to start a service (not reboot persistent) |
| service frobozz stop | systemctl stop frobozz.service | Used to stop a service (not reboot persistent) |
| service frobozz restart | systemctl restart frobozz.service | Used to stop and then start a service |
| service frobozz reload | systemctl reload frobozz.service | When supported, reloads the config file without interrupting pending operations. |
| service frobozz condrestart | systemctl condrestart frobozz.service | Restarts if the service is already running. |
| service frobozz status | systemctl status frobozz.service | Tells whether a service is currently running. |
| ls /etc/rc.d/init.d/ | systemctl list-unit-files --type=service | Used to list all the services and other units |
| chkconfig frobozz on | systemctl enable frobozz.service | Turn the service on, for start at next boot, or other trigger. |
| chkconfig frobozz off | systemctl disable frobozz.service | Turn the service off for the next reboot, or any other trigger. |
| chkconfig frobozz | systemctl is-enabled frobozz.service | Used to check whether a service is configured to start or not in the current environment. |
| chkconfig --list | systemctl list-unit-files --type=service |Print a table of services that lists which runlevels each is configured on or off |
| chkconfig frobozz --list | ls /etc/systemd/system/*.wants/frobozz.service | Used to list what levels this service is configured on or off |
| chkconfig frobozz --add | systemctl daemon-reload | Used when you create a new service file or modify any configuration |

Note that all `/sbin/service/` and `/sbin/chkconfig` lines listed above continue to work on `systemd`, and will be translated to native equivalents as necessary. The only exception is `chkconfig --list`

Any service that defines an additional command in this way would need to define some other, service-specific, way to accomplish this task when writing a native systemd service definition.
Check the package-specific release notes for any services that may have done this.

## Runlevel And Targets

Systemd has a concept of targets which serve a similar purpose as runlevels but act a little different. Each target is named instead of numbered and is intended to serve a specific purpose. Some targets are implemented by inheriting all of the services of another target and adding additional services to it. There are systemd targets that mimic the common sysvinit runlevels so you can still switch targets using the familiar telinit RUNLEVEL command. The runlevels that are assigned a specific purpose on vanilla Fedora installs; 0, 1, 3, 5, and 6; have a 1:1 mapping with a specific systemd target. Unfortunately, there's no good way to do the same for the user-defined runlevels like 2 and 4. If you make use of those it is suggested that you make a new named systemd target as `/etc/systemd/system/$YOURTARGET` that takes one of the existing runlevels as a base (you can look at `/lib/systemd/system/graphical.target` as an example), make a directory `/etc/systemd/system/$YOURTARGET.wants`, and then symlink the additional services that you want to enable into that directory. (The service unit files that you symlink live in `/lib/systemd/system`).

| sysvinit Runlevel | systemd Target | Notes |
| ----------------- | -------------- | ----- |
| 0 | runlevel0.target, poweroff.target | Halt the system. |
| 1, s, single| runlevel1.target, rescue.target | Single user mode. |
| 2, 4 | runlevel2.target, runlevel4.target, multi-user.target | User-defined/Site-specific runlevels. By default, identical to 3. |
| 3 | runlevel3.target, multi-user.target | Multi-user, non-graphical. Users can usually login via multiple consoles or via the network. |
| 5 | runlevel5.target, graphical.target | Multi-user, graphical. Usually has all the services of runlevel 3 plus a graphical login. |
| 6 | runlevel6.target, reboot.target | Reboot |
| emergency | emergency.target | Emergency shell |

To change the root password add `rd.break` at the end of the `linux16` line like so...

```
linux16 /vmlinuz-3.10.0-123.8.1.el7.x86_64 root=UUID=d51061fd-84aa-4589-ab00-3bef08cf3143 ro ... rd.lvm.lv=vg0/root crashkernel=auto  vconsole.keymap=us rhgb quiet rd.break
```

It'll bring you to a shell where you can mount the `/sysroot` (your filesystem) as `rw`.

```
bash# mount -o remount,rw /sysroot
```

Now `chroot` into it and change the root password
```
bash# chroot /sysroot
root@host# passwd root
```

Make sure to relabel or it might not boot again
```
root@host# touch /.autorelabel
```

## Changing runlevels

| sysvinit Command | systemd Command | Notes |
| ---------------- | --------------- | ----- |
| `telinit 3` | systemctl isolate multi-user.target (OR systemctl isolate runlevel3.target OR telinit 3) | Change to multi-user run level. |
| `sed s/^id:.*:initdefault:/id:3:initdefault:/` | systemctl set-default multi-user.target | Set to use multi-user runlevel on next reboot. |
| `who -r` | systemctl get-default | Outputs current "runlevel" |

## Creating a Systemd script

If you have compiled software or need to make sure that something starts on boot (and it doesn't have a systemd script) - you will need to create one

All systemd scripts exist on `/usr/lib/systemd/system` and you usually need to symlink it to `/etc/systemd/system`

I usually just copy and edit an existing script (this is especially useful when you are on a system that has SELinux turned on).

```
[root@iFedora system]# cd /usr/lib/systemd/system/
[root@iFedora system]# cp sshd.service apache2.service
```

Now edit the file to your liking ( left some things in with a comment and will explain each of these options)

(NOTE: I had to add `PidFile "/var/run/apache22.pid"` in the httpd.conf file in order for this to work...but this was more of an Apache issue and not systemd specific)

```
[Unit]
Description=Apache 2.2 server daemon
After=rsyslog.service network.target auditd.service

[Service]
#EnvironmentFile=/etc/sysconfig/apache22
#ExecStart=/usr/local/bin/apachectl $OPTS
Type=forking
PIDFile=/var/run/apache22.pid
ExecStartPre=/bin/rm -f /var/run/apache22.pid
ExecStart=/usr/local/bin/apachectl start
ExecReload=/usr/local/bin/apachectl graceful
ExecStop=/usr/local/bin/apachectl stop
ExecStopPost=/bin/rm -f /var/run/apache22.pid

[Install]
WantedBy=multi-user.target
```

Info:
  * `Description` - A short description of what the service is
  * `After` - What needs to be run before this service is started (i.e. dependencies)
  * `EnvironmentFile` - This is where you put your variables (you can see how it's used on the following ExecStart line)
  *  Format should be something like this (note you are NOT supposed to put these in quotes): `OPTS=-k start`
  * `Type` - this can be simple, forking, oneshot, dbus, notify. In this case it's "forking" because (from the "man systemd.service" page) "This is the behaviour of traditional UNIX daemons"
  * `PIDFile` - self explanatory, path to the pid file (systemd doesn't create this file...this is the file that the process has to create - so make sure it does)
  * `ExecStartPre` - things to run before starting the service
  * `ExecStart` - what to run to start the service
  * `ExecReload` - Optional command to put to reload the service
  * `ExecStop` - Command to stop the service
  * `ExecStopPost` - Any commands to run after service stops, clean ups and such
  * `WantedBy` - Symlinks in the ` .wants/ ` resp.   ` .requires/ ` subdirectory for a unit. This has the effect that when the listed unit name is activated the unit listing it is activated too. (i.e. starts it on boot)

Now that you have the systemd script written. You have to reload the systemd daemon

```
[root@iFedora system]# systemctl --system daemon-reload
```

Make sure it starts on boot

```
[root@iFedora system]# systemctl enable apache2.service
ln -s '/usr/lib/systemd/system/apache2.service' '/etc/systemd/system/multi-user.target.wants/apache2.service'
```

This should be it.

UPDATE:

I set the environment variables like so (in the ` /etc/sysconfig/apache22 ` file)...

```
## Start String
START=start
## Stop String
STOP=stop
## Reload String
RELOAD=graceful
```

Now my ` apache2.service ` looks like this (I removed the "ExecStartPre" and the "ExecStopPost" since apache creates and removes these on it's own...YOUR service may not...it's service specific):

```
[Unit]
Description=Apache 2.2 server daemon
After=rsyslog.service network.target auditd.service

[Service]
EnvironmentFile=/etc/sysconfig/apache22
Type=forking
PIDFile=/var/run/apache22.pid
ExecStart=/usr/local/bin/apachectl $START
ExecReload=/usr/local/bin/apachectl $RELOAD
ExecStop=/usr/local/bin/apachectl $STOP

[Install]
WantedBy=multi-user.target
```

I just needed to remember to reload the daemon after I edited the service file

```
systemctl --system daemon-reload
```

### Oneshot Script

The "oneshot" is good for things that start but aren't a service (for example oc cluster up)

```
[Unit]
Description=demo script
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/demo-start-script.sh
RemainAfterExit=true
ExecStop=/usr/local/bin/demo-stop-script.sh
StandardOutput=journal

[Install]
WantedBy=multi-user.target
```


Example from Red Hat
```
[Unit]
Description=Example Service Script description goes here
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/sbin/example.sh
TimeoutStartSec=0

[Install]
WantedBy=default.target
```

This is what I did to set up hostname on boot

```[root@dhcp-host-29 ~]# cat /etc/systemd/system/set-hostname.service
[Unit]
Description=Sets the hostname based on ip addr
After=rsyslog.service network.target auditd.service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/set-hostname
TimeoutStartSec=0

[Install]
WantedBy=default.target


[root@dhcp-host-29 ~]# cat /usr/local/bin/set-hostname
#!/bin/bash
# Set hostname on boot like cloud-init would
mycounter=0
while [ ${mycounter} -lt 7 ]; do
  if ! ping -c 1 192.168.1.1 > /dev/null 2>&1 ; then
    sleep 5
    mycounter=$[${mycounter} + 1 ]
  else
    /usr/bin/hostnamectl set-hostname $(/usr/bin/dig -x $(/usr/bin/hostname -I | awk '{print $1}') +short)
  fi
done
##
##

```

## RC Local

Make the file executable

```
chmod +x /etc/rc.d/rc.local
```

After that put the commands you want to run on boot inside the `/etc/rc.d/rc.local` file.

## Pager


Set your favorite pager with

```
export PAGER=less
```

But then set an alias for `systemctl` to use no pager

```
alias systemctl='systemctl --no-pager'
```

# Journalctl

Journalctl may be used to query the contents of the systemd journal as written by `systemd-journald.service`

if called without parameters, it will show the full contents of the journal, starting with the oldest entry collected (analogues to "cat /var/log/messages")

```
root@host# journalctl
```

You can "tail" the entries with a '-f' passed

```
root@host# journalctl -f
```

Although you can filter stuff out using switches; you can grep for what you need
```
root@host# journalctl | grep kernel
```

With one match specified, all entries with a field matching the expression are shown:

```
root@host# journalctl _SYSTEMD_UNIT=avahi-daemon.service
```

If two different fields are matched, only entries matching both expressions at the same time are shown:

```
root@host# journalctl _SYSTEMD_UNIT=avahi-daemon.service _PID=28097
```

If two matches refer to the same field, all entries matching either expression are shown:

```
root@host# journalctl _SYSTEMD_UNIT=avahi-daemon.service _SYSTEMD_UNIT=dbus.service
```

If the separator "+" is used, two expressions may be combined in a logical OR. The following will show all messages from the Avahi service process with the PID 28097 plus all messages from the D-Bus service (from any of its processes):

```
root@host# journalctl _SYSTEMD_UNIT=avahi-daemon.service _PID=28097 + _SYSTEMD_UNIT=dbus.service
```

Show all logs generated by the D-Bus executable:

```
root@host# journalctl /usr/bin/dbus-daemon
```

Show all logs of the kernel device node /dev/sda:

```
root@host# journalctl /dev/sda
```

Show all kernel logs from previous boot:

```
root@host# journalctl -k -b -1
```

