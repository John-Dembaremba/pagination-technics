#!/bin/bash

# OWASP ZAP Security Scans
#
# ======== Baseline tests ================
# docker compose run --rm zap /zap/zap-baseline.py \
#   -t "http://app:3030/users/cursor-based?cursor=10&limit=20" \
#   -r "baseline-cursor-report.html" \
#   -z "-config api.key=${ZAP_API_KEY}"

# docker compose run --rm zap /zap/zap-baseline.py \
#   -t "http://app:3030/users/limit-offset?page=1&limit=20" \
#   -r "baseline-limit-offset-report.html" \
#   -z "-config api.key=${ZAP_API_KEY}"

# ============ Full Scans ===================

# OWASP ZAP Active Scans
docker compose run --rm zap /zap/zap-full-scan.py \
  -t "http://app:3030/users/cursor-based?cursor=10&limit=20" \
  -r "full-cursor-report.html" \
  -z "-config api.key=${ZAP_API_KEY} -config spider.maxDepth=1 -config spider.contextName=pagination"

docker compose run --rm zap /zap/zap-full-scan.py \
  -t "http://app:3030/users/limit-offset?page=1&limit=20" \
  -r "full-limit-offset-report.html" \
  -z "-config api.key=${ZAP_API_KEY} -config spider.maxDepth=1 -config spider.contextName=pagination"
