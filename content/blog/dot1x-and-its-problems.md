---
date: 2026-08-29
title: Dot1x deployment and its problems
summary: A summary of some of the problems and cool stuff I encountered when deploying 802.1X
tags:
  - rant
  - 8021x
  - networking
author: Ethan Ashley
created: 2026-08-18
updated: 2026-08-29
draft: false
publish: true
---
The following is less of a formalized post of facts and more of a frustrated listing of thoughts fueled by my experiences.
That does not go to say that it lacks value intrinsically, as i hope that someone exploring deploying 802.1x in its entirety

# Problems
## Client differences
certificate trust is required but it does not reference the device's trust store
	either needs a pinned certificate via management (fine) or user to trust (they don't know whether it is trusted or not)
	it'll do this even when the root certs are installed on the device
android has so many options, even those not offered
	confusing to end user
apple has almost no options
	for BYOD connections EAP-TTLS is not supported
some windows machines have problems with TLS 1.3 with FreeRadius and session resumption still doesn't work right

## External directories
you need a directory which can check both cleartext passwords and NT Hashes (because MSCHAPv2)
LDAP is the only viable option
oauth only works with PAP/EAP-TTLS PAP
most IdPs dont store plaintext passwords (good) so you need some way to get the hash or pass the request to the IdP
	Entra is incredibly inflexible in this reguard
	[GitHub - jimdigriz/freeradius-oauth2-perl](https://github.com/jimdigriz/freeradius-oauth2-perl) exists
		relies on soon to be removed functionality from entra
		solo maintainer (i salute you sir)
		perl script
		easily hits entra rate limits in production

## Entra
we have a connector setup to sync AD to Entra
	to reverse this would be too risky at this point
	cloud created account do not propogate to AD as it is a one way sync
you can setup an AD LDS server, install the connector, and sync entra to it
	this provides LDAP access to entra users and their data
	all problems are solved
		passwords are set but blank lol get fucked hahaha you just spent 5 hours on a new windows server for nothing hahahaha

## Cisco WLC and APs
Documentation says only dTLS (RADSEC UDP) is supported
	configuring RADSEC TCP like a Cisco switch works perfectly
APs only support EAP-TLS, EAP-PEAP MSCHAPv2, or EAP-FAST for authentication
	EAP-TLS requires SCEP certificate installation
	No TTLS or PAP

## Cisco NETCONF
not all IBNS configurations available in the CLI are available via NETCONF
	schema definitions, but my switches (latest stable firmware) do not support these namespaces
	this results in having to automate via CLI (fine) and building out a robust rollback flow in case of lockout or loss of connectivity
Cisco has a [cool tool for browsing schema](https://github.com/CiscoDevNet/yangsuite) which helps

## Mosyle MDM
IOS and tvOS dot1x are outright broken for EAP-TLS (as of now)
Creating the profile using [iMazings](https://imazing.com/profile-editor) tool and installing as a custom profile works
requires SCEP for user certificates unlike Intune

## WebAuth
could be a great solution
no way to automate the client side or provide consistent behavior
intrusive to user flow on managed devices
authorizes based on MAC address, susceptible to spoofing
	also some devices (mobile phones) randomize their MAC addresses and can cause frequent reauthentications

## CoA
requires a separate port and static clients, sending CoA requests over the existing RADSEC tunnel is not well supported
its cool and works but is still UDP, there is little assurance that a CoA went through

## Open ports pre-authentication
a common migration strategy is to have the ports fail open
	the port cannot be open to start with, as many clients will get an IP address, then when they get placed into their proper VLAN, they hold on to the old address and cannot communicate, as the gateway has changed
		CoA can fix this, but is not reliable enough and gets weird with PoE devices

# The cool stuff
using IBNS 2.0 you can apply service profiles, interface templates, and even macros, with freeradius as the central authority
	Can be made better with switches reaching to an SFTP repository for configurations, at the cost of latency and single point of failure
EAP-TLS works wonderfully
	thats the best part of all of this
## Certificates
You now need to maintain a CA and it'll probably be easier to have a SCEP server
You can also manage policy decisions in FreeRadius using Policy OIDs which makes things a lot easier (1.3.6.1.4.1.###)
	The structure and documentation of this must be maintained externally from the CA

# For the future
A more complete solution is more than likely the way to go, something like PacketFence or Cisco ISE
