// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2024 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	indexConfig "github.com/falcosecurity/falcoctl/pkg/index/config"
)

// Fetch fetches the raw index file from an S3 object.
func Fetch(ctx context.Context, conf *indexConfig.Entry) ([]byte, error) {
	o, err := s3ObjectFromURI(conf.URL)
	if err != nil {
		return nil, err
	}

	// Create a new AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		// handle error
		return nil, fmt.Errorf("unable to create AWS config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	// Get the object from S3
	res, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(o.Bucket),
		Key:    aws.String(o.Key),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get S3 object: %w", err)
	}
	defer res.Body.Close()

	// Read the object data
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading S3 object: %w", err)
	}

	return bytes, nil
}
