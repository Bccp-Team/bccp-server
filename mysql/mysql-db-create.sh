#!/bin/bash

# Functions
ok() { echo -e '\e[32m'$1'\e[m'; } # Green

EXPECTED_ARGS=2
E_BADARGS=65
MYSQL=`which mysql`

if [ $# -ne $EXPECTED_ARGS ]
then
	echo "Usage: $0 dbname dbuser"
	exit $E_BADARGS
fi

echo -n "Database password: "
read -s password
echo
echo -n "Confirm password: "
read -s confirmed_password
echo

if [ $password != $confirmed_password ]
then
	echo "Passwords mismatch"
	exit $E_BADARGS
fi

ok "Creating database"
Q1="CREATE DATABASE IF NOT EXISTS $1;"
Q2="GRANT ALL ON *.* TO '$2'@'localhost' IDENTIFIED BY '$password';"
Q3="FLUSH PRIVILEGES;"
SQL="${Q1}${Q2}${Q3}"

$MYSQL -u root -p -e "$SQL"

ok "Database $1 and user $2 created"
