FROM alpine:3.9

COPY upserter /

ENV BIND_HOST=upserter
ENV BIND_PORT=8080
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USERNAME=upserter
ENV DB_PASSWORD=fillmein
ENV DB_NAME=upserter

CMD [ \
  "/upserter" \
]
