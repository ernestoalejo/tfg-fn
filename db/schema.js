
r.db('rethinkdb').table('users').insert({id: 'fn', password: 'fnpass'});

r.dbCreate('fn');
r.db('fn').grant('fn', {read: true, write: true});

r.db('fn').tableCreate('functions', {primary_key: 'Name'});
r.db('fn').tableCreate('secrets', {primary_key: 'Name'});
