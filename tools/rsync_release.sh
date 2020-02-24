#!/bin/bash

echo '============================================================== rsync'

#同步配置
rsync -e "ssh -i ~/.key/id_rsa -p 8222" -cvropg --copy-unsafe-links --exclude-from="tools/excludeR.list" ./bin/  game@xxxx.com:/home/game/stepstep/stepstep/

echo '============================================================== rsync end'

 ssh -i ~/.key/id_rsa -p 8222 game@xxxx.com "sudo -S supervisorctl restart stepstep"
