# syntax=docker/dockerfile:1
FROM gcr.io/distroless/base-debian11
ADD ../cmd/accrual ./
EXPOSE 8080
ENTRYPOINT ["/accrual"]