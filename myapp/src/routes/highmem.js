var express = require('express');
var { highmemCounter } = require('../utils/prometheus');

var router = express.Router();
const fs = require('fs');
const dummyData = require('../utils/dummy.json');

function allocateMemory(size) {
    // Simulate allocation of bytes
    const numbers = size / 8;
    const arr = [];
    arr.length = numbers;
    for (let i = 0; i < numbers; i++) {
        arr[i] = i;
    }
    return arr;
}

const memoryLeakAllocations = [];
    
const field = "heapUsed";
const allocationStep = 1000 * 1024; // 1MB

/* GET home page. */
router.get('/', function(req, res, next) {
    var recCount = req.query.count || 1;
    var length = 0;
    var dataStore = [];
//     for(i=0; i<recCount; i++) {
//         dataStore[i] = [];
//         for(var j=0; j<dummyData.length; j++) {
// //            length+=dummyData[j].creditBalance;
// //            var newData = [...dummyData];
//             dataStore[i].push([...dummyData]);
//         }
//     }


    const allocation = allocateMemory(allocationStep);

    memoryLeakAllocations.push(allocation);

    const mu = process.memoryUsage();
    // # bytes / KB / MB / GB
    const gbNow = mu[field] / 1024 / 1024 / 1024;
    const gbRounded = Math.round(gbNow * 100) / 100;

    console.log(`Heap allocated ${gbRounded} GB`);

    let ts = Date.now();
    let date = new Date(ts);

    highmemCounter.inc({ code: 200 });

    res.send({ 
        api: 'highmem', 
        allocated_gbs: gbRounded, 
        date
    });
});

module.exports = router;