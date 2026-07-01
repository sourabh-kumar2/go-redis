package config

const DEFAULT_PORT = 7379
const DEFAULT_HOST = "0.0.0.0"

var Host string
var Port int
var KeysLimit int = 1 << 20
