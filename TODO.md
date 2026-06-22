
# Appearance
## Add more accent colors
make headers and accent colors alternate, following that of obsidian or vscode

# Blog
## Add a time to read estimator
## Add a word count
exclude code blocks
## Tag filtering for blog in atom feed
Add a howto page of how to use the atom feed
## Search for blog content
## Add a way to leave comments on posts

# Tools
## Make the CIDR tool templatable so that more tools can be made
give the CIDR tool its own page and embed the page into the tool overview
## Make the CIDR tool show verbose information, collapsed by default
add long ipv4 and ipv6 representation
show binary representations of IPs and masks
## Add a power price conversion tool
Also make it so you can put an initial cost and price per kwh and graph it
good for pricing server cost effectivness


# Server
## Harden nginx further
use dhi.io/nginx
https://knowledge.digicert.com/tutorials/enabling-perfect-forward-secrecy
## Get Cloudflared Working
## Make proper Dev and Prod deployments
## Update website via github
~~it pulls from the prod branch automagically~~
use CI/CD
build website container in github if free
setup rollbacks on for webserver

# Other
## Go through useful resources channel in server and make posts about that stuff
## Add a page 'Hire Me'
Have resume and CV available on it
Sensitive info protected by:
- key via get request, part of QR code
- OIDC session
- request form:
	Enter email, name, buisness, purpose, comments
	I get an email notifying of a request. Links to authenticated page on website
	On this website i can view request details, including ip address, browser info, and time
	I can approve or deny. Approving sends them a link where they can access the information
        Perhaps also attach a copy
