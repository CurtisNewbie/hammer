package api

const (
	// event bus to trigger image compression
	CompressImageTriggerEventBus = "event.bus.hammer.image.compress.processing"

	// event bus to trigger video thumbnail generation
	GenVideoThumbnailTriggerEventBus = "event.bus.hammer.video.thumbnail.processing"
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

// Event sent to hammer to trigger an vidoe thumbnail generation.
type GenVideoThumbnailTriggerEvent struct {
	Identifier string // identifier
	FileId     string // file id from mini-fstore
	ReplyTo    string // event bus that will receive event about the generated video thumbnail.
}

// Event replied from hammer about the generated video thumbnail.
type GenVideoThumbnailReplyEvent struct {
	Identifier string // identifier
	FileId     string // file id from mini-fstore
}
