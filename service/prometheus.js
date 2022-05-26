const client = require('prom-client');

const Counter = client.Counter;
const highloadCounter = new Counter({
    name: 'highload_counter',
    help: 'highload access counter',
    labelNames: ['code'],
});
const lowloadCounter = new Counter({
    name: 'lowload_counter',
    help: 'lowload access counter',
    labelNames: ['code'],
});
const hcCounter = new Counter({
    name: 'hc_counter',
    help: 'Health check access counter',
    labelNames: ['code'],
});
const defCounter = new Counter({
    name: 'root_counter',
    help: 'access counter',
    labelNames: ['code'],
});


module.exports = {
    highloadCounter,
    lowloadCounter,
    hcCounter,
    defCounter,
};
 