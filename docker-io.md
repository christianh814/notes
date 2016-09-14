# Overview

This is a 100K view of what docker does. 

From a high level; docker is a way to run applications within LXC containers (much like Solaris Zones). They are not VMS in a technical sense...but you can think of them as a VERY VERY thin VM

First install docker (from EL, you need EPEL, on Fedora it's in the repos)

	
	# yum -y install docker-io


Next, pull a base image

	
	# docker pull centos


Test docker by running hello world 

	
	docker run centos /bin/echo "hello world"
	hello world


What happened here is that you ran the command "hello world" within a container (not the server itself...cool right?)

You can also run an interactive shell inside

	
	docker run -i -t centos /bin/bash


Now you can build an image based on the centos base image. Let's set up an apache instance running php

First create a // Dockerfile //

	
	# DOCKER-VERSION 0.7.2
	
	FROM    centos:6.4
	
	RUN rpm -Uvh http://download.fedoraproject.org/pub/epel/6/i386/epel-release-6-8.noarch.rpm
	
	RUN yum -y install httpd php php-mysql
	
	EXPOSE  80
	ADD index.php /var/www/html/
	ENTRYPOINT ["/usr/sbin/httpd"]
	CMD ["-D", "FOREGROUND"]


Let's break this down...

*  // DOCKER-VERSION // - Is the minimum version this image needs

*  // FROM // - Is the base image you're going to base your container on

*  // RUN // - These are commands you are running during build time

*  // EXPOSE // - Here, you can "expose" a port to the hosting server. (i.e. I "open" port 80 from to container to the housing server)

*  // ADD // - Push files and/or directories to the container (think of it as a CP/MV command). The first option is relative to the server the second is relative to the container (i.e. server:/path/to/file container:/path/to/file)

*  // ENTRYPOINT // - The "magic" of the container. The "entrypoint" allows you to specify a command to run as the default entrypoint for a container. Great explanation can be found [HERE](http://www.kstaken.com/blog/2013/07/06/how-to-use-entrypoint-in-a-dockerfile/)

*  // CMD // - command to run against the entry point

Now build the container (make sure that the // index.php // file exists.)

	
	# docker build -t chrish/centoshttpd .


Once that's done; list your images...

	
	# docker  images
	REPOSITORY           TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
	chrish/centoshttpd   latest              c26790d793fa        30 minutes ago      417 MB
	centos               6.4                 539c0211cd76        9 months ago        300.6 MB
	centos               latest              539c0211cd76        9 months ago        300.6 MB


Now that you see that...run the container...(it'll spit out a UUID if successful)

	
	# docker run -d -p 8080:80 chrish/centoshttpd
	d405a4e7c8ddfe4e03460c2d3a05c3589ad16c68520939812246fd7765fc116c


Let's break this down

*  // docker run // - Command to run containers/images

*  // -d // - This is "detached" mode...the container runs in the background

*  // -p 8080:80 // - This is port mapping. You are taking localport 8080 and binding it to the container's port 80

*  // chrish/centoshttpd // - The "TAG" of the container running (you could have also used the "IMAGE ID")

Test this by going to your server's port 8080 and see the fruits of your labor

Now enter your container

	
	docker exec -it [container-id] bash


For more information visit [Docker's Documentation](http://docs.docker.io) or [Red Hat's How To](https///access.redhat.com/articles/881893)
