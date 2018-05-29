# Nginx 

Nginx is a webserver like apache. 

Good Resources can be found [here](http://www.lifelinux.com/how-to-install-nginx-and-php-fpm-on-centos-6-via-yum) and [here](http://www.howtoforge.com/installing-nginx-with-php5-and-php-fpm-and-mysql-support-on-centos-6.3-p2)

* [Installation](#installation)
* [Configure PHP FPM](#configure-php-fpm)
* [PHP FPM Unix Socket](#php-fpm-unix-socket)
* [Nginx Virtual Hosting Configuration](#nginx-virtual-hosting-configuration)

## Installation

Nginx is in the repos so a straight YUM command will do

```
root@host# yum -y install nginx
```

Make sure it starts on boot

```
root@host# chkconfig nginx on
```

You can start putting HTML files under...

```
/usr/share/nginx/html/
```

## Configure PHP FPM

First install `php-fpm` (in most repos - might be different for EL)

```
root@host# yum -y install php php-fpm
```

`PHP-FPM` is a service so make sure it starts on boot

```
root@host# chkconfig php-fpm on
root@host# service php-fpm start
```

Edit the `/etc/nginx/nginx.conf` file

```
root@host# vi /etc/nginx/nginx.conf
```

And make sure it looks something like this (this is a snippet)

NOTE: The "worker process" should = number of sockets

```
      ...
      worker_processes  4;
      ...
	location / {
            root   /usr/share/nginx/html;
            index  index.html index.htm index.php;
        }
        ...
	location ~ \.php$ {
            root           html;
            fastcgi_pass   127.0.0.1:9000;
            fastcgi_index  index.php;
            fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
            include        fastcgi_params;
        }
        ...
```

Reload nginx

```
root@host# service nginx reload
```

You can test this by creating `/usr/share/nginx/html/info.php` with the following...

```
<?php
phpinfo();
?>
```

## PHP FPM Unix Socket
By default PHP-FPM is listening on port 9000 on 127.0.0.1. It is also possible to make PHP-FPM use a Unix socket which avoids the TCP overhead.

To do this, open `/etc/php-fpm.d/www.conf` and change the following

```
;listen = 127.0.0.1:9000
listen = /dev/shm/php-fpm.sock
```

Now restart `PHP-FPM`

```
root@host# service php-fpm reload
```

Next, you need to tell nginx to use the socket - change the following line (on all your vhosts) to look like this...

```
fastcgi_pass   unix:/dev/shm/php-fpm.sock;
```

Then restart nginx

```
root@host# service nginx restart
```

## Nginx Virtual Hosting Configuration

Creating virtual host config file

```
root@host# cd /etc/nginx/conf.d/
root@host# cp virtual.conf site.conf
```

Open site.conf, enter

```
server {
        server_name  domain.local;
        root /home/www/domain.local/public_html;
        access_log /home/www/domain.local/log/domain.local-access.log;
        error_log /home/www/domain.local/log/domain.local-error.log;

        location / {
                index  index.html index.htm index.php;
        }
        location ~ \.php$ {
                include /etc/nginx/fastcgi_params;
                fastcgi_pass   unix:/dev/shm/php-fpm.sock;
                fastcgi_index index.php;
                fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        }
}
```

Test the config before restarting

```
root@host# service nginx configtest
```


