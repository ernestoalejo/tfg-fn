'use strict';


module.exports = function(request, response) {
  response.write('Hello: ${request.query.name}');
};
