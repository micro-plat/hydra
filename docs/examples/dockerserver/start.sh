sudo docker build -t dockerserver .
sudo docker run --name demo -p 8090:8090 -d dockerserver
