#!/bin/sh
set -e

####################
# Main config
####################
# Set the rule directory in the Elastalert config file to external rules directory.
sed -i -e"s|rules_folder: [[:print:]]*|rules_folder: ${RULES_DIRECTORY}|g" "$ELASTALERT_CONFIG"
# Set the Elasticsearch host that Elastalert is to query.
sed -i -e"s|es_host: [[:print:]]*|es_host: ${ELASTICSEARCH_HOST}|g" "$ELASTALERT_CONFIG"
# Set the port used by Elasticsearch at the above address.
sed -i -e"s|es_port: [0-9]*|es_port: ${ELASTICSEARCH_PORT}|g" "$ELASTALERT_CONFIG"

####################
# Rules config
####################
# NOTIFICATION_EMAILS = comma separated list of email addresses
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
exec elastalert --config "$ELASTALERT_CONFIG"
