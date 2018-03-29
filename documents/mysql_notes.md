# MySQL Notes

Following is notes for MySQL in no paticular order

* [Create Database With User](#create-database-with-user)
* [MySQL Sample Backup and Recovery](#mysql-sample-backup-and-recovery)
* [MySQL Misc Commands](#mysql-misc-commands)
* [MariaDB Notes](#mariadb-notes)

## Create Database With User

Below is an example of what I needed to do (after you connect).

```
root@host# mysql -u root -p
  mysql> create database eticket;
  mysql> connect eticket;
  mysql> GRANT ALL PRIVILEGES ON `eticket` . * TO 'eticket'@'localhost' identified by 'eticket' WITH GRANT OPTION;
  mysql> use eticket
  mysql> quit;
```

You can change `localhost` with an IP or DNS name if it's a remote machine. To allow "all" hosts...change it to ` '%'`; example:

```
GRANT ALL PRIVILEGES ON `ecom`.* TO 'ecom'@'%' IDENTIFIED BY  'Ecom123456' WITH GRANT OPTION;
```

## MySQL Sample Backup and Recovery

Below is an example of a backup and recovery procedure. This is not a "comprehensive" how to; this is just for reference. You can do this manually or you can cron it.

To backup all databases use the mysqldump command with the `-A` or the `--all-databases` option (both do the same thing). The "password" is the password for the database super-user:

```
root@host# mysqldump --user=root --password=password -A > mysql_bkup.sql
```

To recover from this backup use the mysql command and "feed in" the mysql_bkup.sql file.

```
root@host# mysql --user=root --password=password < mysql_bkup.sql
```

Now if you just want to backup a specific database and not all databases you can still use the mysqldump command; all you need to do is specify the database name. This command also requires the database "root" account. NOTE: The database "root" account is different from the root account on the system; they are most likley be different (they should be anyway).

To backup a specific database use the `--opt` switch when specifying the database:

```
root@host# mysqldump --user=root --password=password --opt db_name >  db_name-bkup.sql
```

To recover the database from this backup use the mysql command and "feed in" the `db_name-bkup.sql` file.

```
root@host# mysql --user=root --password=password db_name <  db_name-bkup.sql
```

NOTE: Do Not Use the `--opt` option when doing a `-A` or the `--all-databases` options. The `--opt` loads the dump into memory before commiting to disk; so it can be a problem if dumping all databases because it can get big. 

## MySQL Misc Commands


To list what databases exist within a MySQL installation; use the show command within a MySQL prompt.

```
mysql> show databases;
```

Set up root password with...

```
mysqladmin -u root password NEWPASSWORD
```

However, if you want to change (or update) a root password, then you need to use the following command:

```
mysqladmin -u root -p 'oldpassword' password newpass
```

## MariaDB Notes

Following are quick notes for MariaDB. MariaDB is a "fork" of MySQL and, in principal, should be "backwards compatible"

* [Installation](#installation)

### Installation

Install the server and client tools

```
[root@sandbox ~]# yum -y install mariadb-client mariadb
```

Start and enable the system

```
[root@sandbox ~]# systemctl start mariadb
[root@sandbox ~]# systemctl enable mariadb
```

Verify that MariaDB is listening on all interfaces.

```
[root@sandbox ~]# ss -tulpn | grep mysql
tcp    LISTEN     0      50                     *:3306                  *:*      users:(("mysqld",3711,14))
```

Enable the `skip-networking` directive by opening the file `/etc/my.cnf` in a text editor, and in section `[mysqld]`, add the line `skip-networking=1`

```
[root@sandbox ~]# grep skip-networking /etc/my.cnf
skip-networking=1
```

Restart the service and make sure that it's not listening

```
[root@sandbox ~]# systemctl restart mariadb
[root@sandbox ~]# ss -tulpn | grep mysql
```

Secure the Mariadb service using the `mysql_secure_installation` tool.

```
[root@sandbox ~]# mysql_secure_installation 
/bin/mysql_secure_installation: line 379: find_mysql_client: command not found

NOTE: RUNNING ALL PARTS OF THIS SCRIPT IS RECOMMENDED FOR ALL MariaDB
      SERVERS IN PRODUCTION USE!  PLEASE READ EACH STEP CAREFULLY!

In order to log into MariaDB to secure it, we'll need the current
password for the root user.  If you've just installed MariaDB, and
you haven't set the root password yet, the password will be blank,
so you should just press enter here.

Enter current password for root (enter for none): 
OK, successfully used password, moving on...

Setting the root password ensures that nobody can log into the MariaDB
root user without the proper authorisation.

Set root password? [Y/n] Y
New password: ********
Re-enter new password: ******** 
Password updated successfully!
Reloading privilege tables..
 ... Success!


By default, a MariaDB installation has an anonymous user, allowing anyone
to log into MariaDB without having to have a user account created for
them.  This is intended only for testing, and to make the installation
go a bit smoother.  You should remove them before moving into a
production environment.

Remove anonymous users? [Y/n] Y
 ... Success!

Normally, root should only be allowed to connect from 'localhost'.  This
ensures that someone cannot guess at the root password from the network.

Disallow root login remotely? [Y/n] Y
 ... Success!

By default, MariaDB comes with a database named 'test' that anyone can
access.  This is also intended only for testing, and should be removed
before moving into a production environment.

Remove test database and access to it? [Y/n] Y
 - Dropping test database...
 ... Success!
 - Removing privileges on test database...
 ... Success!

Reloading the privilege tables will ensure that all changes made so far
will take effect immediately.

Reload privilege tables now? [Y/n] Y
 ... Success!

Cleaning up...

All done!  If you've completed all of the above steps, your MariaDB
installation should now be secure.

Thanks for using MariaDB!

```
