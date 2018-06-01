# Yum Notes

Random `yum` notes in no paticular order

* [Create local yum repo](#create-local-yum-repo)

## Create local yum repo

These are quick notes

### From DVD

Copy DVD contents to a directory

```
root@host# mount rhel6-x86_64.iso /mnt -o loop
root@host# mkdir /usr/local/instrepo
root@host# cp -ar /mnt/. /usr/local/instrepo
```

Create repo XML info to that directory

```
root@host# createrepo -v /usr/local/instrepo/Packages
```

Create a repo file under /etc/yum.repos.d/

```
root@host# vi /etc/yum.repos.d/instrepo.repo
  [instrepo]
  name=Install Source Repo
  baseurl=file:///usr/local/instrepo/Packages
  enabled=1
  gpgcheck=0
```

### Share via HTTP

You can create a repo for all your servers to share - Steps are pretty much the same. Except you put your RPMs in the Apache directory.

Copy DVD contents to the Apache directory

```
root@host# mount rhel6-x86_64.iso /mnt -o loop
root@host# mkdir /var/www/html/inst
root@host# cp -ar /mnt/. /var/www/html/inst
```

Create repo XML info to that directory

```
root@host# createrepo -v /var/www/html/inst
```

Start Apache

```
root@host# service httpd start
root@host# chkconfig httpd on
```

On Each Server - Create a repo file under /etc/yum.repos.d/

```
root@host# vi /etc/yum.repos.d/instrepo.repo
  [instrepo]
  name=Install Source Repo
  baseurl=http://inst.example.com/inst
  enabled=1
  gpgcheck=0
```

### Repo sync

You can sync repos you're subscribed to locally. Here is the QnD

```
for repo in rhel-7-server-rpms rhel-7-server-extras-rpms
do
  reposync --gpgcheck -lm --repoid=${repo} --download_path=/var/www/html/repos
  createrepo -v /var/www/html/repos/${repo} -o /var/www/html/repos/${repo}
done
/sbin/restorecon -vR /var/www/html
```

The `my.repo` file would look something like this

```
[rhel-7-server-rpms]
name=rhel-7-server-rpms
baseurl=http://inst.example.com/rhel-7-server-rpms
enabled=1
gpgcheck=0
[rhel-7-server-extras-rpms]
name=rhel-7-server-extras-rpms
baseurl=http://inst.example.com/rhel-7-server-extras-rpms
enabled=1
gpgcheck=0

```
