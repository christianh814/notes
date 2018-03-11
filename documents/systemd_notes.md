# SystemD Notes

This is my `systemd` notes, in no paticular order

* [SystemV to SystemD Cheatsheet](#systemv-to-systemd-cheatsheet)

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

