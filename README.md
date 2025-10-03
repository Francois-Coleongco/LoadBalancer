# THE PLAN

// start the http proxy
// read in new requests passing off to servers based on an injected load balancing session cookie for backend agnosticism
// set this cookie in the response then the client will send it every time automatically. then just route using that cookie
// every time the user sends a request, the user will have their session info (specific to the backend whatever it may be) stored in a Redis instance.
// if the load balancer decides it is optimal, they may pipe a subsequent request by this user to another server safely as the session data is stored in redis, not local to a backend server

// start redis in a docker container btw

// need to track inactivity to delete a session from the redis store

