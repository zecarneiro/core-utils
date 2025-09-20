#!/usr/bin/env bash
# Author: José M. C. Noronha

title-log -m 'CLEANUP SNAP'
evalc "sudo sh -c 'rm -rf /var/lib/snapd/cache/*'"
warn-log "Removes old revisions of snaps"
read -n 1 -s -r -p "Please, CLOSE ALL SNAPS BEFORE RUNNING THIS. PRESS ANY KEY TO CONTINUE"
LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}' | while read snapname revision; do
  echo
  evalc "sudo snap remove \"$snapname\" --revision=\"$revision\""
done
