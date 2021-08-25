#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)

ansible-playbook playbook.yaml --syntax-check

####
# /pkg/ueV1/{playbook.yaml, vars.yaml}
ansible-playbook playbook.yaml --tags run

ansible-playbook playbook.yaml --tags "get_log,get_status"
sleep 20

####
ansible-playbook playbook.yaml --tags execute --extra-vars "call=suspend"
ansible-playbook playbook.yaml --tags get_status
sleep 5

ansible-playbook playbook.yaml --tags execute --extra-vars "call=resume"
ansible-playbook playbook.yaml --tags get_status
sleep 20

####
ansible-playbook playbook.yaml --tags execute --extra-vars "call=kill"
ansible-playbook playbook.yaml --tags "get_log,get_status"
