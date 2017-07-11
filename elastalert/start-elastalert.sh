#!/bin/sh
set -e

####################
# create elasticsearch mappings for index monitoring
# elastlaert will throw an error if these do not already exist
####################
curl -s "$ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT/githubintegration/_search"
if [[ $? != 0 ]]; then
    curl \
        --upload-file /tmp/elasticsearch_mappings.json \
        "$ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT/githubintegration" | python -m json.tool
fi

####################
# Set the timezone
####################
if [ "$SET_CONTAINER_TIMEZONE" = "true" ]; then
    cp /usr/share/zoneinfo/${CONTAINER_TIMEZONE} /etc/localtime && \
    echo "${CONTAINER_TIMEZONE}" >  /etc/timezone && \
    echo "Container timezone set to: $CONTAINER_TIMEZONE"
else
    echo "Container timezone not modified"
fi
# Force immediate synchronisation of the time and start the time-synchronization service.
# In order to be able to use ntpd in the container, it must be run with the SYS_TIME capability.
# In addition you may want to add the SYS_NICE capability, in order for ntpd to be able to modify its priority.
ntpd -s

####################
# Main config
####################
echo ""
echo "rules_folder: $RULES_DIRECTORY" >> "$ELASTALERT_CONFIG"
echo ""
echo "es_host: $ELASTICSEARCH_HOST" >> "$ELASTALERT_CONFIG"
echo ""
echo "es_port: $ELASTICSEARCH_PORT" >> "$ELASTALERT_CONFIG"


####################
# Rules config
####################
# NOTIFICATION_EMAILS = comma separated list of email addresses
echo ""
echo "email:" >> "$RULES_DIRECTORY/new_violation.yaml"
for i in $(echo $NOTIFICATION_EMAILS | tr "," "\n"); do
  echo "- $i" >> "$RULES_DIRECTORY/new_violation.yaml"
done


# Wait until Elasticsearch is online since otherwise Elastalert will fail.
rm -f garbage_file
while ! wget -O garbage_file "$ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT" 2>/dev/null
do
	echo "Waiting for Elasticsearch..."
	rm -f garbage_file
	sleep 1
done
rm -f garbage_file
sleep 5

# Check if the Elastalert index exists in Elasticsearch and create it if it does not.
if ! wget -O garbage_file "$ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT/elastalert_status" 2>/dev/null
then
	echo "Creating Elastalert index in Elasticsearch..."
    elastalert-create-index \
    	--host "$ELASTICSEARCH_HOST" \
    	--port "$ELASTICSEARCH_PORT" \
    	--config "$ELASTALERT_CONFIG" \
    	--index elastalert_status \
    	--old-index ""
else
    echo "Elastalert index already exists in Elasticsearch."
fi
rm -f garbage_file

echo "Starting Elastalert..."
exec elastalert --config "$ELASTALERT_CONFIG" --verbose
