"use strict";

var system = require('system');
var page = require('webpage').create();

if (system.args.length === 1) {
  console.log(system.args[1] + ' [USERNAME] [PASSWORD] [REQUEST URL] [OUTPUT FILE]');
  phantom.exit();
}

// authenticate
page.settings.userName = system.args[1];
page.settings.password = system.args[2];
var pageUrl = system.args[3];
var outputPath = system.args[4];

console.log(pageUrl, outputPath);

page.viewportSize = {width: 1280, height: 768};

page.open(pageUrl, function() {
  page.evaluate(function() {
    var r = document.getElementById('ribbon');
    if (r) {
      r.remove();
    }
  });
  waitFor(
      function() {
        return document.body.textContent.indexOf("Loading") === -1
      },
      function() {
        setTimeout(
            function() {
              page.render(outputPath);
              phantom.exit();
            }, 2000)
      },
      60000
  )
});


/**
 * taken from the phantomJS examples directory
 */
function waitFor(testFx, onReady, timeOutMillis) {
  var maxtimeOutMillis = timeOutMillis ? timeOutMillis : 3000, //< Default Max Timout is 3s
      start = new Date().getTime(),
      condition = false,
      interval = setInterval(function() {
        if ((new Date().getTime() - start < maxtimeOutMillis) && !condition) {
          // If not time-out yet and condition not yet fulfilled
          condition = (typeof(testFx) === "string" ? eval(testFx) : testFx()); //< defensive code
        } else {
          if (!condition) {
            // If condition still not fulfilled (timeout but condition is 'false')
            console.log("'waitFor()' timeout");
            phantom.exit(1);
          } else {
            // Condition fulfilled (timeout and/or condition is 'true')
            console.log("'waitFor()' finished in " + (new Date().getTime() - start) + "ms.");
            typeof(onReady) === "string" ? eval(onReady) : onReady(); //< Do what it's supposed to do once the condition is fulfilled
            clearInterval(interval); //< Stop this interval
          }
        }
      }, 250); //< repeat check every 250ms
}
