apt-get update

# Install redis
apt install redis-server
systemctl restart redis.service

# Install MySql
apt install mysql-server

# Set Up SQL DB
mysql <<EOF
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
FLUSH PRIVILEGES;
EOF
mysql -u root -proot -e "CREATE DATABASE if not exists resolver;"
