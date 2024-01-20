# hammer

Hammer - Simple Image Compression Service.

## Dependencies

- Consul
- RabbitMQ
- [mini-fstore](https://github.com/curtisnewbie/mini-fstore)

## Updates

- Since v0.0.5, hammer nolonger uses libvip (external dependency) to compress images. It now uses [github.com/disintegration/gift](github.com/disintegration/gift) instead, it might be slower and more resource intensive, but it's written purely in Go.

## The Sequence

```mermaid
sequenceDiagram

    participant b as backend
    participant m as mini-fstore
    participant h as hammer

    b->>m:Upload image
    b--)h:(MQ) trigger image compression
    h->>m:download original file
    h->>h:compress image
    h->>m:upload compressed image
    h--)b:(MQ) reply file_id of the compressed image
```