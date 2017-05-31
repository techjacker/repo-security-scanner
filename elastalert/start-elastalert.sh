#!/bin/sh

set -e

# Set the timezone.
# if [ "$SET_CONTAINER_TIMEZONE" = "true" ]; then
# 	setup-timezone -z ${CONTAINER_TIMEZONE} && \
# 	echo "Container timezone set to: $CONTAINER_TIMEZONE"
# else
# 	echo "Container timezone not modified"
# fi

# Force immediate synchronisation of the time and start the time-synchronization service.
# In order to be able to use ntpd in the container, it must be run with the SYS_TIME capability.
# In addition you may want to add the SYS_NICE capability, in order for ntpd to be able to modify its priority.
ntpd -s

# Wait until Elasticsearch is online since otherwise Elastalert will fail.
# rm -f garbage_file
# while ! wget -O garbage_file ${ELASTICSEARCH_HOST}:${ELASTICSEARCH_PORT} 2>/dev/null
# do
# 	echo "Waiting for Elasticsearch..."
# 	rm -f garbage_file
# 	sleep 1
# done
# rm -f garbage_file
# sleep 5

# # Check if the Elastalert index exists in Elasticsearch and create it if it does not.
# if ! wget -O garbage_file ${ELASTICSEARCH_HOST}:${ELASTICSEARCH_PORT}/elastalert_status 2>/dev/null
# then
# 	echo "Creating Elastalert index in Elasticsearch..."
#     elastalert-create-index --host ${ELASTICSEARCH_HOST} --port ${ELASTICSEARCH_PORT} --config ${ELASTALERT_CONFIG} --index elastalert_status --old-index ""
# else
#     echo "Elastalert index already exists in Elasticsearch."
# fi
# rm -f garbage_file

echo "Starting Elastalert..."

# Set the rule directory in the Elastalert config file to external rules directory.
sed -i -e"s|rules_folder: [[:print:]]*|rules_folder: ${RULES_DIRECTORY}|g" ${ELASTALERT_CONFIG}
# Set the Elasticsearch host that Elastalert is to query.
sed -i -e"s|es_host: [[:print:]]*|es_host: ${ELASTICSEARCH_HOST}|g" ${ELASTALERT_CONFIG}
# Set the port used by Elasticsearch at the above address.
sed -i -e"s|es_port: [0-9]*|es_port: ${ELASTICSEARCH_PORT}|g" ${ELASTALERT_CONFIG}

exec elastlaert --config ${ELASTALERT_CONFIG}
