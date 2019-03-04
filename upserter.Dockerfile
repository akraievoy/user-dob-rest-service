FROM alpine:3.9

COPY upserter/upserter /

CMD ["/upserter"]
