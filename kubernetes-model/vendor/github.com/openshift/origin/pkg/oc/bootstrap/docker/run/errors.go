/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package run

import (
	"bytes"
	"fmt"

	"github.com/docker/engine-api/types/container"
	"github.com/openshift/origin/pkg/oc/bootstrap/docker/errors"
)

type runError struct {
	error
	stdOut string
	stdErr string
	config *container.Config
}

func newRunError(rc int, cause error, stdOut, stdErr string, config *container.Config) error {
	return &runError{
		error:  errors.NewError("Docker run error rc=%d", rc).WithCause(cause),
		stdOut: stdOut,
		stdErr: stdErr,
		config: config,
	}
}

func (e *runError) Details() string {
	out := &bytes.Buffer{}
	fmt.Fprintf(out, "Image: %s\n", e.config.Image)
	fmt.Fprintf(out, "Entrypoint: %v\n", e.config.Entrypoint)
	fmt.Fprintf(out, "Command: %v\n", e.config.Cmd)
	if len(e.stdOut) > 0 {
		errors.PrintLog(out, "Output", e.stdOut)
	}
	if len(e.stdErr) > 0 {
		errors.PrintLog(out, "Error Output", e.stdErr)
	}
	return out.String()
}
