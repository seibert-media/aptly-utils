# Aptly Utils

Package provide some aptly utils

## Install

`go get github.com/bborbe/aptly_utils/bin/aptly_clean_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_copy_package`

`go get github.com/bborbe/aptly_utils/bin/aptly_create_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_delete_package`

`go get github.com/bborbe/aptly_utils/bin/aptly_delete_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_package_lister`

`go get github.com/bborbe/aptly_utils/bin/aptly_package_version`

`go get github.com/bborbe/aptly_utils/bin/aptly_upload`

## Create repository

```
aptly_create_repo \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable \
-architecture=amd64,all
```

## Delete repository

```
aptly_delete_repo \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable
```

## Clean repository

```
aptly_clean_repo \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable
```

## Upload Debian package

```
aptly_upload \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-file=booking_1.0.1-b47.deb \
-repo=unstable
```

## List packages

```
aptly_package_lister \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable
```

## Delete package

```
aptly_delete_package \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable \
-name=booking \
-version=1.0.1-b47
```

## Copy package from source to target repo

### Copy package with version

```
aptly_copy_package \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-source=unstable \
-target=stable \
-name=booking \
-version=1.0.1-b47
```

### Copy latest version

```
aptly_copy_package \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-source=unstable \
-target=stable \
-version=latest \
-name=booking 
```

### Copy latest version of each package

```
aptly_copy_package \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-source=unstable \
-target=stable \
-name=all \
-version=latest
```

## Version of Package

```
aptly_package_version \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=/etc/aptly_api_password \
-repo=unstable \
-name=booking
```

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Aptly-Utils/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
