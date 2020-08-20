cd ~/work/src/github.com/hydra-pkg/docs/webserver
sudo npm run build
cd ~/work/src/github.com/hydra-pkg/docs/webserver/dist
sudo rm -rf menus.json
sudo rm -rf favicon.ico
sudo rm -rf 03-settings.html
sudo rm -rf ./img

cd ~/work/src/github.com/micro-plat/hydra/docs/bin/webserver/src
sudo rm -rf ./css
sudo rm -rf ./fonts
sudo rm -rf ./js
sudo rm -rf ./index.html



cp -r ~/work/src/github.com/hydra-pkg/docs/webserver/dist/* ./
