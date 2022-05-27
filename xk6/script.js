import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';

/* scenario specs */
const preallocVUs = 2;
const maxVUs = 50;
const timeUnit = '1m';

const scenarioStages = {
  'highload' : [
    { duration: '1m', target: 100 },
    { duration: '3m', target: 300 },
    { duration: '30s', target: 0 }  
  ],
  'highmem' : [
    { duration: '1m', target: 100 },
    { duration: '3m', target: 300 },
    { duration: '30s', target: 0 }  
  ],
  'highcpu' : [
    { duration: '1m', target: 100 },
    { duration: '3m', target: 300 },
    { duration: '30s', target: 0 }  
  ],
  'lowload' : [
    { duration: '1m', target: 100 },
    { duration: '3m', target: 300 },
    { duration: '30s', target: 0 }  
  ]
}

const highcpuCount = 800; // Count variable to control CPU consumed by each highcpu API call.
const highmemCount = 80;  // Count variable to control Mem consumed by each highmem API call.

/* End scenario specs */
function generateScenarioObj(scenarioName) {
  return {
    executor: 'ramping-arrival-rate',
    exec: scenarioName,
    preAllocatedVUs: preallocVUs,
    timeUnit,
    maxVUs,
    startRate: scenarioStages[scenarioName][0].target,
    stages: scenarioStages[scenarioName]
  }
}

function generateScenarios() {
  var scenarios = {};
  Object.keys(scenarioStages).forEach(element => {
    module.exports[element] = prepareExecFn(element);
    scenarios[element] = generateScenarioObj(element);
  });
  return scenarios;
}

export const options = {
  scenarios: generateScenarios(),
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

const hostname = __ENV.MY_HOSTNAME;

function prepareExecFn(scenarioName) {
  return () => {
    const res = http.get('http://'+hostname+'/'+scenarioName);
    check(res, {
      'verify homepage text': (r) =>
        r.body.includes(scenarioName),
    });
    sleep(1);  
  }
}
