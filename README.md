# hammer

Hammer - Simple Image Compression Service.

### The Sequence


```mermaid
sequenceDiagram

    participant b as browser (bolobao)
    participant m as mini-fstore
    participant v as vfm
    participant h as hammer

    b->>m:Upload file
    m->>b:return file_id (the fake one)
    b->>v:Create file recrod
    v->>m:get the real file_id
    m->>v:actual file info
    v->>v:check if it's an image
    v--)h:(MQ) trigger image compression
    h->>m:download original file
    m->>h:original image
    h->>h:compress image
    h->>m:upload compressed image
    h--)v:(MQ) file_id of the compressed image
```