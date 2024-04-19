#!/bin/bash

# Function to enable bandwidth limitation for specified port
enable_bandwidth() {
    tc qdisc add root dev lo handle 1: htb default 12
}

enable_bandwidth_limitation() {
    local port=$1
    local rate=$2  # rate in Mbps
    local id=$3

    tc class add dev lo parent 1: classid 1:${id} htb rate ${rate}mbit
    tc filter add dev lo protocol ip parent 1:0 prio 1 u32 match ip dport ${port} 0xffff flowid 1:${id}

    echo "Bandwidth limitation enabled for port ${port} with rate ${rate}Mbps"
}

# Function to disable bandwidth limitation for specified port
disable_bandwidth_limitation() {
    tc -s qdisc del root dev lo
}

# Main script
if [[ "$1" == "enable" ]]; then
    # Enable bandwidth limitation for specified ports
    enable_bandwidth
    enable_bandwidth_limitation 10000 10 1 
    enable_bandwidth_limitation 10011 10 2
    enable_bandwidth_limitation 10022 10 3
    enable_bandwidth_limitation 10033 10 4
    enable_bandwidth_limitation 10044 10 4
    enable_bandwidth_limitation 10055 10 4
    enable_bandwidth_limitation 10066 10 4
    enable_bandwidth_limitation 10077 10 4
    enable_bandwidth_limitation 10088 10 4
elif [[ "$1" == "disable" ]]; then
    # Disable bandwidth limitation for specified ports
    disable_bandwidth_limitation
elif [[ "$1" == "check" ]]; then
    tc class show dev lo
    tc filter show dev lo
else
    echo "Usage: $0 <enable|disable>"
    exit 1
fi
