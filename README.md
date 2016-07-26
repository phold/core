# Core

The core is split in two projects: the site administration (starting/stopping the server, adding a new website to the proxy, etc) and the api that allows all of that happen from a client.

# Phold workflow 

Website Creation
================

A user access one of our clients, say, the http one. He logs in, and creates a website (you can create as much websites as you like). He then proceeds to configure it (we copy this configuration back to a Hugo backend) and selecting a theme, uploading images, altering templates, writing content, etc.

After is all said and done, the user clicks a "publish" button.

**Now, this is where the fun begins**

The idea is to take all of the previous steps (i.e. the assembled website, created and curated by Hugo and already static) and make a .tar.gz out of it. We then pick one of our many machines available, and copy via ssh this .tar.gz to it.

We then move over to the machine we copied the compiled website to, unpack it to a directory, and flush out a standard Dockerfile that runs our website. We start the website on this container and assure that everything is OK.

After that, the machine has webserver that proxies requests from port 80 to different backend webservers. This needs to be updated with the port exposed by the Docker container we just created. This involves a graceful restart of the proxy in order to include the new website.

And there you have it. The docker container is serving a website, and being proxied from the internet.

The name of the website is configurable on the web application. It will be a subdomain of `phold`. So, if you choose the name `correct-horse-battery-staple`, the name of your website shall be `correct-horse-battery-staple.phold.io`. The availability of such subdomains will be delegated to Phold. Two subdomains must not clash.

Website Update
==============

The user can still alter all the details of the website he created. The workflow is almost the same as above, but involves tearing down the Docker container associated with his website, altering the website (just overwrite it, but make sure the user understands this), creating a new container (of course, the image *will* be cached) and proxy reloading.

Website Deletion
================

The user can delete the website at any given point in time. We take the metadata relating to the Docker container, tear it down and remove the website from the proxy.

Migrations
==========

We will support migrations that the Hugo tools support: Jekyll, WordPress, Blogger, etc.

Command Line
============

Since the web is nothing but a client to an external API, we will have a command line client to upload websites.

Much like the web uses Hugo behind the scenes, we'll have a route that receives a compiled website from Hugo and host it, too.

For this to work, an API key is going to be needed, and the user will have to grab it from the Phold website. Also, email validation happens in this step.

Features
========

The free version will have a wide array of features, including:

- The ability to create a website from the webpage **and** putting it online in less than 2 minutes (like the Hugo video)
- Making changes to that website and have it updated virtually immediately
- The ability to provide a Hugo project and put it online
- The ability to migrate from WordPress/Blogger/Jekyll to our website
- The ability to download the Hugo-compiled website (but not the Hugo scaffolding - this is going to be paid)
- Command line interface for superusers that don't want to host any content online.
  - With this, you can create a hugo project, compile it, and put it up via command line with a name. This will spin up a Docker container to host your website, but there will be no Hugo information for your website, only metadata regarding the machine and container PID. Updates don't exist here either: every time you publish a website, it will overwrite your previous one (tearing down the Docker container and putting it back up with the new website). Also, we are going to need an API key control for this to work. That way, we'll identify the user uploading websites (in case someone exploits a breach in our systems)


## Notes

- Discuss interesting features available for the paid version of the web client:
 - Fully qualified hostnames
 - Especialized support
 - Website download as a Hugo project
