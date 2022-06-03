var express = require('express');
var { highmemCounter } = require('../utils/prometheus');

var router = express.Router();
const fs = require('fs');
const dummyData = require('../utils/dummy.json');

function allocateMemory(size) {
    // Simulate allocation of bytes
    // return new ArrayBuffer(size);

    const numbers = size;
    const arr = [];
    arr.length = numbers;
    for (let i = 0; i < numbers; i++) {
        arr[i] = 0;
    }
    return arr;
}

// const memoryLeakAllocations = [];
    
const field = "heapUsed";

/* GET home page. */
router.get('/', function(req, res, next) {
    var recCount = req.query.count || 1;
    var length = 0;
    var dataStore = [];
    const allocationStep = recCount * 1024 * 1024; // recCount = 1000 ==> 1MB

    const allocation = Buffer.allocUnsafe(allocationStep).fill(0);

    // memoryLeakAllocations.push(allocation); // To induce leak

    const mu = process.memoryUsage();
    // // # bytes / KB / MB / GB
    // const gbNow = mu[field] / 1024 / 1024 / 1024;
    // const gbRounded = Math.round(gbNow * 1000) / 1000;

    // console.log(`Heap allocated ${gbRounded} GB`);

    let ts = Date.now();
    let date = new Date(ts);

    highmemCounter.inc({ code: 200 });

    res.send({ 
        api: 'highmem', 
        allocated_gbs: mu, 
        date
    });
});

module.exports = router;