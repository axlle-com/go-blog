FROM docker.elastic.co/kibana/kibana:8.14.3

USER root
COPY init.sh /usr/local/bin/kibana-init.sh
RUN chmod +x /usr/local/bin/kibana-init.sh
USER kibana
