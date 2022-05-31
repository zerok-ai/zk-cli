import { check } from 'k6';
import http from 'k6/http';
import { sleep } from 'k6';
import { Trend } from 'k6/metrics';

/* scenario specs */
const preallocVUs = 1200;
const maxVUs = 1200;
const timeUnit = '1m';

const scenarioStages = {
/*/
  'highmem' : [
    { duration: '3m', target: 999999999 },
    { duration: '3m', target: 999999999 },
    { duration: '5s', target: 999999999 }  
  ],
/**/
  'highcpu' : [
    { duration: '1m', target: 10000 },
    { duration: '3m', target: 10000 },
    { duration: '30s', target: 10000 }  
  ],
/*/	
  'highload' : [
    { duration: '1m', target: 10000 },
    { duration: '1m', target: 10000 },
    { duration: '1m', target: 10000 },
    { duration: '1m', target: 10000 }
  ],
/**/
  'lowload' : [
    { duration: '1m', target: 1000 },
    { duration: '3m', target: 3000 },
    { duration: '30s', target: 0 }  
  ]
/**/
}

const verticalScaleCount = {
  // Count variable to control Mem consumed by each highmem API call.
  'highmem': 1200,
   // Count variable to control CPU consumed by each highcpu API call.
  'highcpu': 50
}

const scenarioMetrics = ['waiting', 'duration']

/* End scenario specs */
var myTrend = {};

function generateScenarioObj(scenarioName) {
  return {
    executor: 'constant-arrival-rate',
    exec: scenarioName,
    preAllocatedVUs: preallocVUs,
    timeUnit,
    duration: '4m',
    maxVUs,
    rate: scenarioStages[scenarioName][0].target,
//    startRate: scenarioStages[scenarioName][0].target,
//    stages: scenarioStages[scenarioName]
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
  return scenarios;
}

export const options = {
  noConnectionReuse: true,
  scenarios: generateScenarios(),
  VUs: preallocVUs,
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
    const res = http.get('http://'+hostname+'/'+scenarioName+'?count='+verticalScaleCount[scenarioName]);
    check(res, {
      'verify homepage text': (r) =>
        r.body.includes(scenarioName),
    });
    scenarioMetrics.forEach((metric) => {
    	myTrend[scenarioName][metric].add(res.timings[metric], {tag: `${scenarioName}_${metric}`});
    })
    sleep(1);  
  }
}
