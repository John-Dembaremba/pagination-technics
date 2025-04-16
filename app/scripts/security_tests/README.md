
# Security Tests for Pagination App

This document outlines the security testing suite for the `pagination-app`, a Golang-based web application with REST API endpoints (`/users/cursor-based`, `/users/limit-offset`). The tests leverage [OWASP ZAP](https://www.zaproxy.org/) to identify vulnerabilities such as XSS, SQL injection, CSRF, SSRF, and misconfigurations. Results are saved as HTML reports in `app/zap-reports/` for review. The suite includes **baseline** and **active** scans, executed via the script `./scripts/security_tests/run_tests.sh`.

## Test Types

The security tests are designed to assess the application’s resilience against various attack vectors:

1. **Baseline Scan** (`zap-baseline.py`):
   - **Purpose**: Detects common vulnerabilities and misconfigurations with minimal impact.
   - **Details**: Performs passive scanning on `/users/cursor-based` and `/users/limit-offset`, analyzing headers, cookies, and responses for issues like missing security headers (e.g., `X-Content-Type-Options: nosniff`) or weak configurations (e.g., Spectre isolation).
   - **Goal**: Identify low-hanging issues (e.g., OWASP Top 10 misconfigurations) without stressing the app, ensuring no critical errors contribute to issues like `cursor-based` failures.

2. **Active Scan** (`zap-full-scan.py`):
   - **Purpose**: Tests for dynamic vulnerabilities by simulating attacks.
   - **Details**: Actively injects payloads into parameters like `cursor` and `limit` to detect XSS, SQL injection, CSRF, SSRF, and other exploits. Examples include `cursor=<script>alert(1)</script>` for XSS or `cursor=10; DROP TABLE users` for SQL injection.
   - **Goal**: Confirm the app resists injection attacks and handles errors securely, particularly for `cursor-based`, which shows HTTP 500 errors in k6 performance tests.

## Prerequisites

Before running the tests, ensure the following:

- **Docker and Docker Compose**: Installed to manage `app`, `zap`, and other services.
- **Project Setup**: Clone the repository and navigate to the `app/` directory:
  ```bash
  cd app
  ```
- **Environment Configuration**: Copy `env-example` to `.env` and generate a `ZAP_API_KEY` (see below).
- **ZAP Service**: Configured in `docker-compose.yml`:
  ```yaml
  zap:
    image: owasp/zap2docker-stable
    container_name: zap
    volumes:
      - ./zap-reports:/zap/wrk
    environment:
      - ZAP_PORT=8080
      - ZAP_API_KEY=${ZAP_API_KEY}
    networks:
      - pagination-app
    command: zap.sh -daemon -host 0.0.0.0 -port 8080 -config api.addrs.addr.name=.* -config api.key=${ZAP_API_KEY}
  ```

## Generating ZAP_API_KEY

The `ZAP_API_KEY` is a secure, random string required for ZAP’s API authentication. Generate it using `openssl`:

1. **Generate the Key**:
   ```bash
   openssl rand -hex 32
   ```
   Example output: `f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2`

2. **Add to `.env`**:
   ```bash
   cp env-example .env
   echo "ZAP_API_KEY=$(openssl rand -hex 32)" >> .env
   ```

3. **Verify `.env` Security**:
   Ensure `.env` is ignored by version control:
   ```bash
   grep ".env" .gitignore || echo ".env" >> .gitignore
   ```

## Running Security Tests

The security tests are automated via a script that runs ZAP baseline and active scans for both endpoints.

### Instructions

1. **Navigate to the Project Directory**:
   ```bash
   cd app
   ```

2. **Ensure Services Are Running**:
   Start `app` and `zap` services:
   ```bash
   docker compose up -d app zap
   ```

3. **Create Report Directory**:
   ```bash
   mkdir -p zap-reports
   ```

4. **Run All Tests**:
   Execute the security test suite:
   ```bash
   ./scripts/security_tests/run_tests.sh
   ```

   **What It Does**:
   - Runs baseline scans (`zap-baseline.py`) on `/users/cursor-based?cursor=10&limit=20` and `/users/limit-offset?page=1&limit=20` to detect headers and configuration issues.
   - Runs active scans (`zap-full-scan.py`) to test for XSS, SQL injection, CSRF, SSRF, and other vulnerabilities.
   - Saves reports to `app/zap-reports/`:
     - `baseline-cursor-report.html`
     - `baseline-limit-offset-report.html`
     - `full-cursor-report.html`
     - `full-limit-offset-report.html`

5. **Monitor Execution**:
   - Check terminal output for scan progress and errors.
   - View ZAP logs if issues arise:
     ```bash
     docker compose logs zap
     ```
   - Monitor `app` logs for errors (e.g., HTTP 500 on `cursor-based`):
     ```bash
     docker compose logs app
     ```

## Visualizing Results with ZAP Reports

ZAP generates HTML reports for each scan, providing detailed vulnerability analysis.

### Steps to Visualize

1. **Access Reports**:
   - List reports:
     ```bash
     ls app/zap-reports
     ```
   - Open in a browser:
     - `baseline-cursor-report.html`: Passive scan results for `cursor-based`.
     - `baseline-limit-offset-report.html`: Passive scan results for `limit-offset`.
     - `full-cursor-report.html`: Active scan results for `cursor-based`.
     - `full-limit-offset-report.html`: Active scan results for `limit-offset`.

2. **Analyze Findings**:
   - **Baseline Scans**:
     - Check for:
       - Missing headers (e.g., `X-Content-Type-Options: nosniff`).
       - Weak configurations (e.g., Spectre isolation warnings).
       - Severity: Low/Medium (e.g., `WARN-NEW` alerts).
     - Example: Previous scans flagged `X-Content-Type-Options Header Missing`.
   - **Active Scans**:
     - Look for:
       - **XSS**: Payloads like `<script>alert(1)</script>` in `cursor`.
       - **SQL Injection**: Attempts like `cursor=10' OR '1'='1`.
       - **SSRF**: URLs in `cursor` (e.g., `http://malicious.com`).
       - **CSRF**: Token issues (unlikely for GET APIs).
       - **Server Errors**: HTTP 500, explaining `cursor-based` k6 failures.
     - Severity: High/Medium (e.g., injection vulnerabilities).
   - Filter by severity (High, Medium, Low) to prioritize fixes.

3. **Correlate with k6 Metrics** (Optional):
   - If `cursor-based` shows HTTP 500 errors, cross-check with Grafana:
     - Open `http://localhost:3000`.
     - Query:
       ```flux
       from(bucket: "k6")
         |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
         |> filter(fn: (r) => r["_measurement"] == "http_req_failed")
         |> filter(fn: (r) => r["endpoint"] == "cursor-based")
         |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
       ```
   - High `http_req_failed` may align with ZAP’s server error alerts.

4. **Troubleshoot Issues**:
   - **No Alerts**: Run a manual GUI scan:
     ```bash
     docker compose up -d zap
     ```
     Open `http://localhost:8080`, target `http://app:3030/users/cursor-based?cursor=10&limit=20`, and run Active Scan.
   - **cursor-based Failures**: If HTTP 500 persists, test:
     ```bash
     docker compose run --rm k6 sh -c "curl -v http://app:3030/users/cursor-based?cursor=10&limit=20"
     ```
     Review `app/internal/api/cursor.go` or `app/internal/repo/db.go`.

## Notes

- **Scan Configuration**: Update `run_tests.sh` if endpoints change (e.g., `cursor` vs. `after`).
- **Resource Limits**: Active scans may stress `app`. Monitor:
  ```bash
  docker stats app
  ```
- **Security**: Keep `ZAP_API_KEY` unique and secure in `.env`.
- **Known Issues**:
  - `cursor-based` fails in k6 `load`, `stress`, and `spike` tests (100% `http_req_failed`). Active scans may reveal SQL injection or logic errors.
  - Previous warnings: `X-Content-Type-Options Header Missing`, `Insufficient Site Isolation Against Spectre`.
- **Customization**: Add API scans with `bruno.json`:
  ```bash
  docker compose run --rm zap /zap/zap-api-scan.py -t "/zap/wrk/bruno.json" -f openapi -r "api-report.html" -z "-config api.key=${ZAP_API_KEY}"
  ```

## Example Commands

- Run a single baseline scan:
  ```bash
  docker compose run --rm zap /zap/zap-baseline.py -t "http://app:3030/users/cursor-based?cursor=10&limit=20" -r "baseline-cursor-report.html" -z "-config api.key=${ZAP_API_KEY}"
  ```

- Run a single active scan:
  ```bash
  docker compose run --rm zap /zap/zap-full-scan.py -t "http://app:3030/users/cursor-based?cursor=10&limit=20" -r "full-cursor-report.html" -z "-config api.key=${ZAP_API_KEY}"
  ```

- Check ZAP setup:
  ```bash
  docker compose logs zap
  ```

- Verify reports:
  ```bash
  ls app/zap-reports
  ```
