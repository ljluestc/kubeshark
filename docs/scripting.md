# Kubeshark Scripting Guide

## PCAP Handling in Scripts

When working with scripts in Kubeshark, especially those that process PCAP files, it's important to understand how the PCAP TTL (Time To Live) affects your scripts.

### Understanding PCAP TTL

By default, PCAP files in Kubeshark are stored for 60 seconds before being automatically deleted to conserve storage space. This TTL value can be adjusted in your configuration:
