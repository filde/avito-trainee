package database

import "avito-trainee/external/httpserver"

var _ httpserver.StorageItf = &Database{}
