'use strict';
const express = require('express')
const helmet = require('helmet')
const app = express()

const apiMetrics = require('prometheus-api-metrics');
app.use(apiMetrics({
    metricsPrefix: "loadtest_"
}));

var promCounters = require('./prometheus');
var sysInfo = require('./sysinfo');
 
const port = 3000;

app.use(helmet());
app.listen(port, () => {
	console.log('Server ready at http://'+ sysInfo.hostname + ':' + port);
	console.log('Supports Prometheus');
})

/* APIs */
app.get('/hc', (req, res) => {
	promCounters.hcCounter.inc({ code: 200 });
	res.send({"success": true});
});

app.get('/info1', (req, res) => 
{
	let ts = Date.now();
	let date = new Date(ts);

	let response = {
		ipconfig: sysInfo.ipconfig,
		hostname: sysInfo.hostname, 
		port,
		date,
		"api": "info1"
	};
	promCounters.info1Counter.inc({ code: 200 });
	res.send(response);
})

app.get('/info2', (req, res) => 
{
	let ts = Date.now();
	let date = new Date(ts);

	let response = {
		ipconfig: sysInfo.ipconfig,
		hostname: sysInfo.hostname, 
		port,
		date,
		"api": "info2"
	};
	promCounters.info2Counter.inc({ code: 200 });
	res.send(response);
})

app.get('/', (req, res) => {
	promCounters.defCounter.inc({ code: 200 });
	res.send({"Nothing to show here": true});
});  

