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
- Portainer (optional)

## Why Portainer?
I use Portainer to manage my Docker containers. It is not required, but it makes it easier to manage them. You can
easily add each stack directly from repository and have it automatically updated at chosen interval.

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
