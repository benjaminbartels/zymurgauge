ZYM_PATH=$HOME/.zymurgauge

# create directories 
mkdir -p ${ZYM_PATH}/nginx
mkdir -p ${ZYM_PATH}/influxdb/data
mkdir -p ${ZYM_PATH}/influxdb/init
mkdir -p ${ZYM_PATH}/telegraf

# copy conf files that are to be updated to their respective data directories
cp config/nginx.conf ${ZYM_PATH}/nginx
cp config/influxdb.conf ${ZYM_PATH}/influxdb
cp config/init.iql ${ZYM_PATH}/influxdb/init
cp config/telegraf.conf ${ZYM_PATH}/telegraf

# initialze influxdb
docker run --rm \
    -v $ZYM_PATH/influxdb/influxdb.conf:/etc/influxdb/influxdb.conf:ro \
    -v $ZYM_PATH/influxdb/data:/var/lib/influxdb \
    -v $ZYM_PATH/influxdb/init:/docker-entrypoint-initdb.d \
    arm32v7/influxdb:1.8.10 /init-influxdb.sh

# enable flex
sed -i 's/^  flux-enabled = false$/    flux-enabled = true/' ${ZYM_PATH}/influxdb/influxdb.conf