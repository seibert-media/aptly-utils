#!/bin/sh

SOURCEDIRECTORY="github.com/bborbe/aptly_utils"
INSTALLS="github.com/bborbe/aptly_utils/bin/aptly_package_lister github.com/bborbe/aptly_utils/bin/aptly_package_version github.com/bborbe/aptly_utils/bin/aptly_copy_package github.com/bborbe/aptly_utils/bin/aptly_create_repo github.com/bborbe/aptly_utils/bin/aptly_delete_package github.com/bborbe/aptly_utils/bin/aptly_delete_repo github.com/bborbe/aptly_utils/bin/aptly_upload github.com/bborbe/aptly_utils/bin/aptly_clean_repo"
VERSION="1.0.1-b${BUILD_NUMBER}"
NAME="aptly-utils"

################################################################################

export GOROOT=/opt/go
export PATH=/opt/go2xunit/bin/:/opt/utils/bin/:/opt/aptly_utils/bin/:/opt/aptly/bin/:/opt/debian_utils/bin/:/opt/debian/bin/:$GOROOT/bin:$PATH
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports
DEB="${NAME}_${VERSION}.deb"
rm -rf $REPORT_DIR ${WORKSPACE}/*.deb ${WORKSPACE}/pkg
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find $SOURCEDIRECTORY -name "*_test.go" | dirof | unique`
FAILED=false
for PACKAGE in $PACKAGES
do
  XML=$REPORT_DIR/`pkg2xmlname $PACKAGE`
  OUT=$XML.out
  go test -i $PACKAGE
  go test -v $PACKAGE | tee $OUT
  cat $OUT
  go2xunit -fail=true -input $OUT -output $XML
  rc=$?
  if [ $rc != 0 ]
  then
    echo "Tests failed for package $PACKAGE"
    FAILED=true
  fi
done

if $FAILED
then
  echo "Tests failed => skip install"
  exit 1
else
  echo "Tests success"
fi

echo "Tests completed, install to $GOPATH"

go install $INSTALLS

echo "Install completed, create debian package"

create_debian_package \
-loglevel=DEBUG \
-version=$VERSION \
-config=src/$SOURCEDIRECTORY/create_debian_package_config.json || exit 1

echo "Create debian package completed, upload"

aptly_upload \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-file=$DEB \
-repo=unstable || exit 1

echo "Upload completed"
