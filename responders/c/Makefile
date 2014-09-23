# Assuming czmq is include directory is in ../../../czmq/include
CC=gcc
CFLAGS=-I../../../czmq/include -lczmq -lzmq

responder:
	$(CC) $(CFLAGS) responder.c -o responder

clean:
	rm -rf *o responder