#!/bin/sh

# Workaround for default kernel settings: https://github.com/tatsushid/go-fastping/issues/25#issuecomment-236203705
sysctl -w net.ipv4.ping_group_range="0 65535"

exec /pingrok $@
