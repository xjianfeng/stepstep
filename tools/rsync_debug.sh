#!/bin/bash


echo '============================================================== rsync'

#同步配置
rsync -e "ssh -i ~/.key/id_rsa -p 8222" -cvropg --copy-unsafe-links --exclude-from="tools/excludeR.list" ./bin/ game@xxxxx.com:/home/game/stepstep/
echo '============================================================== rsync end'
