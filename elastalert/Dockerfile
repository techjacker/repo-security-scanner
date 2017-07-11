FROM iron/python:2

RUN mkdir -p /var/empty
 # Set this environment variable to true to set timezone on container start.
ENV SET_CONTAINER_TIMEZONE true
# Default container timezone as found under the directory /usr/share/zoneinfo/.
ENV CONTAINER_TIMEZONE Europe/London
# URL from which to download Elastalert.

WORKDIR /opt
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates openssl-dev libffi-dev python-dev gcc musl-dev tzdata openntpd curl && \
    rm -rf /var/cache/apk/* && \
# Install pip - required for installation of Elastalert.
    wget https://bootstrap.pypa.io/get-pip.py && \
    python get-pip.py && \
    rm get-pip.py && \
# Download and unpack Elastalert.
    wget https://github.com/Yelp/elastalert/archive/master.zip && \
    unzip *.zip && \
    rm *.zip && \
    mv e* elastalert

# Install Elastalert.
ENV ELASTALERT_ROOT /opt/elastalert
WORKDIR ${ELASTALERT_ROOT}
RUN python setup.py install && \
    pip install -e . && \
    pip install notifications-python-client && \
    pip uninstall twilio --yes && \
    pip install twilio==6.0.0

WORKDIR /opt
ENV RULES_DIRECTORY /opt/rules
COPY rules ${RULES_DIRECTORY}
ENV ELASTALERT_CONFIG /opt/elastalert_config.yaml
COPY elastalert_config.yaml ${ELASTALERT_CONFIG}
COPY modules ${ELASTALERT_ROOT}/elastalert_modules

COPY elasticsearch_mappings.json /tmp
COPY start-elastalert.sh /opt/
RUN chmod +x /opt/start-elastalert.sh
CMD ["/opt/start-elastalert.sh"]
