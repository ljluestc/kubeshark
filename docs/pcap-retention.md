# PCAP Retention in Kubeshark

## Overview

Kubeshark captures network traffic as PCAP (Packet Capture) files for analysis. By default, these files are automatically cleaned up after a configurable TTL (Time To Live) period to manage disk usage. However, in some scenarios, especially when using scripts for Digital Forensics and Incident Response (DFIR), you may need to retain specific PCAP files for longer periods.

## Default PCAP TTL

The default TTL for PCAP files is 60 seconds. This means that if a PCAP file is not accessed or marked for extended retention, it will be deleted 60 seconds after creation.

## Configuring PCAP TTL

You can adjust the PCAP TTL in your Kubeshark configuration:

```yaml
tap:
  misc:
    pcapTTL: 300s  # Set to 5 minutes
