FROM iron/python:2


# ENV ELASTICSEARCH_HOST elasticsearchhost
# ENV ELASTICSEARCH_PORT 9200


# Set this environment variable to true to set timezone on container start.
ENV SET_CONTAINER_TIMEZONE false
# Default container timezone as found under the directory /usr/share/zoneinfo/.
ENV CONTAINER_TIMEZONE Europe/London
# URL from which to download Elastalert.
ENV ELASTALERT_URL https://github.com/Yelp/elastalert/archive/master.zip
# Directory to which Elastalert and Supervisor logs are written.
ENV LOG_DIR /opt/logs
# Elastalert home directory name.
ENV ELASTALERT_DIRECTORY_NAME elastalert
# Elastalert home directory full path.
ENV ELASTALERT_HOME /opt/${ELASTALERT_DIRECTORY_NAME}

# Install software required for Elastalert and NTP for time synchronization.
WORKDIR /opt
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates openssl-dev libffi-dev python-dev gcc musl-dev tzdata openntpd && \
    rm -rf /var/cache/apk/* && \
# Install pip - required for installation of Elastalert.
    wget https://bootstrap.pypa.io/get-pip.py && \
    python get-pip.py && \
    rm get-pip.py && \
# Download and unpack Elastalert.
    wget ${ELASTALERT_URL} && \
    unzip *.zip && \
    rm *.zip && \
    mv e* ${ELASTALERT_DIRECTORY_NAME}

# Install Elastalert.
WORKDIR ${ELASTALERT_HOME}
RUN python setup.py install && \
    pip install -e . && \
    pip uninstall twilio --yes && \
    pip install twilio==6.0.0

ENV RULES_DIRECTORY /opt/rules
COPY rules ${RULES_DIRECTORY}
ENV ELASTALERT_CONFIG /opt/elastalert_config.yaml
COPY elastalert_config.yaml ${ELASTALERT_CONFIG}

RUN mkdir -p ${LOG_DIR} && \
    mkdir -p /var/empty




    # apk del python-dev && \
    # apk del musl-dev && \
    # apk del gcc && \
    # apk del openssl-dev && \
    # apk del libffi-dev && \

COPY start-elastalert.sh /opt/
RUN chmod +x /opt/start-elastalert.sh
CMD ["/opt/start-elastalert.sh"]
