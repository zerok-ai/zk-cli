import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

var prealloc = {
	highload: 3000, // 11000,
	highmem:  1500, // 300,
	highcpu:  1500, // 300,
	lowload:  1500, // 300
};
var highcpuCount = 800;
var highmemCount = 80;


export const options = {
  scenarios: {
    highload: {
      executor: 'ramping-arrival-rate',
      exec: 'highload',
      preAllocatedVUs: prealloc.highload,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '3m', target: prealloc.highload },
        { duration: '30s', target: 0 }
      ]    
    },
    highmem: {
      executor: 'ramping-arrival-rate',
      exec: 'highmem',
      preAllocatedVUs: prealloc.highmem,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '3m', target: prealloc.highmem },
        { duration: '30s', target: 0 }
      ]
    },
    highcpu: {
      executor: 'ramping-arrival-rate',
      exec: 'highcpu',
      preAllocatedVUs: prealloc.highcpu,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '3m', target: prealloc.highcpu },
        { duration: '30s', target: 0 }
      ]
    },
    lowload: {
      executor: 'ramping-arrival-rate',
      exec: 'lowload',
      preAllocatedVUs: prealloc.lowload,
      stages: [
        { duration: '1m', target: 300 },
        { duration: '3m', target: prealloc.lowload },
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

const hostname = 'a3ae3226c4e37450ca10ff855f2fed15-453678147.us-east-2.elb.amazonaws.com';

export function highload () {
  const res = http.get('http://'+hostname+'/highload');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('highload'),
  });
  sleep(1);
}

export function highcpu() {
  const res = http.get('http://'+hostname+'/highcpu?count='+highcpuCount);
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('highcpu'),
  });
  sleep(1);
}

export function highmem() {
  const res = http.get('http://'+hostname+'/highmem?count='+highmemCount);
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('highmem'),
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


