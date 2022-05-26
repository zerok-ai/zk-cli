var express = require('express');
var { highmemCounter } = require('../utils/prometheus');

var router = express.Router();
var jsondata = require('../utils/dummy.json');

function hackySack(dummyData, recCount) {
    if(recCount == 0) {
        return {};
    }
    var newData = JSON.parse(JSON.stringify(dummyData));
    newData.depth = recCount;
    newData.child = JSON.parse(JSON.stringify(hackySack(newData, recCount - 1)));
    return newData;
}

/* GET home page. */
router.get('/', function(req, res, next) {    
    var recCount = req.query.count || 10;
    var finalData = hackySack(jsondata, recCount);

    let ts = Date.now();
    let date = new Date(ts);

    highmemCounter.inc({ code: 200 });
    res.send({ 
        api: 'highmem',
        count: finalData.depth,
        date
    });
});

module.exports = router;
