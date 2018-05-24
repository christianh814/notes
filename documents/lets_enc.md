# Let's Encrypt

This  will show you how to set up a TLS/SSL certificate from [Let's Encrypt](http://letsencrypt.org/) to use for webservers/websreveces. Additionally, we will cover how to automate the certificate renewal process using a cron job.

SSL certificates are used within web servers to encrypt the traffic between server and client, providing extra security for users accessing your application. Let’s Encrypt provides an easy way to obtain and install trusted certificates for free.

For this purposes we will be using `feodra` and the certification automation tool called `certbot`

Also `acme.sh` is available.

* [Installation](#installation)
* [Request Certificate](#request-certificate)
* [Renewal](#renewal)

## Installation

First, make sure the system is updated

```
dnf -y clean all
dnf -y update
```

Next, install `certbot`

```
dnf -y install certbot
```

For EL7 users, you might need to do this

```
subscription-manager repos --disable=*
subscription-manager repos --enable="rhel-7-server-rpms" \
--enable="rhel-7-server-extras-rpms" --enable=rhel-7-server-optional-rpms
yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
```

The help menu should be enough to get you started

```
$ certbot --help

  certbot [SUBCOMMAND] [options] [-d DOMAIN] [-d DOMAIN] ...

Certbot can obtain and install HTTPS/TLS/SSL certificates.  By default,
it will attempt to use a webserver both for obtaining and installing the
cert. The most common SUBCOMMANDS and flags are:

obtain, install, and renew certificates:
    (default) run   Obtain & install a cert in your current webserver
    certonly        Obtain or renew a cert, but do not install it
    renew           Renew all previously obtained certs that are near expiry
   -d DOMAINS       Comma-separated list of domains to obtain a cert for

  (the certbot apache plugin is not installed)
  --standalone      Run a standalone webserver for authentication
  (the certbot nginx plugin is not installed)
  --webroot         Place files in a server's webroot folder for authentication
  --manual          Obtain certs interactively, or using shell script hooks

   -n               Run non-interactively
  --test-cert       Obtain a test cert from a staging server
  --dry-run         Test "renew" or "certonly" without saving any certs to disk

manage certificates:
    certificates    Display information about certs you have from Certbot
    revoke          Revoke a certificate (supply --cert-path)
    delete          Delete a certificate

manage your account with Let's Encrypt:
    register        Create a Let's Encrypt ACME account
  --agree-tos       Agree to the ACME server's Subscriber Agreement
   -m EMAIL         Email address for important account notifications

More detailed help:

  -h, --help [TOPIC]    print this message, or detailed help on a topic;
                        the available TOPICS are:

   all, automation, commands, paths, security, testing, or any of the
   subcommands or plugins (certonly, renew, install, register, nginx,
   apache, standalone, webroot, etc.)
```

## Request Certificate

First; you need to allow `80` and `443` in your firewall config

```
firewall-cmd --add-service=http
firewall-cmd --add-service=https
firewall-cmd --runtime-to-permanent
```

Simply pass `-d domain.tld` to request for a certificate (for multiple domains pass multiple `-d`)

```
certbot certonly --standalone -d example.com --agree-tos -m example@example.com
```

Note that `example.com` must be resolve-able

Your certificates (along with a README) will be placed under:

```
/etc/letsencrypt/live/<your domain>
```

## Renewal

Let's Encrypt certs are temporary; you can cron the recreation with...

```
15 1 * * 7 /usr/bin/certbot renew --renew-hook "systemctl restart httpd" >> /var/log/le-renew.log
```

Certs expire every 90 days or so...so I set it to run every other week. It's probably better to create a wrapper script for this (in case you need to do a "song and dance" like for OpenShift)
