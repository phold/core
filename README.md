# Phold - Core

The core is split in two projects: the site administration (starting/stopping the server, adding a new website to the proxy, etc) and the api that allows all of that happen from a client.

## Phold workflow 

Website Creation
================

A user access one of our clients, say, the http one. He logs in, and creates a website (you can create as much websites as you like). He then proceeds to configure it (we copy this configuration back to a Hugo backend) and selecting a theme, uploading images, altering templates, writing content, etc.

After is all said and done, the user clicks a "publish" button.

**Now, this is where the fun begins**

The idea is to take all of the previous steps (i.e. the assembled website, created and curated by Hugo and already static) and make a .tar.gz out of it. We then pick one of our many machines available, and copy via ssh this .tar.gz to it.

We then move over to the machine we copied the compiled website to, unpack it to a directory, and flush out a standard Dockerfile that runs our website. We start the website on this container and assure that everything is OK.

After that, the machine has webserver that proxies requests from port 80 to different backend webservers. This needs to be updated with the port exposed by the Docker container we just created. This involves a graceful restart of the proxy in order to include the new website.

And there you have it. The docker container is serving a website, and being proxied from the internet.

Website Update
==============

The user can still alter all the details of the website he created. The workflow is almost the same as above, but involves tearing down the Docker container associated with his website, altering the website (just overwrite it, but make sure the user understands this) and proxy reloading.

Website Deletion
================

The user can delete the website at any given point in time. We take the metadata relating to the Docker container, tear it down and remove the website from the proxy.

## Notes

 - Should we version websites and allow for easy rollback (maybe in the paid version)?
