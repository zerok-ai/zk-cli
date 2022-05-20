import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 30 },
    { duration: '2m', target: 100 },
    { duration: '2m', target: 100 }
  ],
};

export default function () {
  const res = http.get('http://af73403b39a704289af7dc2014fd202b-900425903.us-east-2.elb.amazonaws.com/info2');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('info2'),
  });

  sleep(1);
}
