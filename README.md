# Aptly Utils

Package provide some aptly utils

## Install

`go get github.com/bborbe/aptly_utils/bin/aptly_clean_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_copy_package`

`go get github.com/bborbe/aptly_utils/bin/aptly_create_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_delete_package`

`go get github.com/bborbe/aptly_utils/bin/aptly_delete_repo`

`go get github.com/bborbe/aptly_utils/bin/aptly_package_lister`

`go get github.com/bborbe/aptly_utils/bin/aptly_package_versions`

`go get github.com/bborbe/aptly_utils/bin/aptly_package_latest_version`

`go get github.com/bborbe/aptly_utils/bin/aptly_upload`

## List repositories

```
aptly_repo_lister \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password
```

## Create repository

```
aptly_create_repo \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable \
-architecture=amd64,all
```

## Delete repository

```
aptly_delete_repo \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable
```

## Clean repository

```
aptly_clean_repo \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable
```

## Upload Debian package

```
aptly_upload \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-file=booking_1.0.1-b47.deb \
-repo=unstable
```

## List packages

```
aptly_package_lister \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable
```

## Delete package

```
aptly_delete_package \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable \
-name=booking \
-version=1.0.1-b47
```

## Copy package from source to target repo

### Copy package with version

```
aptly_copy_package \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-source=unstable \
-target=stable \
-name=booking \
-version=1.0.1-b47
```

### Copy latest version

```
aptly_copy_package \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-source=unstable \
-target=stable \
-version=latest \
-name=booking 
```

### Copy latest version of each package

```
aptly_copy_package \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-source=unstable \
-target=stable \
-name=all \
-version=latest
```

## Latest Version of Package

```
aptly_package_latest_version \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable \
-name=booking
```

## Versions of Package

```
aptly_package_versions \
-logtostderr \
-v=2 \
-url=https://aptly.benjamin-borbe.de \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-repo=unstable \
-name=booking
```
