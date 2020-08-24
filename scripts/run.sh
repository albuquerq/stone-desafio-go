#/bin/sh


export DB_HOST=""
export DB_PORT=""
export DB_USER=""
export DB_NAME=""
export DB_PASS=""
export PORT=""

echo "Executing in $(pwd)"

go run ../cmd/bankingd/bankingd.go