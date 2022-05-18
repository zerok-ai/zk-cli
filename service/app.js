'use strict';
const express = require('express')
const helmet = require('helmet')
const os = require('os');

const app = express()

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

// Prometheus
const client = require('prom-client');
const collectDefaultMetrics = client.collectDefaultMetrics;
// Probe every 5th second.
collectDefaultMetrics({ timeout: 5000 });
const counter = new client.Counter({
  name: 'info1 - node_request_operations_total',
  help: 'The total number of processed requests'
});
const histogram = new client.Histogram({
  name: 'info1 - node_request_duration_seconds',
  help: 'Histogram for the duration in seconds.',
  buckets: [1, 2, 5, 6, 10]
});


app.get('/hc', (req, res) => {
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
	
	//Simulate a sleep
	var start = new Date()
	var simulateTime = 1000

	setTimeout(function(argument) {
		// execution time simulated with setTimeout function
	    var end = new Date() - start
	    histogram.observe(end / 1000); //convert to seconds
	}, simulateTime)

	counter.inc();

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
	
	res.send(response);
})

app.get('/', (req, res) => {
	counter.inc();
	res.send({"Nothing to show here": true});
});

app.use(helmet())

// Metrics endpoint
app.get('/metrics', (req, res) => {
  res.set('Content-Type', client.register.contentType)
  res.end(client.register.metrics())
})

app.listen(port, () => 
	console.log('Server ready at http://'+ hostname + ':' + port));