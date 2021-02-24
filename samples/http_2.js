import http from "k6/http";
import { check } from "k6";

if (__VU == 1) {
    http.setResponseCallback(http.expectedStatuses(200))
} else {
    http.setResponseCallback(http.expectedStatuses(300))
}
export default function () {
    console.log(JSON.stringify(http));
  check(http.get("https://test-api.k6.io/"), {
    "status is 200": (r) => r.status == 200,
    "protocol is HTTP/2": (r) => r.proto == "HTTP/2.0",
  });

}
