FROM alpine:3.9

COPY sender wait-for.sh /

ENV IN_FILE_PATH=localhost
ENV BATCH_SIZE=32
ENV UPSERTER_HOST=upserter
ENV UPSERTER_PORT=8080


CMD [ \
  "/bin/sh", "-c", \
  "/wait-for.sh ${UPSERTER_HOST}:${UPSERTER_PORT} -- /sender -in-file-path ${IN_FILE_PATH} -batch-size ${BATCH_SIZE} -upserter-host ${UPSERTER_HOST} -upserter-port ${UPSERTER_PORT}" \
]