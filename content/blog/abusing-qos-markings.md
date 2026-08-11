---
date: 2026-08-10
title: Abusing QoS Markings
summary: Some notes from when I had the splended idea to see if my university had a trust boundary configured.
tags:
  - networking
  - projects
  - curious
  - linux
  - qos
author: Ethan Ashley
created: 2025-04-07
updated: 2025-04-12
draft: false
publish: true
---

# Configuration
```
# Set DSCP
#sudo iptables -t mangle -A OUTPUT -j TOS --set-tos 184  # This line just changes the DSCP. See 'dscp' in `man iptables-extensions`
sudo iptables -t mangle -A OUTPUT -j DSCP --set-dscp 46
sudo iptables -t mangle -L -v --line-numbers 
sudo ip6tables -t mangle -A OUTPUT -j DSCP --set-dscp 46
sudo ip6tables -t mangle -L -v --line-numbers
sudo netfilter-persistent save
```

to remove entries: `sudo iptables -t mangle -D OUTPUT #`
OR
using firewalld
```
sudo firewall-cmd --permanent --direct --add-rule ipv4 mangle OUTPUT 0 -j DSCP --set-dscp 46
sudo firewall-cmd --permanent --direct --add-rule ipv6 mangle OUTPUT 0 -j DSCP --set-dscp 46
sudo firewall-cmd --reload
sudo firewall-cmd --direct --get-all-rules
```

mtr command to test: `mtr -rnc 30 1.1.1.1`

What about PCP?
	this requires the traffic to be tagged, and this doesnt work well on access ports, especially across networks which you do not know the configuration of

# Environment
Configuration and testing was done on an Ubuntu Linux 24.04 laptop
	Dell XPS 9520
	Tests were done back to back in the same location over WiFi
	this system has a lot of other things done to it from over the years, but they shouldnt affect results
Tests were ran from Michigan Technological University's campus
	Tests ran at McNair ResNet were done at about 10:30 PM on a sunday night
	Tests ran at Fisher Hall were done at about 10:00 AM on a monday morning during a large lecture (college physics 2?)

# Results
### Ping from McNair ResNet
ping before changes (after they were undone cause testing)
```
--- 1.1.1.1 ping statistics ---
50 packets transmitted, 50 received, 0% packet loss, time 49066ms
rtt min/avg/max/mdev = 18.864/29.052/139.863/27.824 ms
```
ping DIRECTLY after changes
```
--- 1.1.1.1 ping statistics ---
50 packets transmitted, 50 received, 0% packet loss, time 49068ms
rtt min/avg/max/mdev = 18.873/19.774/25.593/1.415 ms
```

### mtr from McNair ResNet 
without changes
```
╰─$ mtr -rnc 30 1.1.1.1
Start: 2025-04-06T23:10:00-0400
HOST: DIANE                       Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- 141.219.220.3              0.0%    30    1.6   2.2   1.5   5.2   0.9
  2.|-- 172.31.255.182             0.0%    30    2.1   2.9   1.9  10.7   1.6
  3.|-- 141.219.183.113            0.0%    30    3.5   5.1   2.0  19.4   4.5
  4.|-- 207.75.40.9                0.0%    30    2.3   3.5   2.1   9.5   1.8
  5.|-- 207.72.231.19              0.0%    30    7.5  18.2   3.6 142.8  29.9
  6.|-- 207.72.231.5               0.0%    30    6.9   7.1   6.0  14.4   1.5
  7.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  8.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  9.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
 10.|-- 207.72.231.31              0.0%    30   13.0  14.0  12.1  31.4   3.5
 11.|-- 208.115.136.180           60.0%    30   26.2  26.8  19.7  59.1  10.9
 12.|-- 141.101.73.212             0.0%    30   22.0  28.2  18.8  70.1  14.0
 13.|-- 1.1.1.1                    0.0%    30   19.2  22.8  18.8  73.1  12.7
```

