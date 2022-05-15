import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  const res = http.get('http://a0fda3b0dbd5d409a90c2e1022c55741-1376665223.us-east-2.elb.amazonaws.com:3000');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('Hello Mudit Mathur! How are you?'),
  });
  sleep(1);
}
