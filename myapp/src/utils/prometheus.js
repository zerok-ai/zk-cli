/**
 * Prometheus
 */

const apiMetrics = require('prometheus-api-metrics');
const client = require('prom-client');

const Counter = client.Counter;
const highloadCounter = new Counter({
    name: 'highload_counter',
    help: 'highload access counter',
    labelNames: ['code'],
});
const highcpuCounter = new Counter({
    name: 'highcpu_counter',
    help: 'highcpu access counter',
    labelNames: ['code'],
});
const highmemCounter = new Counter({
    name: 'highmem_counter',
    help: 'highmem access counter',
    labelNames: ['code'],
});
const lowloadCounter = new Counter({
    name: 'lowload_counter',
    help: 'lowload access counter',
    labelNames: ['code'],
});


module.exports = {
    highloadCounter,
    lowloadCounter,
    highmemCounter,
    highcpuCounter
};
 