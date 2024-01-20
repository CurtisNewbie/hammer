package api

const (
	// Name of the event bus that receives evnets to trigger image compression.
	CompressImageTriggerEventBus = "event.bus.hammer.image.compress.processing"
)

// Event sent to hammer to trigger an image compression.
type ImageCompressTriggerEvent struct {
	Identifier string // identifier
	FileId     string // file id from mini-fstore
	ReplyTo    string // event bus that will receive event about the compressed image
}

// Event replied from hammer about the compressed image.
type ImageCompressReplyEvent struct {
	Identifier string // identifier
	FileId     string // file id from mini-fstore
}
