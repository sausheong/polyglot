#include "czmq.h"
#include <uuid/uuid.h>
#define ROUTEID "GET/_/hello/c"

int main (void) {
    zctx_t *ctx = zctx_new ();
    void *responder = zsocket_new (ctx, ZMQ_REQ);

    char identity [37];
    uuid_t uuid;
    uuid_generate(uuid);
    uuid_unparse_lower(uuid, identity);
    
    zmq_setsockopt (responder, ZMQ_IDENTITY, identity, strlen (identity));
    zsocket_connect (responder, "tcp://localhost:4321");

    printf ("%s - %s responder ready\n", ROUTEID, identity);
    zstr_send(responder, ROUTEID);
    while (true) {
        char *msg = zstr_recv (responder);
        if (!msg) {
          printf ("No message received from broker");
          break;
        }        
        zstr_sendm (responder, ROUTEID);
        zstr_sendm (responder, "200");
        zstr_sendm (responder, "{\"Content-Type\": \"text/html\"}");
        zstr_send (responder, msg);
        zstr_free (&msg);
    }
    zctx_destroy (&ctx);
    return 0;
}