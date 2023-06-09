package hammer

import (
	"fmt"
	"os"

	"github.com/h2non/bimg"
)

// Compress image
func CompressImage(file string, output string) error {
	buffer, err := bimg.Read(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	thumbnail, err := bimg.NewImage(buffer).Thumbnail(256)
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail, %v", err)
	}

	if err = bimg.Write(output, thumbnail); err != nil {
		return fmt.Errorf("failed to write thumbnail file, %v", err)
	}
	return nil
}
