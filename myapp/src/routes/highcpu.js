var express = require('express');
var { highcpuCounter } = require('../utils/prometheus');

var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
    var num1 = 1;
    var num2 = 1;
    var temp;
    for(var i=0; i<req.query.count; i++) {
        temp = num2;
        num2 = num1 + num2;
        num1 = temp;
    }

    let ts = Date.now();
    let date = new Date(ts);

    highcpuCounter.inc({ code: 200 });

    res.send({ 
        api: 'highcpu', 
        fib: num2, 
        date
    });
});

module.exports = router;
