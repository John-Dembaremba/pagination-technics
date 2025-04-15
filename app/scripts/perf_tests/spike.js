import http from "k6/http";
import { sleep, check } from "k6";

export const options = {
  stages: [
    { duration: "4m", target: 500 }, // fast ramp-up to a high point
    { duration: "1m", target: 0 }, // quick ramp-down to 0 users
  ],
  thresholds: {
    "http_req_duration{endpoint:limit-offset}": ["p(95)<500"], // Expect delays
    "http_req_duration{endpoint:cursor-based}": ["p(95)<500"],
    "http_req_failed{endpoint:limit-offset}": ["rate<0.1"],
    "http_req_failed{endpoint:cursor-based}": ["rate<0.1"],
  },
};

export default function () {
  const test_type = "spike";

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
