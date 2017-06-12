FROM iron/python:2

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
    wget https://github.com/Yelp/elastalert/archive/master.zip && \
    unzip *.zip && \
    rm *.zip && \
    mv e* elastalert

# Install Elastalert.
WORKDIR /opt/elastalert
RUN python setup.py install && \
    pip install -e . && \
    pip uninstall twilio --yes && \
    pip install twilio==6.0.0

ENV RULES_DIRECTORY /opt/rules
COPY rules ${RULES_DIRECTORY}
ENV ELASTALERT_CONFIG /opt/elastalert_config.yaml
COPY elastalert_config.yaml ${ELASTALERT_CONFIG}

COPY start-elastalert.sh /opt/
RUN chmod +x /opt/start-elastalert.sh
CMD ["/opt/start-elastalert.sh"]
