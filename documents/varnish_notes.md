# Varnish Notes

Varnish Notes in no paticular order

* [Installation](#installation)
* [Configuration](#configuration)

There are good resources

* [Drupal Example](http://www.lullabot.com/articles/varnish-multiple-web-servers-drupal)
* [Officail Docs](https://www.varnish-cache.org/docs/2.1/tutorial/advanced_backend_servers.html)

## Installation

[Install instructions](https://www.varnish-cache.org/installation/redhat)

Install the repo
```
root@host# rpm --nosignature -i http://repo.varnish-cache.org/redhat/varnish-3.0/el5/noarch/varnish-release-3.0-1.noarch.rpm
```

Install YUM priorities

```
root@host# yum -y install yum-plugin-priorities
```

Make sure the ` /etc/yum.repos.d/varnish.repo ` file has a high priority (in this case ` priority=1 `)

```
[varnish-3.0]
name=Varnish 3.0 for Enterprise Linux 5 - $basearch
baseurl=http://repo.varnish-cache.org/redhat/varnish-3.0/el5/$basearch
priority=1
enabled=1
gpgcheck=0
#gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-VARNISH
```

Use YUM to install

```
root@host# yum -y install varnish
```

## Configuration

Edit the ` /etc/sysconfig/varnish ` file and change the "listen" port to 80

```
VARNISH_LISTEN_PORT=80
```

Next create a VLC file under ` /etc/varnish/ `

Sample ` default.vlc ` file
```
# Define the list of backends (web servers).
# Port 80 Backend Servers

#
# "backend" is your host declarations. You set up
#  the host IP address the url to 'curl' for health
#  check and how often. You also set up the threshold
#  until a host is  declared 'down'
#

# mm.la3.4over.com
backend mm_la3 {
	.host = "192.168.114.144";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

# nn.la3.4over.com
backend nn_la3 {
	.host = "192.168.114.145";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

# qq.la3.4over.com
backend qq_la3 {
	.host = "192.168.114.115";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

# rr.la3.4over.com
backend rr_la3 {
	.host = "192.168.114.116";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

# ss.la3.4over.com
backend ss_la3 {
	.host = "192.168.114.118";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

# tt.la3.4over.com
backend tt_la3 {
	.host = "192.168.114.119";
	.probe = {
		.url = "/health.php";
		.interval = 5s;
		.timeout = 3s;
		.window = 5;
		.threshold = 3;
	 }
	.first_byte_timeout = 300s;
}

#
# Define the "director" that determines how to distribute incoming requests.
# it tells what load-balancing method to use and what the 'backends' are
#

# Production Fundink Director
director prod_fundink round-robin {
  { .backend = qq_la3; }
  { .backend = rr_la3; }
}

# Staging Fundink Director
director staging_fundink round-robin {
  { .backend = mm_la3; }
  { .backend = nn_la3; }
}

# Holding Fundink Director
director holding_fundink round-robin {
  { .backend = ss_la3; }
  { .backend = tt_la3; }
}

#
# This sets vcl_recv - so what do do with the request.
#
sub vcl_recv {
# choose a backend depending on domain
        if (req.http.host ~ "^(www\.|admin\.|api\.)?fundink.com$") {
                set req.backend = prod_fundink;
        }
	else if (req.http.host ~ "^(.*)?staging.(.*)?fundink.com$") {
                set req.backend = staging_fundink;
        }
        else if (req.http.host ~ "^holding.fundink.com$") {
                set req.backend = holding_fundink;
        }
	else {
              	error 418;
        }
        # remove and reset the clients IP address
	# So that Apache sees the client adress and not
	# varnish's address
        remove req.http.X-Forwarded-For;
        set req.http.X-Forwarded-For=client.ip;
}
#
#-30-
```

Test the config file

```
root@host# varnishd -C -f /etc/varnish/default.vlc
```

It should spit out your config (good) or throw an error (bad)

If you're all set go ahead and restart

```
root@host# service varnish restart
```
