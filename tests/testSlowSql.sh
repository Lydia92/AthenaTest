#!/bin/bash

#config monitor database server
athena_db_host="192.168.21.210"
athena_db_port=3326
athena_db_user="athena"
athena_db_password="123456"
athena_db_database="athena"

#config mysql server
mysql_client="/usr/local/mysql/bin/mysql"
mysql_host="192.168.21.208"
mysql_port=3301
mysql_user="dba"
mysql_password="123456"

#config slowqury
slowquery_dir="/data/mysql/mysqldata3301/log/"
slowquery_long_time=1
slowquery_file=`$mysql_client -h$mysql_host -P$mysql_port -u$mysql_user -p$mysql_password  -e "show variables like 'slow_query_log_file'"|grep log|awk '{print $2}'`
pt_query_digest="/usr/bin/pt-query-digest"

#config server_id
athena_server_id="192.168.21.208-3301"

#collect mysql slowquery log into athena database
$pt_query_digest --user=$athena_db_user --password=$athena_db_password --port=$athena_db_port --review h=$athena_db_host,D=$athena_db_database,t=mysql_slow_query_review  --history h=$athena_db_host,D=$athena_db_database,t=mysql_slow_query_review_history  --no-report --limit=100%  --filter=" \$event->{add_column} = length(\$event->{arg}) and \$event->{serverid}='$athena_server_id'  " $slowquery_file > /tmp/athena_slowquery.log

#config mysql slowquery
tmp_log=`$mysql_client -h$mysql_host -P$mysql_port -u$mysql_user -p$mysql_password -e "select concat('$slowquery_dir','slowquery_',date_format(now(),'%Y_%m_%d_%H_%i'),'.log');"|grep log|sed -n -e '2p'`


$mysql_client -h$mysql_host -P$mysql_port -u$mysql_user -p$mysql_password -e "set global slow_query_log=1;set global long_query_time=$slowquery_long_time;"
$mysql_client -h$mysql_host -P$mysql_port -u$mysql_user -p$mysql_password -e "set global slow_query_log_file = '$tmp_log'; "

#delete log before 7 days
/usr/bin/find $slowquery_dir -name 'slowquery_*' -mtime +7|xargs rm -rf ;

####END####