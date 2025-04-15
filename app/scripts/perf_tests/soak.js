import http from "k6/http";
import { sleep, check } from "k6";

export const options = {
  stages: [
    { duration: "5m", target: 100 }, // traffic ramp-up from 1 to 100 users over 5 minutes.
    { duration: "1h", target: 100 }, // stay at 100 users for 30 mins!!!
    { duration: "3m", target: 0 }, // ramp-down to 0 users
  ],
  thresholds: {
    "http_req_duration{endpoint:limit-offset}": ["p(95)<250"],
    "http_req_duration{endpoint:cursor-based}": ["p(95)<250"],
    "http_req_failed{endpoint:limit-offset}": ["rate<0.02"],
    "http_req_failed{endpoint:cursor-based}": ["rate<0.02"],
  },
};

export default function () {
  const test_type = "soak";

  const limitOffset = http.get(
    "http://app:3030/users/limit-offset?page=1&limit=20",
    {
      tags: { endpoint: "limit-offset", test_type: test_type },
    },
  );
  check(limitOffset, { "limit-offset status is 200": (r) => r.status === 200 });

  const cursorBased = http.get(
    "http://app:3030/users/cursor-based?cursor=22&limit=20",
    {
      tags: { endpoint: "cursor-based", test_type: test_type },
    },
  );
  check(cursorBased, { "cursor-based status is 200": (r) => r.status === 200 });

  sleep(1);
}
