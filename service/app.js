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
	res.send({"Nothing to show here": true});
});

app.use(helmet())

app.listen(3000, () => 
	console.log('Server ready at http://'+ hostname + ':3000'));