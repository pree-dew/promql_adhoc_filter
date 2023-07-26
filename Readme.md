Small binary to modify any promql by adding any adhoc filter to any promql

Steps to use it:

1. git clone git@github.com:pree-dew/promql_adhoc_filter.git
2. go build .
3. ./adhoc_filter 'sum(rate(http_request_total{}[1m])) by (job, container) - sum(http{}) by (container)' '{"job": "test"}'

   parsed expr: sum(http_request_total{job="test"}[1m]) by(job,container) - sum(http) by(container)
