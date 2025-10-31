package buildeventrecorder

import (
	"strconv"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
)

func getFileDigestFromBesFile(file *bes.File) *string {
	var digest string
	if file == nil {
		return nil
	}
	if file.GetDigest() != "" {
		digest = file.GetDigest()
		return &digest
	}
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		if len(stringArr) < 2 {
			return nil
		}
		digest = stringArr[len(stringArr)-2]
		return &digest
	}
	return nil
}

func getFileSizeBytesFromBesFile(file *bes.File) *int64 {
	var sizeBytes int64
	if file == nil {
		return nil
	}
	if file.GetLength() != 0 {
		sizeBytes = file.GetLength()
		return &sizeBytes
	}
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		sizeBytes, err := strconv.ParseInt(stringArr[len(stringArr)-1], 10, 64)
		if err != nil {
			return nil
		}
		return &sizeBytes
	}
	return nil
}

func getFileDigestFunctionFromBesFile(file *bes.File) *string {
	if file == nil {
		return nil
	}
	defaultDigestFunction := strings.ToLower(remoteexecution.DigestFunction_SHA256.String())
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		if len(stringArr) < 3 {
			return &defaultDigestFunction
		}
		digestFunction := stringArr[len(stringArr)-3]
		if _, ok := remoteexecution.DigestFunction_Value_value[strings.ToUpper(digestFunction)]; ok {
			return &digestFunction
		}
	}
	return &defaultDigestFunction
}
