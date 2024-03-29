#!/bin/bash
ZYM_PATH=$HOME/.zymurgauge

# check for previous setup
if [ -d $ZYM_PATH ] 
then
    echo $'\e[31m'"$ZYM_PATH"' directory already exists'
    exit
fi

unset zym_username
unset zym_password
unset brewfather_user_id
unset brewfather_key
unset brewfather_log_url
unset influxdb_username
unset influxdb_password

# set credentials
read -p $'\e[32m?\e[0m Enter Zymurgauge admin account username : ' zym_username

prompt=$'\e[32m?\e[0m Enter Zymurgauge admin account password : '
while IFS= read -p "$prompt" -r -s -n 1 char
do
    if [[ $char == $'\0' ]]
    then
        break
    fi
    prompt='*'
    zym_password+="$char"
done

if [ ${#zym_password} -lt 8 ]; 
then
    echo
    echo $'\e[31mZymurgauge password is must be at least 8 characters\e[0m'
    exit
fi

echo

read -p $'\e[32m?\e[0m Enter Brewfather API User ID : ' brewfather_user_id

read -p $'\e[32m?\e[0m Enter Brewfather API Key : ' brewfather_key

read -p $'\e[32m?\e[0m Enter Brewfather API Log URL : ' brewfather_log_url

read -p $'\e[32m?\e[0m Enter InfluxDB admin account username : ' influxdb_username

prompt=$'\e[32m?\e[0m Enter InfluxDB admin account password : '
while IFS= read -p "$prompt" -r -s -n 1 char
do
    if [[ $char == $'\0' ]]
    then
        break
    fi
    prompt='*'
    influxdb_password+="$char"
done

if [ ${#influxdb_password} -lt 8 ]; 
then
    echo
    echo $'\e[31mInfluxDB password is must be at least 8 characters\e[0m'
    exit
fi

echo

# create directories 
mkdir -p $ZYM_PATH/data
mkdir -p $ZYM_PATH/nginx
mkdir -p $ZYM_PATH/influxdb
mkdir -p $ZYM_PATH/telegraf

influxdb_token=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 64 ; echo '')

# download config files that are to be updated to their respective data directories
wget -P $ZYM_PATH/nginx/ https://raw.githubusercontent.com/benjaminbartels/zymurgauge/master/config/nginx.conf
wget -P $ZYM_PATH/telegraf https://raw.githubusercontent.com/benjaminbartels/zymurgauge/master/config/telegraf.conf

# create self signed cert
openssl req \
    -new \
    -newkey rsa:4096 \
    -days 365 \
    -nodes \
    -x509 \
    -subj "/C=US/CN=$(hostname).local" \
    -keyout $ZYM_PATH/nginx/cert.key \
    -out $ZYM_PATH/nginx/cert.pem

# set token in telegraf.conf
sed -i 's/^#   token = ""$/  token = "'${influxdb_token}'"/' $ZYM_PATH/telegraf/telegraf.conf

echo $'\e[32mSetting up InfluxDB\e[0m'

# initalize influxdb
docker run -d -p 8086:8086 \
    --name influxdb_setup \
    -v $ZYM_PATH/influxdb:/var/lib/influxdb2 \
    influxdb:2.2.0

# wait for influx docker container to be ready
sleep 10s

# create admin credentians and initial org and bucket
docker exec influxdb_setup influx setup -f \
    --username $influxdb_username \
    --password $influxdb_password \
    --token $influxdb_token \
    --org zymurgauge \
    --bucket telegraf

# get bucket id
bucket_id=$(docker exec influxdb_setup influx bucket ls --name telegraf --hide-headers | \
    while read -a array; do echo "${array[0]}" ; done)

# create read_token
read_token=$(docker exec influxdb_setup influx auth create --hide-headers --org zymurgauge --read-bucket $bucket_id | \
    while read -a array; do echo "${array[1]}" ; done)

# kill and delete influxdb_setup container
docker kill influxdb_setup
docker rm influxdb_setup

echo $'\e[32mInfluxDB setup complete\e[0m'

echo $'\e[32mSetting up Zymurgauge\e[0m'

influx_url=https://$(hostname).local:8086

# initalize zymurgauge
docker run --rm -v $ZYM_PATH/data:/data \
    ghcr.io/benjaminbartels/zymurgauge:latest init --username=$zym_username --password=$zym_password \
    --brewfather-user-id=$brewfather_user_id \
    --brewfather-key=$brewfather_key \
    --brewfather-log-url=$brewfather_log_url \
    --influx-dburl=$influx_url \
    --influx-db-token=$read_token \
    --stats-d-address=telegraf:8125

echo $'\e[32mZymurgauge setup complete\e[0m'