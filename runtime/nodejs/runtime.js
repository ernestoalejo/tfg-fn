'use strict';

var http = require('http');
var url = require('url');


var server = http.createServer((request, response) => {
  var body = '';
  request.on('data', data => body += data);
  request.on('end', () => {
    var context = JSON.parse(body);
    require(context.call)(context.request, {
      write: data => response.end(data),
    });
  });
});
server.listen(50050);
