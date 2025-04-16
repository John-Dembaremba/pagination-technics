
# Performance Tests for Pagination App

This document outlines the performance testing suite for the `pagination-app`, a Golang-based web application with REST API endpoints (`/users/cursor-based`, `/users/limit-offset`). The tests use [k6](https://k6.io/) to evaluate the application’s performance under various conditions, with results stored in InfluxDB and visualized using Grafana. The suite includes **load**, **stress**, **spike**, **soak**, and **breakpoint** tests, executed via a script (`run_all_test.sh`).

## Test Types

The performance tests are designed to assess different aspects of the application’s behavior:

1. **Load Test** (`load.js`):
   - **Purpose**: Measures performance under expected traffic conditions.
   - **Details**: Simulates a steady number of virtual users (VUs) accessing `/users/cursor-based` and `/users/limit-offset` to verify response times, throughput, and error rates.
   - **Goal**: Ensure the app handles typical load (e.g., 50 VUs) with `http_req_duration` < 500ms and `http_req_failed` < 1%.

2. **Stress Test** (`stress.js`):
   - **Purpose**: Tests the app’s limits by increasing load beyond normal capacity.
   - **Details**: Gradually ramps up VUs (e.g., from 50 to 200) to identify performance degradation, bottlenecks, or failures (e.g., HTTP 500 errors).
   - **Goal**: Find the point where response times exceed acceptable thresholds or errors spike.

3. **Spike Test** (`spike.js`):
   - **Purpose**: Evaluates handling of sudden traffic surges.
   - **Details**: Rapidly increases VUs (e.g., from 10 to 500 in seconds) to simulate a traffic spike, then scales down.
   - **Goal**: Confirm the app recovers gracefully without crashing, with minimal error rates.

4. **Soak Test** (`soak.js`):
   - **Purpose**: Assesses stability over extended periods.
   - **Details**: Runs a moderate load (e.g., 50 VUs) for hours to detect memory leaks, database issues, or resource exhaustion.
   - **Goal**: Ensure consistent performance with no degradation over time.

5. **Breakpoint Test** (`breakpoint.js`):
   - **Purpose**: Determines the app’s breaking point by continuously increasing load.
   - **Details**: Incrementally adds VUs until the app fails (e.g., `http_req_failed` > 5% or `http_req_duration` > 2s).
   - **Goal**: Identify maximum capacity and failure modes (e.g., database timeouts, CPU saturation).

## Prerequisites

Before running the tests, ensure the following:

- **Docker and Docker Compose**: Installed to manage `app`, `influxdb`, `grafana`, and `k6` services.
- **Project Setup**: Clone the repository and navigate to the `app/` directory:
  ```bash
  cd app
  ```
- **Environment Configuration**: Copy `env-example` to `.env` and update if needed:
  ```bash
  cp env-example .env
  ```
  - Ensure `INFLUXDB_TOKEN`, `GRAFANA_ADMIN_PASSWORD`, etc., are set.
- **Grafana and InfluxDB**: Running to store and visualize k6 metrics (configured in `docker-compose.yml`).

## Running Performance Tests

The performance tests are automated via a script that runs all k6 tests sequentially.

### Instructions

1. **Navigate to the Project Directory**:
   ```bash
   cd app
   ```

2. **Ensure Services Are Running**:
   Start `app`, `influxdb`, `grafana`, and other dependencies:
   ```bash
   docker-compose up -d app influxdb grafana
   ```

3. **Run All Tests**:
   Execute the test suite using `run_all_test.sh`:
   ```bash
   ./scripts/perf_tests/run_all_test.sh
   ```

   **What It Does**:
   - Builds the `k6` Docker image.
   - Starts `app`, `influxdb`, and `grafana` (if not already running).
   - Runs `load.js`, `stress.js`, `spike.js`, `soak.js`, and `breakpoint.js` in sequence.
   - Outputs metrics to InfluxDB (`k6` bucket).

4. **Monitor Execution**:
   - Check terminal output for test progress and errors.
   - View container logs if issues arise:
     ```bash
     docker-compose logs k6
     docker-compose logs app
     ```

## Visualizing Results with Grafana (k6 UI)

k6 metrics are stored in InfluxDB and visualized using Grafana for real-time analysis.

### Steps to Visualize

1. **Access Grafana**:
   - Open `http://localhost:3000` in a browser.
   - Log in:
     - Username: `admin`
     - Password: Set in `.env` (`GRAFANA_ADMIN_PASSWORD`).

2. **Configure Data Source**:
   - If not already set, add InfluxDB as a data source:
     - Go to **Configuration > Data Sources > Add data source**.
     - Select **InfluxDB**.
     - Set:
       - URL: `http://influxdb:8086`
       - Bucket: `k6`
       - Token: From `.env` (`INFLUXDB_TOKEN`).
     - Save and test.

3. **Import k6 Dashboard**:
   - Use a pre-configured k6 dashboard or create one:
     - Go to **Dashboards > Import**.
     - Use a k6 template (e.g., ID `2587` from Grafana’s dashboard library) or import `observability/dashboards/grafana/go-runtime.json` if customized.
   - Alternatively, create a new dashboard:
     - Add a panel with this Flux query for `cursor-based` errors:
       ```flux
       from(bucket: "k6")
         |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
         |> filter(fn: (r) => r["_measurement"] == "http_req_failed")
         |> filter(fn: (r) => r["endpoint"] == "cursor-based")
         |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
       ```

4. **Analyze Metrics**:
   - View key metrics in Grafana:
     - **http_req_duration**: Response times for each endpoint.
     - **http_req_failed**: Error rates (critical for `cursor-based`).
     - **vus**: Virtual users over time.
   - Filter by `test_type` (`load`, `stress`, `spike`, `soak`, `breakpoint`).
   - Example: For `load` test, check if `http_req_duration` < 500ms and `http_req_failed` < 1%.

5. **Troubleshoot Issues**:
   - If `cursor-based` shows high `http_req_failed` (e.g., 100% in `load`, `stress`, `spike`):
     - Check `app` logs:
       ```bash
       docker-compose logs app
       ```
     - Query InfluxDB directly:
       ```bash
       docker exec -it influxdb influx query 'from(bucket:"k6") |> range(start:-24h) |> filter(fn: (r) => r.endpoint == "cursor-based") |> limit(n:10)'
       ```

## Notes

- **Test Configuration**: Adjust VUs, durations, or thresholds in `load.js`, `stress.js`, etc., based on your app’s capacity (e.g., `app/docker_builds/app/Dockerfile` resources).
- **Resource Limits**: Monitor `docker stats` during tests to avoid overloading `app` or `influxdb`.
- **Security**: Ensure `.env` is in `.gitignore` to protect secrets.
- **Known Issues**:
  - `cursor-based` fails in `load`, `stress`, and `spike` tests. Check ZAP reports (`app/zap-reports/`) for security issues or `app/internal/api/cursor.go` for logic errors.
- **Customization**: Add new tests by creating `*.js` files in `scripts/perf_tests/` and updating `run_all_test.sh`.

## Example Commands

- Run a single test (e.g., `load`):
  ```bash
  docker-compose -f docker-compose.yml run --rm k6 run scripts/perf_tests/load.js
  ```

- Check Grafana setup:
  ```bash
  docker-compose logs grafana
  ```

- View InfluxDB data:
  ```bash
  docker exec -it influxdb influx query 'from(bucket:"k6") |> range(start:-24h) |> limit(n:10)'
  ```
