cd ~/work/src/github.com/micro-plat/hydra/docs/bin/dist
rm -rf ./*


cd ~/work/src/github.com/micro-plat/hydra/docs/bin/webserver/src
zip -r webserver.zip ./

mv ./webserver.zip ~/work/src/github.com/micro-plat/hydra/docs/bin/dist
cd ~/work/src/github.com/micro-plat/hydra/docs/bin/webserver

go build
cp ./webserver ~/work/src/github.com/micro-plat/hydra/docs/bin/dist

