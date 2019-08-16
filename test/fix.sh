#!/usr/bin/env bash

iptables -F INPUT
iptables -F OUTPUT
iptables -F FORWARD
iptables -F
iptables -t nat -F
iptables -t mangle -F
iptables -t filter -F