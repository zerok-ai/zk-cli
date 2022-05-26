const client = require('prom-client');

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
    info1Counter,
    info2Counter,
    hcCounter,
    defCounter,
};
 