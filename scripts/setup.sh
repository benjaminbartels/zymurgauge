ZYM_PATH=$HOME/.zymurgauge

# TODO: run zym init with user and pass
# TODO: dont run containers as root

# create directories 
mkdir -p ${ZYM_PATH}/data
mkdir -p ${ZYM_PATH}/nginx
mkdir -p ${ZYM_PATH}/influxdb/data
mkdir -p ${ZYM_PATH}/influxdb/config
mkdir -p ${ZYM_PATH}/influxdb/init
mkdir -p ${ZYM_PATH}/telegraf

# copy conf files that are to be updated to their respective data directories
cp config/nginx.conf ${ZYM_PATH}/nginx
cp config/setup-v1.sh ${ZYM_PATH}/influxdb/init
cp config/telegraf.conf ${ZYM_PATH}/telegraf
