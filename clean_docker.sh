sudo docker rm -f $(sudo docker ps -aq)
sudo docker network prune
sudo docker volume prune
cd fixtures && docker-compose up -d
# setup mysql
cd ..
# mysql data
cd app_server/mysql/
docker run -id -p 3306:3306 --name=mysql -v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root mysql
#docker exec -it mysql /bin/bash
#mysql -uroot -proot
#create database cryptology
cd ..
#rm education
#go build
#./education