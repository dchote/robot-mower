#!/bin/sh
rsync --delete --ignore-errors -rlu  . robot-mower.local:/home/mower/gocode/src/github.com/dchote/robot-mower/src/
ssh -t robot-mower.local /home/mower/gocode/src/github.com/dchote/robot-mower/src/build_run.sh
