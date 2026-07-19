---
title: Migrating to Arcane
summary: I needed to move my Minecraft servers anyway so I switch from Puffer Panel to Arcane
tags:
  - selfhosting
  - homelab
  - xcp-ng
  - projects
  - virtualization
  - proxmox
  - containers
  - podman
  - docker
author: Ethan Ashley
date: 2026-07-18
created: 2026-07-18
updated: 2026-07-19
draft: false
publish: true
---

So as a part of my home lab, as many do, I host game servers. The game servers relevant here are my two Minecraft servers.
Yesterday i had that server was not performing how I wanted it to. I had it running on an old server and it just doesn't have the power for this kind of thing. So it was time to migrate to something else.

# The old setup
4C4T CPU 
16GB RAM
Proxmox installed

The Minecraft servers are running in Docker containers managed by [Puffer Panel](https://www.pufferpanel.com/) inside an LXC container on Proxmox.
Besides from not having enough power now, this is a unique configuration in my environment. In this same container, there are other Docker containers running, not managed by Puffer Panel.

## Background

I introduced Puffer Panel because I wanted some UI and an easier way to manage my game servers. If I need to change something about one of my services, I don't mind having to be in the CLI, as that is the purpose of that time. If I have to change something for a game server, it's usually a quick change or because it is not working, and I don't want to always have to pull out my laptop and hop in a terminal for that, and that's not how I want to be spending my time at that moment, I want to be playing the game.


![[attachments/e9f1cd8b.png]]


Puffer Panel works well. I chose it because [Pterodactyl](https://pterodactyl.io/) was overkill, had too much overhead, and seemed like it would possibly become more of a headache than it was worth.
I moved to Puffer Panel after running my game servers the old fashion way and eventually begin to feel the drawbacks that comes with.
With Puffer Panel, I can manage and edit files, view load, see the logs and send commands, easily change settings, make backups, and give other users access.
If I didn't have any other reason to migrate away from this system, I would still be using Puffer Panel, but it's time for the new shiny.

I have also had a new shiny thing in the back of my head since i heard about it on the [selfh.st blog](https://selfh.st/post/2025-favorite-new-apps/) and later on the [ServersatHome Arcane vs Dockhand video](https://www.youtube.com/watch?v=3TaDWpYgGtE). 
I've been using Docker for 7 years, starting on LXC 9 years ago (as of 2026). When I last looked at Docker GUIs, they weren't anything crazy, and they surely were not a stable complete solution. Since I already have structure to how I manage my compose files and am comfortable with using Docker directly, there was no point. But, with this problem these game servers have introduced, there is now.

The new shiny is Arcane. Similar to Puffer Panel, it has templates, but these templates are much more native to how Docker (or others) work. There are not many game templates in the community repo, but templates are easy to make.
In addition, Arcane is for any container services, not just game servers. 
It supports all the same things Puffer Panel does: manage, edit, and backup files; see load, logs, and configs; send commands; perform updates; create new servers.
Arcane also supports a couple of other things: Docker Swarm, Podman (Beta), Git, automatic updates, vulnerability scanning, and OIDC.

# The new setup
2 4C8T CPUs
96GB RAM

This server is running XCP-NG with a twin node which it replicates to for disaster recovery, and the core design of these two systems is mobility, as in, I can move it elsewhere or use a different WAN and things keep working.
To make a new server for this, it was a few clicks to create a new VM from my Alma Linux template, attaching my on boarding cloud-init config, and sending it. After a while, I had a fresh new VM ready for services.


![[attachments/d4f32be5.png]]

## Arcane

Setting up Arcane was easy too: just a few commands and a compose file and I was off
[Installation](https://getarcane.app/docs/setup/installation)

Arcane was built for Docker first, and Podman support is still in beta, but we try new things around here, and I prefer Podman these days.

```
sudo dnf install podman podman-compose podman-docker 
sudo dnf install zip unzip  # if you need them
systemctl --user enable --now podman.socket
loginctl enable-linger $UID
```

`compose.yml`
```
services:
  arcane:
    image: ghcr.io/getarcaneapp/manager:latest
    container_name: arcane
    ports:
      - '3552:3552'
    volumes:
      - /run/user/1000/podman/podman.sock:/var/run/docker.sock
      - ./data:/app/data
      - ./builds:/builds
      - ./backups:/backups
      # Optional host project mount:
      # - /path/to/projects:/app/data/projects
    security_opt:
      - label:disable
    environment:
      - APP_URL=http://localhost:3552
      - PUID=1000
      - PGID=1000
      - ENCRYPTION_KEY=REDACT
      - JWT_SECRET=REDACT
    restart: always
```

Then start it with `podman compose up` and the setup password will be in the console.

> [!NOTE]- Tangent for developers
> I would like to take a tangent though: having setup passwords with static credentials are still a security vulnerability. A random password should be generated instead.
> Developers would be surprised by how many default instances are left exposed on the internet, through carelessness, forgetfulness, or a misplaced `restart: always`.
> Use random passwords.


![[attachments/ad454d77.png]]

![[attachments/f5b83a79.png]]

There are a couple of patches I had to make, more than likely due to using Podman instead. Again, it's in beta so this is to be expected.
To see when SELinux is causing the issues: `sudo ausearch -m avc -ts recent`.
	Don't be lazy and disable SELinux, learn to use it instead.

`/etc/systemd/system/user@.service.d/podman.conf`
```
[Service]
Delegate=cpu cpuset io memory pids

[Slice]
TasksMax=80%
```

```
sudo systemctl daemon-reload
```

```
sudo semanage fcontext -a -t container_var_run_t "/var/run/user/1000/podman(/.*)?"
sudo restorecon -Rv /run/user/1000/podman

ls -Zd /run/user/1000/podman/podman.sock

sudo semanage fcontext -a -t container_file_t "/home/local-admin/arcane/backups(/.*)?"
sudo restorecon -Rv /home/local-admin/arcane/backups

sudo semanage fcontext -a -t container_file_t "/home/local-admin/arcane/builds(/.*)?"
sudo restorecon -Rv /home/local-admin/arcane/builds

sudo semanage fcontext -a -t container_file_t "/home/local-admin/arcane/data(/.*)?"
sudo restorecon -Rv /home/local-admin/arcane/data

sudo semanage fcontext -a -t container_file_t "/home/local-admin/.local/share/containers/storage/volumes(/.*)?"
sudo restorecon -Rv /home/local-admin/.local/share/containers/storage/volumes
```

## Minecraft
Now, the only thing I am missing is how I am going to deploy my Minecraft servers. During the search, I found this: [Minecraft Server Configurator for Docker](https://setupmc.com) which uses [this container](https://github.com/itzg/docker-minecraft-server). Configurator was useful for getting started, but it doesn't have [all of the options.](https://github.com/itzg/docker-minecraft-server/tree/master/docs/configuration) 
Now this container is impressive, it can do more server configurations than I knew about, and the configurator makes a stupid easy.

So now we have our compose, but this isn't very templatable for Arcane, so this is what I ended up with:
`compose.yaml`
```yaml
services:
  mc:
    container_name: ${NAME}
    image: itzg/minecraft-server:${TAG}
    restart: unless-stopped
    tty: true
    stdin_open: true
    ports:
      - "${PORT}:25565"
    volumes:
      - data:/data:Z
      
volumes:
  data:
    name: "${NAME}-data"
```
`.env`
```env
NAME="mc-frens"
CREATE_CONSOLE_IN_PIPE=true
TAG="stable"
PORT="25565"

EULA="TRUE"
VERSION="LATEST"
INIT_MEMORY="1G"
MAX_MEMORY="6G"
MAX_PLAYERS="10"
MOTD="Ethan Minecraft"
ICON="https://ethanashley.net/static/img/duck_mini.png"
USE_AIKAR_FLAGS="true"
USE_MEOWICE_FLAGS="true"
TZ="America/New_York"
DIFFICULTY="2"
VIEW_DISTANCE="32"
REGION_FILE_COMPRESSION="lz4"
ENABLE_WHITELIST="true"
HIDE_ONLINE_PLAYERS="true"
ENFORCE_SECURE_PROFILE="false"
PLAYER_IDLE_TIMEOUT="60"
ALLOW_FLIGHT="true"
ENABLE_ROLLING_LOGS="true"
```

Then allow the port: `sudo firewall-cmd --permanent --add-port=25565/tcp`
And then a similar config for the RLCraft server.
Effectively I took only the environment variable output of the configurator and put it in my env file.
This is not a perfect template, as a perfect template would use a `compose.override.yaml` file, but this is good enough for me.

So the new setup is services (including my Minecraft server) running on Podman managed by Arcane in a VM on XCP-NG.

# Future
Next I'd like to experiment with Arcane's ability to monitor and pull Git repos and automatically build and deploy from them.
Currently I have my website running in its own VM via Podman, and in the future I'd like for it to be on Kubernetes, but its the perfect project to experiment with this feature.
