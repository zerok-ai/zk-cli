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
  const res = http.get('http://a8e0a85d68de643199a7d49f14969f6c-480821876.us-east-2.elb.amazonaws.com/highload');
  check(res, {
    'verify homepage text': (r) =>
      r.body.includes('highload'),
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
