#!/bin/bash

go install code.cloudfoundry.org/go-pubsub/pubsub-gen

$GOPATH/bin/pubsub-gen \
  --struct-name=code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/end2end.X \
  --package=end2end_test \
  --traverser=StructTraverser \
  --output=$GOPATH/src/code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/end2end/generated_traverser_test.go \
  --pointer \
  --interfaces='{"message":["M1","M2","M3"]}' \
  --include-pkg-name=true \
  --imports=code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/end2end \
  --slices='{"X.RepeatedY":"I","RepeatedEmpty":""}'

gofmt -s -w .
