FROM consul:1.15.2

COPY config-file.json /consul/config

ENTRYPOINT ["consul", "agent", "-config-file=/consul/config/config-file.json"]