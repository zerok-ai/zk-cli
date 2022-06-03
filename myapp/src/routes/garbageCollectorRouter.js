var express = require('express');
var router = express.Router();

/* GET users listing. */
router.get('/', function(req, res, next) {
  try {
    if (global.gc) {global.gc();}
    res.send("GC done")
  } catch (e) {
    console.log(e);
    res.send("GC failed")
  }
});

module.exports = router;
