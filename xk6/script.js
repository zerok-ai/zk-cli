import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';


export const options = {
  scenarios: {
    highload: {
      executor: 'ramping-arrival-rate',
      exec: 'highload',
      preAllocatedVUs: 11000,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '1m', target: 900 },
        { duration: '30s', target: 10000 },
        { duration: '1m', target: 1000 },
        { duration: '1m', target: 0 }
      ]    
    },
    lowload: {
      executor: 'ramping-arrival-rate',
      exec: 'lowload',
      preAllocatedVUs: 300,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '3m', target: 300 },
        { duration: '30s', target: 0 }
      ]
    }
  },
  ext: {
    loadimpact: {
      apm: [
        {
          provider: 'prometheus',
          remoteWriteURL: 'http://localhost:9090/api/v1/write',
          includeDefaultMetrics: true,
          includeTestRunId: true,
          resampleRate: 3,
        },
      ],
    },
  },
};

const hostname = 'aa7db044357554675bdf584d6d1fdffc-1919935816.us-east-2.elb.amazonaws.com';

export function highload () {
  const res = http.get('http://'+hostname+'/highload');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('highload'),
  });
  sleep(1);
}

export function lowload () {
  const res = http.get('http://'+hostname+'/lowload');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('lowload'),
  });
  sleep(1);
}


