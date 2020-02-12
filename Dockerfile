FROM scratch

COPY "bin/properties.bin" /
ENTRYPOINT ["/properties.bin"]
