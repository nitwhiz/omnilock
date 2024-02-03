# omnilock

A zero dependency locker.

## How does it work?

This locker uses TCP connections to keep track of locks.
This way, a lock is freed as soon as the locking process ends (some way or another) or it's unlocked explicitly.

## Protocol

1. Open a TCP connection to the omnilock server
2. Send a lock command through the established connection (e.g. `lock my_lock`)
3. Read the response, it may be either `success`, `failed` or `error: something went wrong`
4. Do your things on the client side
5. Either unlock the lock explicitly: `unlock my_lock` or just let the process exit so the OS closes the connection

Note: All commands and responses end with `\n` (newline).

## Usage as docker image

Start the server with:

```shell
docker run --rm -it -p 7194:7194 ghcr.io/nitwhiz/omnilock:latest
```

And connect to it via port 7194.
