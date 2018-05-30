# Rocket On Fedora

Here are some brief notes and simple use case on getting `rkt` working on Fedora

* [Set Up](#set-up)
* [Fedora Image](#fedora-image)
* [Create Basic Apache with PHP Image](#create-basic-apache-with-php-image)
* [Build The Container](#build-the-container)
* [Running Rocket in systemd](#running-rocket-in-systemd)
* [Unit File](#unit-file)
* [Advanced Systemd Unit](#advanced-systemd-unit)

## Set Up

After you install Fedora; disable `firewalld` and `selinux` as `rkt` doesn't play well with them just yet

```
systemctl disable firewalld
sed -i.bak 's/SELINUX\=enforcing/SELINUX\=disabled/g' /etc/selinux/config
```

Install `rkt` it's there in the default repos

```
dnf -y install rkt golang-github-appc-spec
```

Visit [this page](https://github.com/containers/build/releases) and download the latest `acbuild` binaries. At the time of this writing it's 0.4.0

```
cd /usr/local/src/
wget https://github.com/containers/build/releases/download/v0.4.0/acbuild-v0.4.0.tar.gz
tar -xzf acbuild-v0.4.0.tar.gz 
cp acbuild-v0.4.0/* /usr/local/bin/
```

Make sure everything is working well

```
acbuild --help
rkt --version
```

At this point go ahead and update the Fedora vm and reboot

```
dnf -y clean all && dnf -y update && systemctl reboot
```

## Fedora Image

Now, visit [here](https://github.com/fedora-cloud/docker-brew-fedora) or [here](https://koji.fedoraproject.org/koji/tasks?start=0&state=all&view=tree&method=image&order=-id) and download the latest base Fedora Docker image. We will be extracting this and using it as our base image.

```
cd ~
mkdir fedora-base
wget https://raw.githubusercontent.com/fedora-cloud/docker-brew-fedora/0b198330cf019e615e1393affc4338664a9dd332/fedora-25-20170420.tar.xz 
xz -d Fedora-Docker-Base-25-20170310.0.x86_64.tar.xz 
tar -xf Fedora-Docker-Base-25-20170310.0.x86_64.tar -C fedora-base
```

Now go into this dir and remove everything you don't need (basically extract the base `rootfs` dir)

```
cd fedora-base
HASH=$(cat repositories | awk -F '"latest": "' '{ print $2 }' | awk '{ sub(/[^a-zA-Z0-9]+/, ""); print }')
echo $HASH
mv $HASH/layer.tar .
rm -rf $HASH repositories
tar -xf layer.tar --same-owner --preserve-permissions
rm layer.tar 
```

Now you should have your `fedora-base` dir with the following contents

```
[root@fbox fedora-base]# pwd
/root/fedora-base
[root@fbox fedora-base]# ls -l
total 64
lrwxrwxrwx  1 root root    7 Feb  3  2016 bin -> usr/bin
dr-xr-xr-x  2 root root 4096 Feb  3  2016 boot
drwxr-xr-x  2 root root 4096 Mar  9 21:48 dev
drwxr-xr-x 45 root root 4096 Mar  9 21:49 etc
drwxr-xr-x  2 root root 4096 Mar  9 21:49 home
lrwxrwxrwx  1 root root    7 Feb  3  2016 lib -> usr/lib
lrwxrwxrwx  1 root root    9 Feb  3  2016 lib64 -> usr/lib64
drwx------  2 root root 4096 Mar  9 21:48 lost+found
drwxr-xr-x  2 root root 4096 Feb  3  2016 media
drwxr-xr-x  2 root root 4096 Feb  3  2016 mnt
drwxr-xr-x  2 root root 4096 Feb  3  2016 opt
drwxr-xr-x  2 root root 4096 Mar  9 21:48 proc
dr-xr-x---  2 root root 4096 Mar  9 21:49 root
drwxr-xr-x  2 root root 4096 Mar  9 21:48 run
lrwxrwxrwx  1 root root    8 Feb  3  2016 sbin -> usr/sbin
drwxr-xr-x  2 root root 4096 Feb  3  2016 srv
drwxr-xr-x  2 root root 4096 Mar  9 21:48 sys
drwxrwxrwt  7 root root 4096 Mar  9 21:49 tmp
drwxr-xr-x 12 root root 4096 Mar  9 21:49 usr
drwxr-xr-x 19 root root 4096 Mar  9 21:49 var
```

## Create Basic Apache with PHP Image

Use this `fedora-base` to create an image that has Apache and PHP. First copy this dir

```
cd ~
cp -a fedora-base fedora-httpd
```

Now use `dnf` to install the required packages

```
dnf -y --installroot=/root/fedora-httpd --setopt=tsflags='nodocs' --setopt="override_install_langs=en_US.utf8" install httpd php
```

Add an `info.php` file

```
echo '<?php phpinfo(); ?>' >> /root/fedora-httpd/var/www/html/info.php
```

## Build The Container

Use `acbuild` to build the container.

```
acbuild begin /root/fedora-httpd
acbuild set-name chx.local/fedora-httpd
acbuild label add version latest
acbuild set-exec  -- /usr/sbin/httpd -D FOREGROUND
acbuild annotation add appc.io/executor/supports-systemd-notify true
acbuild port add http tcp 80
acbuild write fedora-httpd.aci
acbuild end
```

What we're doing here is:

* Telling acbuild to use our `fedora-httpd` rootfs that we extracted as the rootfs for the container
* Setting the name to `your domain>/fedora` For example, I would use `chx.local/fedora-httpd`
* We set the label "version" to "latest"
* Set our exec (i.e. "entrypoint") to apache running in the foreground
* Set our container to notify systemd (more on this later)
* Apache listens on 80 so we are going to name this port (not required); so we can reference it later
* We write out everything to fedora-httpd.aci. This what we will actually run with rkt.
* Tell acbuild we're done

Use `actool` to validate (it should return nothing if all is okay)

```
actool validate fedora-httpd.aci 
```

## Running Rocket in systemd

Run your new container in systemd

```
systemd-run --slice=machine rkt run --dns=8.8.8.8 --net=all --insecure-options=all --port=http:8080 /root/fedora-httpd.aci
```

In a few; you can visit your server's IP's address on port 8080 and you should see the fedora test page.

List your machines; stop it at this point

```
[root@fbox ~]# machinectl list
MACHINE                                  CLASS     SERVICE
rkt-6c62e5d9-144a-4dd3-9582-3676e5f08f23 container rkt    

1 machines listed.
[root@fbox ~]# machinectl stop rkt-6c62e5d9-144a-4dd3-9582-3676e5f08f23
```

## Unit File

Create a simple unit file called `fedora-httpd.service` under `etc/systemd/system`

```
[Unit]
Description=fedora-apache

[Service]
Slice=machine.slice
ExecStart=/usr/bin/rkt run --dns=8.8.8.8 --net=all --insecure-options=all --port=http:8080 chx.local/fedora-httpd:latest
KillMode=mixed
Restart=always
```

After you add this run

```
systemctl --system daemon-reload
```

Now you can start it like any other systemd service

```
[root@fbox ~]# systemctl start fedora-httpd.service
[root@fbox ~]# systemctl status fedora-httpd.service
● fedora-httpd.service - fedora-apache
   Loaded: loaded (/etc/systemd/system/fedora-httpd.service; static; vendor preset: disabled)
   Active: active (running) since Mon 2017-03-20 16:01:20 PDT; 5s ago
 Main PID: 2698 (systemd-nspawn)
    Tasks: 35 (limit: 4915)
   CGroup: /machine.slice/fedora-httpd.service
           ├─2698 /usr/bin/systemd-nspawn --boot --register=true --link-journal=try-guest --keep-unit --quiet --uuid=d4de8c30-98b8-4880-8116-ad3c3bbec9c9 --machine=rkt-d4de8
           ├─init.scope
           │ └─3136 /usr/lib/systemd/systemd --default-standard-output=tty --log-target=null --show-status=0
           └─system.slice
             ├─fedora-httpd.service
             │ ├─3143 /usr/sbin/httpd -D FOREGROUND
             │ ├─3162 /usr/sbin/httpd -D FOREGROUND
             │ ├─3163 /usr/sbin/httpd -D FOREGROUND
             │ ├─3164 /usr/sbin/httpd -D FOREGROUND
             │ ├─3165 /usr/sbin/httpd -D FOREGROUND
             │ └─3167 /usr/sbin/httpd -D FOREGROUND
             └─systemd-journald.service
               └─3138 /usr/lib/systemd/systemd-journald

Mar 20 16:01:20 fbox systemd[1]: Started fedora-apache.
Mar 20 16:01:20 fbox rkt[2698]: image: using image from local store for image name coreos.com/rkt/stage1-host:1.11.0-1.git1ec4c60.fc25
Mar 20 16:01:20 fbox rkt[2698]: image: using image from local store for image name chx.local/fedora-httpd:latest
Mar 20 16:01:20 fbox rkt[2698]: networking: loading networks from /etc/rkt/net.d
Mar 20 16:01:20 fbox rkt[2698]: networking: loading network default with type ptp
Mar 20 16:01:20 fbox rkt[2698]: [ 1820.604808] fedora-httpd[5]: AH00557: httpd: apr_sockaddr_info_get() failed for rkt-d4de8c30-98b8-4880-8116-ad3c3bbec9c9
Mar 20 16:01:20 fbox rkt[2698]: [ 1820.605128] fedora-httpd[5]: AH00558: httpd: Could not reliably determine the server's fully qualified domain name, using 127.0.0.1.
```

## Advanced Systemd Unit

Create a more complex unit file (note I changed the port number so you can have both running) as `fedora-httpd2.service` under `etc/systemd/system`

```
[Unit]
# Metadata
Description=fedora-apache2
Documentation=https://apache.org
# Wait for networking
Requires=network-online.target
After=network-online.target

[Service]
Slice=machine.slice
# Resource limits
Delegate=true
CPUShares=512
MemoryLimit=1G
#Only use 30% of CPU
CPUQuota=30%
# Pin to cpu3
CPUAffinity=0,3
# Env vars
#Environment=TMPDIR=/var/tmp
# Fetch the app (not strictly required, `rkt run` will fetch the image if there is not one)
ExecStartPre=/usr/bin/rkt fetch chx.local/fedora-httpd:latest
# Start the app
ExecStart=/usr/bin/rkt run --dns=8.8.8.8 --net=all --insecure-options=all --port=http:8081 chx.local/fedora-httpd:latest
ExecStopPost=/usr/bin/rkt gc --mark-only
KillMode=mixed
Restart=always
```

Reload and start the new service

```
[root@fbox system]# systemctl --system daemon-reload
[root@fbox system]# systemctl start fedora-httpd2.service
[root@fbox system]# systemctl status fedora-httpd2.service
● fedora-httpd2.service - fedora-apache2
   Loaded: loaded (/etc/systemd/system/fedora-httpd2.service; static; vendor preset: disabled)
   Active: active (running) since Mon 2017-03-20 16:05:44 PDT; 4s ago
     Docs: https://apache.org
  Process: 3222 ExecStartPre=/usr/bin/rkt fetch chx.local/fedora-httpd:latest (code=exited, status=0/SUCCESS)
 Main PID: 3246 (systemd-nspawn)
    Tasks: 35 (limit: 4915)
   Memory: 35.2M (limit: 1.0G)
      CPU: 457ms
   CGroup: /machine.slice/fedora-httpd2.service
           ├─3246 /usr/bin/systemd-nspawn --boot --register=true --link-journal=try-guest --keep-unit --quiet --uuid=33289317-d408-4a67-837b-2345cd47188d --machine=rkt-33289
           ├─init.scope
           │ └─3687 /usr/lib/systemd/systemd --default-standard-output=tty --log-target=null --show-status=0
           └─system.slice
             ├─fedora-httpd.service
             │ ├─3696 /usr/sbin/httpd -D FOREGROUND
             │ ├─3717 /usr/sbin/httpd -D FOREGROUND
             │ ├─3718 /usr/sbin/httpd -D FOREGROUND
             │ ├─3719 /usr/sbin/httpd -D FOREGROUND
             │ ├─3720 /usr/sbin/httpd -D FOREGROUND
             │ └─3721 /usr/sbin/httpd -D FOREGROUND
             └─systemd-journald.service
               └─3691 /usr/lib/systemd/systemd-journald

Mar 20 16:05:44 fbox systemd[1]: Starting fedora-apache2...
Mar 20 16:05:44 fbox rkt[3222]: image: using image from local store for image name chx.local/fedora-httpd:latest
Mar 20 16:05:44 fbox rkt[3222]: sha512-1dbb277c5c8a678c739075fb4495a9c2
Mar 20 16:05:44 fbox systemd[1]: Started fedora-apache2.
Mar 20 16:05:44 fbox rkt[3246]: image: using image from local store for image name coreos.com/rkt/stage1-host:1.11.0-1.git1ec4c60.fc25
Mar 20 16:05:44 fbox rkt[3246]: image: using image from local store for image name chx.local/fedora-httpd:latest
Mar 20 16:05:44 fbox rkt[3246]: networking: loading networks from /etc/rkt/net.d
Mar 20 16:05:44 fbox rkt[3246]: networking: loading network default with type ptp
Mar 20 16:05:45 fbox rkt[3246]: [ 2085.867176] fedora-httpd[5]: AH00557: httpd: apr_sockaddr_info_get() failed for rkt-33289317-d408-4a67-837b-2345cd47188d
Mar 20 16:05:45 fbox rkt[3246]: [ 2085.867603] fedora-httpd[5]: AH00558: httpd: Could not reliably determine the server's fully qualified domain name, using 127.0.0.1.
```

If you list your Machines there should be two...

```
[root@fbox system]# machinectl list
MACHINE                                  CLASS     SERVICE
rkt-33289317-d408-4a67-837b-2345cd47188d container rkt    
rkt-d4de8c30-98b8-4880-8116-ad3c3bbec9c9 container rkt    

2 machines listed.
```

# Red Hat and Centos

Currently; you can't run `rkt` on Red Hat because of an [incompatable version of systemd](https://github.com/coreos/rkt/issues/1305)

You __CAN__ however; build a RHEL image on a RHEL 7 box and just copy the ACI over and run it on your Fedora Box.

**UPDATE:** Looks like you can now. 

Info [here](https://wiki.centos.org/Cloud/rkt) and [here](http://www.jebriggs.com/blog/2016/11/linux-rkt-on-centos7-is-just-too-easy/)

