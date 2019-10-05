FROM prom/prometheus:v2.10.0

USER root

RUN mkdir -p /tsdb/data && \
    chown -R nobody:nogroup /tsdb/data

USER       nobody
EXPOSE     9090
VOLUME     ["/prometheus"]
WORKDIR    /prometheus
ENTRYPOINT ["/bin/prometheus"]
CMD        ["--config.file=/etc/prometheus/prometheus.yml", \
            "--storage.tsdb.path=/prometheus", \
            "--web.console.libraries=/usr/share/prometheus/console_libraries", \
            "--web.console.templates=/usr/share/prometheus/consoles"]
