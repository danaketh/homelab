# HomeLab
This repository contains all the information about my homelab. I will try to keep it up to date as much as possible.
If you have any questions, feel free to contact me.

You can use this as a reference for your own homelab, but please keep in mind that this is a hobby project and not
a production environment. I am not responsible for any damage or data loss.4

Since my homelab is running in local network, most of the containers use simple passwords and are not secured. If you
plan to use any of these containers in production, please make sure to:
- use strong passwords
- use HTTPS
- use firewall to limit access to the containers
- make sure you DO NOT store passwords in the repository

## Requirements
- some HW to run 24/7
- some Linux to run on the HW
- some knowledge about Docker
- Pulumi

## Why Pulumi?
Because it's the real infrastructure as code. You're not forced to learn some weird language that is used only for
infrastructure. In this case, I'm using Python, but Pulumi also support TypeScript, JavaScript, Go, C#, Java, and YAML.

### How to use it?
Let's assume you have [installed Pulumi](https://www.pulumi.com/docs/install/) on your workstation and you have also installed Docker on the machine where
your containers will be running. We'll have to do a tiny bit of configuration before we can start deploying tho.
I'm using Ubuntu, so if you're one something else, you'll have to figure out some of the steps yourself.

Login to your server console and run:
```bash
sudo systemctl edit docker.service
```

This should open an editor for you to modify the startup script for Docker. Ignore the long comment in the file
where everything looks like you can just uncomment it. Instead, look for the `### Anything between here and the comment below will become the new contents of the file` and `### Lines below this comment will be discarded` comments, and add the following lines between them:

```bash
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock -H tcp://0.0.0.0:2375
```

This will open Docker API to the world. In case of our homelab, it's not a big deal, but if you're going to use
this repository as a boilerplate for your production environment, you should do with further configuration
and [add TLS](https://linuxhandbook.com/docker-remote-access/). You may also go above and beyond and use some
firewall to limit access to the API either to a list of IP addresses, a VPN, or some kind of authentication.

Save the file and restart Docker:
```bash
sudo systemctl restart docker.service
```

You can verify that Docker API is now listening on port 2375:
```bash
sudo lsof -i -P -n | grep LISTEN | grep 2375
```

This should return something like this:
```text
dockerd   1727     root    3u  IPv6  28222      0t0  TCP *:2375 (LISTEN)
```

And you're done with the server configuration. Now you can clone this repository to your workstation and from there,
please refer to the [Pulumi documentation](https://www.pulumi.com/docs/cli/) on how to use Pulumi.

## Available stacks

### AdGuard
AdGuard is a DNS server with ad blocking capabilities. Since I have quite a number of devices at home, I decided to
use it as a DNS server for my network. Blocking ads on mobile devices is a bit tricky, so using AdGuard as a DNS
allows me to protect kids from unwanted content. I'm also using it for forwarding DNS requests to my local services.

### Domoticz
Domoticz is a home automation system. I'm using it to control my lights, heating, etc. Just a hobby project to turn
my home into a smart home without spending a lot of money and using proprietary solutions.

### Dozzle
Dozzle is a simple container log viewer for Docker. You get all logs in one place and you can easily filter them.

### Gitea
Gitea is a self-hosted Git service. I'm using it to test workflows before pushing them to GitHub and also as a backup
of all my repositories.

### Grafana
Grafana is a monitoring platform. I'm mostly just playing with it to learn something new and play with different data.

### Traefik
Traefik is a reverse proxy. It allows me to expose my services to the network without exposing them directly and having
to use different ports for each service. Instead I can just use domain names and Traefik will take care of the rest.

### Watchtower
Watchtower is a container that automatically updates other containers. Because I'm lazy!
