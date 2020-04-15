cd ../migrations
echo "Up test migrations"
goose postgres "user=$1 password=$2 dbname=$3 sslmode=disable" up
echo "Down test migrations"
goose postgres "user=$1 password=$2 dbname=$3 sslmode=disable" down
cd ~