FROM elastic/filebeat:8.14.3

COPY filebeat.yml /usr/share/filebeat/filebeat.yml

WORKDIR /usr/share/filebeat

CMD ["filebeat", "-e"]
