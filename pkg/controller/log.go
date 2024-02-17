/*
Copyright 2024 Jason Kölker.

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

package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Info logs a message at the info level to the logger in the context.
func Info(ctx context.Context, msg string, keysAndValues ...any) {
	log.FromContext(ctx).WithCallDepth(1).Info(msg, keysAndValues...)
}

// Error logs a message at the error level to the logger in the context.
func Error(ctx context.Context, err error, msg string, keysAndValues ...any) {
	log.FromContext(ctx).WithCallDepth(1).Error(err, msg, keysAndValues...)
}

// WrapError wraps an error with a message and logs it at the error level to
// the logger in the context.
func WrapError(ctx context.Context, err error, msg string, keysAndValues ...any) error {
	log.FromContext(ctx).WithCallDepth(1).Error(err, msg, keysAndValues...)

	return fmt.Errorf("%s: %w", msg, err)
}
