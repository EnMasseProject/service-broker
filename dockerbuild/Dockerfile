FROM alpine
ADD broker /broker
ADD dev.config.yaml /dev.config.yaml
ENTRYPOINT ["/broker"]
CMD ["-c", "dev.config.yaml"]
