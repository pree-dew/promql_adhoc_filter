Small binary to modify any promql by adding any adhoc filter to any promql

Steps to use it:

1. git clone git@github.com:pree-dew/promql_adhoc_filter.git
2. go build .
3. ./adhoc_filter 'sum(http_request_total{}) by (job)' '{"job": "prod"}'