with changes (this time ran first)
```
╰─$ mtr -rnc 30 1.1.1.1
Start: 2025-04-06T23:08:30-0400
HOST: DIANE                       Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- 141.219.220.3              0.0%    30    2.3   2.4   1.5  10.5   1.8
  2.|-- 172.31.255.182             0.0%    30    3.5   2.7   1.9   9.5   1.6
  3.|-- 141.219.183.113            0.0%    30    2.5   3.4   2.0  11.4   1.9
  4.|-- 207.75.40.9                0.0%    30    2.2   2.7   2.0   5.5   0.8
  5.|-- 207.72.231.19              0.0%    30   10.8  13.4   3.2  65.2  15.6
  6.|-- 207.72.231.5               0.0%    30    6.3   6.8   6.2   8.7   0.7
  7.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  8.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  9.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
 10.|-- 207.72.231.31              0.0%    30   12.6  14.4  12.0  39.9   5.1
 11.|-- 208.115.136.180           36.7%    30   20.7  30.3  19.0 122.4  23.8
 12.|-- 141.101.73.212             0.0%    30   19.4  22.2  19.1  40.5   4.8
 13.|-- 1.1.1.1                    0.0%    30   19.0  19.3  19.0  21.0   0.4
```
### mtr from Fisher
without changes
```
╰─$ mtr -rnc 30 1.1.1.1                                                                                          
Start: 2025-04-07T10:00:24-0400
HOST: DIANE                       Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- 141.219.220.3              0.0%    30    7.9   6.9   1.7  37.8   6.9
  2.|-- 172.31.255.182             0.0%    30    2.7   6.2   2.5  20.4   5.1
  3.|-- 141.219.183.113            0.0%    30    4.6   8.9   2.7  44.7   8.9
  4.|-- 207.75.40.9                0.0%    30    3.8   7.2   2.8  20.3   5.2
  5.|-- 207.72.231.19              0.0%    30    7.1  16.0   3.8  91.0  18.5
  6.|-- 207.72.231.5               0.0%    30    7.5  10.1   6.8  27.6   4.8
  7.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  8.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  9.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
 10.|-- 207.72.231.31              0.0%    30   26.4  17.8  12.6  41.2   6.6
 11.|-- 208.115.136.180           66.7%    30   21.4  31.4  21.0  52.4  10.0
 12.|-- 141.101.73.212             0.0%    30   38.3  26.0  19.5  44.8   6.5
 13.|-- 1.1.1.1                    0.0%    30   22.2  22.2  19.3  37.8   3.7
```

with changes
```
╰─$ mtr -rnc 30 1.1.1.1
Start: 2025-04-07T10:02:10-0400
HOST: DIANE                       Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- 141.219.220.3              0.0%    30    3.1   3.3   2.1   7.9   1.4
  2.|-- 172.31.255.182             0.0%    30    3.1   3.9   2.5  12.3   1.9
  3.|-- 141.219.183.113            0.0%    30    3.3   5.1   2.0  35.1   5.9
  4.|-- 207.75.40.9                0.0%    30    3.5   3.6   2.6   5.4   0.8
  5.|-- 207.72.231.19              0.0%    30    4.0  17.6   3.2  90.4  20.6
  6.|-- 207.72.231.5               0.0%    30    7.5   8.5   6.6  28.5   3.9
  7.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  8.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
  9.|-- ???                       100.0    30    0.0   0.0   0.0   0.0   0.0
 10.|-- 207.72.231.31              0.0%    30   16.2  13.9  12.3  21.1   1.7
 11.|-- 208.115.136.180           73.3%    30   48.1  41.8  19.9 141.3  41.2
 12.|-- 141.101.73.212             0.0%    30   21.2  22.0  19.6  32.5   3.1
 13.|-- 1.1.1.1                    0.0%    30   20.3  21.5  19.2  36.1   4.0
```

### fast.com speedtest from Fisher Hall
without changes
```
120
Mbps
Latency
Unloaded
12 ms
Loaded
61 ms
Upload
Speed
120 Mbps
Client   Marquette, US   141.219.223.132   
Server(s) Ann Arbor, US  |  Chicago, US
```

with changes
```
120
Mbps
Latency
Unloaded
12 ms
Loaded
123 ms
Upload
Speed
240 Mbps
Client   Marquette, US   141.219.223.132   
Server(s) Ann Arbor, US  |  Chicago, US
```


# References
[Type of Service (ToS) and DSCP Values - LinuxReviews](https://linuxreviews.org/Type_of_Service_(ToS)_and_DSCP_Values)

