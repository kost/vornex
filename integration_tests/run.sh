#!/bin/bash

echo "::group::Build vornex"
rm integration-test vornex 2>/dev/null
cd ../cmd/vornex
go build
mv vornex ../../integration_tests/vornex
echo "::endgroup::"

echo "::group::Build vornex integration-test"
cd ../integration-test
go build
mv integration-test ../../integration_tests/integration-test 
cd ../../integration_tests
echo "::endgroup::"

sudo ./integration-test
if [ $? -eq 0 ]
then
  exit 0
else
  exit 1
fi
