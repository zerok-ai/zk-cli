import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 60 },
    { duration: '2m', target: 250 },
    { duration: '2m', target: 250 }
  ],
};

export default function () {
  const res = http.get('http://a8b53a31b008142e9b6239cfc81bef31-1654016547.us-east-2.elb.amazonaws.com/info1');
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
