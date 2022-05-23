'use strict';
const express = require('express')
const helmet = require('helmet')
const os = require('os');
const app = express()



/**
 * Prometheus
 */

const apiMetrics = require('prometheus-api-metrics');
app.use(apiMetrics({
	metricsPrefix: "loadtest_"
}));
// const register = apiMetrics.register;
// // const client = apiMetrics.HttpMetricsCollector;
const client = require('prom-client');
// const register = client.register;

// Enable collection of default metrics
// client.collectDefaultMetrics({
// 	gcDurationBuckets: [0.001, 0.01, 0.1, 1, 2, 5], // These are the default buckets.
// });
// new client.Counter({
// 	name: 'scrape_counter',
// 	help: 'Number of scrapes (example of a counter with a collect fn)',
// 	collect() {
// 		// collect is invoked each time `register.metrics()` is called.
// 		this.inc();
// 	},
// });


const Counter = client.Counter;
const info1Counter = new Counter({
	name: 'info1_counter',
	help: 'info1 access counter',
	labelNames: ['code'],
});
const info2Counter = new Counter({
	name: 'info2_counter',
	help: 'info2 access counter',
	labelNames: ['code'],
});
const hc = new Counter({
	name: 'hc_counter',
	help: 'Health check access counter',
	labelNames: ['code'],
});
const defCounter = new Counter({
	name: 'root_counter',
	help: '/ access counter',
	labelNames: ['code'],
});

/**********/

// get hosts and ips
const nets = os.networkInterfaces();
const ipconfig = Object.create(null); // Or just '{}', an empty object
const hostname = os.hostname();
const port = 3000;
for (const name of Object.keys(nets)) {
	for (const net of nets[name]) {
		if (!ipconfig[name]) {
			ipconfig[name] = [];
		}
		ipconfig[name].push(net.address);
	}
}

app.get('/hc', (req, res) => {
	hc.inc({ code: 200 });
	res.send({"success": true});
});

app.get('/info1', (req, res) => 
{
	let ts = Date.now();
	let date = new Date(ts);

	let response = {
		ipconfig,
		hostname, 
		port,
		date,
		"api": "info1"
	};
	info1Counter.inc({ code: 200 });
	res.send(response);
})

app.get('/info2', (req, res) => 
{
	let ts = Date.now();
	let date = new Date(ts);

	let response = {
		ipconfig,
		hostname, 
		port,
		date,
		"api": "info2"
	};
	info2Counter.inc({ code: 200 });
	res.send(response);
})

// Setup server to Prometheus scrapes:
// app.get('/metrics', async (req, res) => {
// 	try {
// 		res.set('Content-Type', register.contentType);
// 		res.end(await register.metrics());
// 	} catch (ex) {
// 		res.status(500).end(ex);
// 	}
// });

app.get('/', (req, res) => {
	defCounter.inc({ code: 200 });
	res.send({"Nothing to show here": true});
});

app.use(helmet())

app.listen(port, () => {
	console.log('Server ready at http://'+ hostname + ':' + port);
	console.log('Supports Prometheus');
})
