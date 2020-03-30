/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package buildcontext

import (
	"errors"
	"strings"

	"github.com/GoogleContainerTools/kaniko/pkg/constants"
	"github.com/GoogleContainerTools/kaniko/pkg/util"
)

const (
	TarBuildContextPrefix = "tar://"
)

// BuildContext unifies calls to download and unpack the build context.
type BuildContext interface {
	// Unpacks a build context and returns the directory where it resides
	UnpackTarFromBuildContext() (string, error)
}

// GetBuildContext parses srcContext for the prefix and returns related buildcontext
// parser
func GetBuildContext(srcContext string) (BuildContext, error) {
	split := strings.SplitAfter(srcContext, "://")
	prefix := split[0]
	context := split[1]

	switch prefix {
	case constants.GCSBuildContextPrefix:
		return &GCS{context: context}, nil
	case constants.S3BuildContextPrefix:
		return &S3{context: context}, nil
	case constants.LocalDirBuildContextPrefix:
		return &Dir{context: context}, nil
	case constants.GitBuildContextPrefix:
		return &Git{context: context}, nil
	case constants.HTTPSBuildContextPrefix:
		if util.ValidAzureBlobStorageHost(srcContext) {
			return &AzureBlob{context: srcContext}, nil
		}
		return nil, errors.New("url provided for https context is not in a supported format, please use the https url for Azure Blob Storage")
	case TarBuildContextPrefix:
		return &Tar{context: context}, nil
	}
	return nil, errors.New("unknown build context prefix provided, please use one of the following: gs://, dir://, tar://, s3://, git://, https://")
}
