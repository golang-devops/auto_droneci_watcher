package config

import (
	"fmt"
	"strings"
)

//ParseSecretLine will parse the secret line according to its expected format/layout
func ParseSecretLine(line string) (*Secret, error) {
	//Each in format 'docker/image1,docker/image2 SECRET_KEY=secret_value'
	firstSpaceIndex := strings.Index(line, " ")
	if firstSpaceIndex == -1 {
		return nil, fmt.Errorf("Expected a space character in the secret line '%s'", line)
	}

	imagesPart := strings.TrimSpace(line[:firstSpaceIndex])
	keyValuePart := strings.TrimSpace(line[firstSpaceIndex+1:])

	imageList := strings.Split(imagesPart, ",")

	indexEqualInKeyVal := strings.Index(keyValuePart, "=")
	if indexEqualInKeyVal == -1 {
		return nil, fmt.Errorf("Expected an '=' character in the secret line '%s' (was looking for it in the part '%s')", line, keyValuePart)
	}

	key := keyValuePart[:indexEqualInKeyVal]
	val := keyValuePart[indexEqualInKeyVal+1:]

	return &Secret{
		Key:    key,
		Value:  val,
		Images: imageList,
	}, nil
}
