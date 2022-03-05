
PASSWORD=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
ZYM_PATH=$HOME/.zymurgauge

# create directories 
mkdir -p ${ZYM_PATH}/nginx
mkdir -p ${ZYM_PATH}/influxdb/data
mkdir -p ${ZYM_PATH}/influxdb/init
mkdir -p ${ZYM_PATH}/telegraf

# copy conf files that are to be updated to their respective data directories
cp config/nginx.conf ${ZYM_PATH}/nginx
cp config/influxdb.conf ${ZYM_PATH}/influxdb
cp config/create-telegraf.iql ${ZYM_PATH}/influxdb/init
cp config/telegraf.conf ${ZYM_PATH}/telegraf

# update conf files with generated password
sed -i 's/^  # password = "metricsmetricsmetricsmetrics"$/  password = "'${PASSWORD}'"/' ${ZYM_PATH}/telegraf/telegraf.conf
sed -i 's/metricsmetricsmetricsmetrics/'$PASSWORD'/g' ${ZYM_PATH}/influxdb/init/create-telegraf.iql