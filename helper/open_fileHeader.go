package helper

import "mime/multipart"

func OpenFileHeader(fileHeader *multipart.FileHeader) multipart.File {
	file, err := fileHeader.Open()
	if err != nil {
		return nil
	}

	return file
}
