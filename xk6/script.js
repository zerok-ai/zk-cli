import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';
import { Trend } from 'k6/metrics';

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
const scenarioMetrics = ['waiting', 'duration']

/* End scenario specs */
var myTrend = {};

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
    scenarioMetrics.forEach((metric) => {
	myTrend[element] = myTrend[element] || {};
    	myTrend[element][metric] = new Trend(`custom_${element}_${metric}`);
    })
    module.exports[element] = prepareExecFn(element);
    scenarios[element] = generateScenarioObj(element);
  });
  console.log(scenarios)
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
    scenarioMetrics.forEach((metric) => {
    	myTrend[scenarioName][metric].add(res.timings[metric], {tag: `${scenarioName}_${metric}`});
    })
    console.log(res);
    sleep(1);  
  }
}
