import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 300 },
    { duration: '2m', target: 900 },
    { duration: '2m', target: 1500 },
    { duration: '2m', target: 2000 },
    { duration: '3m', target: 2300 },
    { duration: '6m', target: 2000 },
    { duration: '2m', target: 1000 },
    { duration: '1m', target: 500 },
    { duration: '1m', target: 0 }
  ],
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

export default function () {
  const res = http.get('http://a1268da6f6ffb4780b47832ec41d452c-1058680272.us-east-2.elb.amazonaws.com/info1');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('info1'),
  });
  sleep(1);
}

// export function handleSummary(data) {
//   console.log('Preparing the end-of-test summary...');

//   //Can change 7 to 2 for longer results.
//   let r = (Math.random() + 1).toString(36).substring(7);
//   const filepath = "results/"+r+".json";
//   console.log(filepath);
//   return {
//     filepath: JSON.stringify(data)
//   }
// }
