var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
const apiMetrics = require('prometheus-api-metrics');

var indexRouter = require('./routes/index');
var healthCheckRouter = require('./routes/hc');
var lowloadRouter = require('./routes/lowload');
var highloadRouter = require('./routes/highload');
var highcpuRouter = require('./routes/highcpu');
var highmemRouter = require('./routes/highmem');

var app = express();

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(apiMetrics({metricsPrefix: "loadtest_"}));

app.use('/', indexRouter);
app.use('/hc', healthCheckRouter);
app.use('/lowload', lowloadRouter);
app.use('/highload', highloadRouter);
app.use('/highcpu', highcpuRouter);
app.use('/highmem', highmemRouter);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.send(err);
});

module.exports = app;
