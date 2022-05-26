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
  const res = http.get('http://afcdb7908a66d4ebdad8531a29ccdc52-1591718493.us-east-2.elb.amazonaws.com/lowload');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('lowload'),
  });

  sleep(1);
}
