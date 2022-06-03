var express = require('express');
var { highcpuCounter } = require('../utils/prometheus');

var router = express.Router();

function fibonacci(nth) {
    var num1 = 1;
    var num2 = 1;
    var temp;
    for(var i=0; i<nth; i++) {
        temp = num2;
        num2 = num1 + num2;
        num1 = temp;
    }
    return num2;
}

/* GET home page. */
router.get('/', function(req, res, next) {    
    var recCount = (req.query.count || 10);
    var finalData = 1; // hackySack(jsondata, recCount);
    
    for(var i=0; i<recCount; i++) {
        finalData += 1;
        Math.random();
        // finalData += Math.random()*recCount; //fibonacci(i);
    }

    let ts = Date.now();
    let date = new Date(ts);

    highcpuCounter.inc({ code: 200 });
    res.send({ 
        api: 'highcpu',
        count: finalData,
        date
    });
});

module.exports = router;
