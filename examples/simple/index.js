'use strict';


module.exports = function(request, response) {
  setTimeout(() => response.write(`Hello: ${request.form.name}`), 70);
};
