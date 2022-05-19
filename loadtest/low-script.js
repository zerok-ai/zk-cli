import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  const res = http.get('http://a91d8bb2c6a8f404cbac9a4f85bda814-622073186.us-east-2.elb.amazonaws.com/info2');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('info2'),
  });

  sleep(1);
}
