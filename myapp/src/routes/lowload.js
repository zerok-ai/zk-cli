var express = require('express');

var { lowloadCounter } = require('../utils/prometheus');
var { ipconfig, hostname } = require('../utils/sysinfo'); 

var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  let ts = Date.now();
  let date = new Date(ts);

  let response = {
      ipconfig,
      hostname, 
      date,
      "api": "lowload"
  };

  lowloadCounter.inc({ code: 200 });

  res.send(response);
});

module.exports = router;

