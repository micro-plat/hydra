cd /home/yanglei/work/src/github.com/micro-plat/hydra/docs/bin/dist
rm -rf ./*


cd /home/yanglei/work/src/github.com/micro-plat/hydra/docs/bin/webserver/src
zip -r webserver.zip ./

mv ./webserver.zip /home/yanglei/work/src/github.com/micro-plat/hydra/docs/bin/dist
cd /home/yanglei/work/src/github.com/micro-plat/hydra/docs/bin/webserver

go build
cp ./webserver /home/yanglei/work/src/github.com/micro-plat/hydra/docs/bin/dist

