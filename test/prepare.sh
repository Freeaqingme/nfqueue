#!/usr/bin/env bash

iptables -F
iptables -t nat -F
iptables -t mangle -F
iptables -t filter -F

iptables -A INPUT -j NFQUEUE --queue-balance 0:3
iptables -A OUTPUT -j NFQUEUE --queue-balance 0:3
iptables -A FORWARD -j NFQUEUE --queue-balance 0:3