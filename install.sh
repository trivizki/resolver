# Install golang 
apt-get update
apt-get install golang-go

# Getting dependencies
apt-get install libpcap-dev

go get -v github.com/google/gopacket 
go get -v github.com/spf13/viper
go get -v github.com/go-redis/redis
go get -v github.com/mattn/go-sqlite3
go get -v github.com/go-sql-driver/mysql 

# Install redis
sudo apt install redis-server
sudo systemctl restart redis.service

# Install MySql
sudo apt install mysql-server

# Set Up SQL DB
sudo mysql <<EOF
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
FLUSH PRIVILEGES;
EOF
sudo mysql -u root -proot -e "CREATE DATABASE if not exists resolver;"




