# PHP Composer

PHP Has a method to install 3rd party libraries called "Composer".

I installed it on my fedora box using

```
dnf -y install composer
```

Next, cd into your project's dir and run `composer init` to create the json file
```
composer init
```

The questions are pretty self explanatory; but here's an example of my answers
```
[root@fed-webserver cms]# composer init
Do not run Composer as root/super user! See https://getcomposer.org/root for details

                                            
  Welcome to the Composer config generator  
                                            


This command will guide you through creating your composer.json config.

Package name (<vendor>/<name>) [root/mycms]: 
Description []: 
Author [, n to skip]: Christian Hernandez <c_hernand1982@yahoo.com>
Minimum Stability []: 
Package Type (e.g. library, project, metapackage, composer-plugin) []: 
License []: 

Define your dependencies.

Would you like to define your dependencies (require) interactively [yes]? yes
Search for a package: phpmailer

Found 15 packages matching phpmailer

   [0] phpmailer/phpmailer
   [1] kruisdraad/phpmailer
   [2] phpmailerflamin/phpmailer
   [3] gulltour/phpmailer
   [4] opencart-patches/phpmailer
   [5] minmb/phpmailer
   [6] phpmailer/phpmailer
   [7] nabu-3/provider-phpmailer-phpmailer
   [8] blackleg/elgg-phpmailer
   [9] zyx/zyx-phpmailer
  [10] qu-modules/qu-phpmailer
  [11] geolysis/silverstripe-phpmailer
  [12] codegun/phpmailer-lite
  [13] adrianorsouza/codeigniter-phpmailer
  [14] locomotivemtl/charcoal-email

Enter package # to add, or the complete package name if it is not listed: 0
Enter the version constraint to require (or leave blank to use the latest version): 

Using version ^5.2 for phpmailer/phpmailer
Search for a package: Would you like to define your dev dependencies (require-dev) interactively [yes]? no

{
    "name": "root/mycms",
    "require": {
        "phpmailer/phpmailer": "^5.2"
    },
    "authors": [
        {
            "name": "Christian Hernandez",
            "email": "c_hernand1982@yahoo.com"
        }
    ]
}

Do you confirm generation [yes]? yes
Would you like the vendor directory added to your .gitignore [yes]? 
```

In the end you'll have a file called `composer.json` in the root of your project. Mine looks like this
```
{
    "name": "root/mycms",
    "require": {
        "phpmailer/phpmailer": "^5.2"
    },
    "authors": [
        {
            "name": "Christian Hernandez",
            "email": "c_hernand1982@yahoo.com"
        }
    ]
}

```

To install this  package, just run `composer install`

```
[root@fed-webserver cms]# composer install
Do not run Composer as root/super user! See https://getcomposer.org/root for details
Loading composer repositories with package information
Updating dependencies (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing phpmailer/phpmailer (v5.2.24): Downloading (100%)         
phpmailer/phpmailer suggests installing league/oauth2-google (Needed for Google XOAUTH2 authentication)
Writing lock file
Generating autoload files
```

This creates a directory called `vendor` with all your libraries installed.

```
[root@fed-webserver cms]# ll vendor/
total 4
-rw-r--r--. 1 root root 178 Aug 24 12:53 autoload.php
drwxr-xr-x. 2 root root 203 Aug 24 12:53 composer
drwxr-xr-x. 3 root root  23 Aug 24 12:53 phpmailer
```

Load these librarys with either

```
require __DIR__ . '/vendor/autoload.php';
```

or 

```
require 'vendor/autoload.php';
```

If you add packages to your `composer.json` file (like this)
```
{
    "name": "root/mycms",
    "require": {
        "phpmailer/phpmailer": "^5.2",
        "sendgrid/sendgrid": "~6.0"
    },
    "authors": [
        {
            "name": "Christian Hernandez",
            "email": "c_hernand1982@yahoo.com"
        }
    ]
}
```

Run the following to install them

```
composer update
```

Here is some good info on `composer`:

  * `composer install` - installs the vendor packages according to composer.lock (or creates composer.lock if not present),
  * `composer update` - always regenerates composer.lock and installs the lastest versions of available packages based on composer.json
  * `composer dump-autoload` - won’t download a thing. It just regenerates the list of all classes that need to be included in the project (autoload_classmap.php). Ideal for when you have a new class inside your project.
  * Ideally, you execute `composer dump-autoload -o` - a faster load of your webpages. The only reason it is not default, is because it takes a bit longer to generate (but is only slightly noticable)

# PHP Passwords

You can use builtin php functions to encrypt passwords

# Encrypt

To encrypt the password do the following (assign this to a variable)
```
password_hash($password_given_in_post, PASSWORD_BCRYPT, array('cost' => 12));
```

# Validate

Validate it with the following command (after extracting the saved password in the db)
```
password_verify($password_given_in_post, $enc_password_in_db);
```
