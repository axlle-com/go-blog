#!/bin/bash

# Ждем запуска Elasticsearch
until curl -u elastic:changeme -s http://elasticsearch:9200/_cluster/health | grep '"status":"green"' > /dev/null; do
  echo "Waiting for Elasticsearch to be ready..."
  sleep 5
done

# Создаем файл конфигурации Kibana
cat <<EOF > /usr/share/kibana/config/kibana.yml
elasticsearch.hosts: ["http://elasticsearch:9200"]
elasticsearch.username: "kibana"
elasticsearch.password: "name"
EOF

echo "Kibana configured to use the 'elastic' user."

# Запуск Kibana
/usr/share/kibana/bin/kibana
