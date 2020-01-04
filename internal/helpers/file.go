package helpers

// GetExtensionFile ..
func GetExtensionFile(contentType string) (extensionFile string) {
	switch contentType {
	case "image/jpeg":
		extensionFile = "jpeg"
	case "image/jpg":
		extensionFile = "jpg"
	case "image/png":
		extensionFile = "png"
	case "image/svg+xml":
		extensionFile = "svg"
	case "image/bmp":
		extensionFile = "bmp"
	case "image/gif":
		extensionFile = "gif"
	case "image/vnd.microsoft.icon":
		extensionFile = "ico"
	case "image/tiff":
		extensionFile = "tiff"
	case "image/webp":
		extensionFile = "webp"
	default:
		extensionFile = "png"
	}
	return extensionFile
}
