package qingcloud

import "time"

const TIMEOUT_CREATE_CLUSTER = 600 * time.Second
const TIMEOUT_START_CLUSTER = 600 * time.Second
const TIMEOUT_STOP_CLUSTER = 600 * time.Second
const TIMEOUT_DELETE_CLUSTER = 600 * time.Second
const TIMEOUT_RECOVER_CLUSTER = 600 * time.Second
const TIMEOUT_CEASE_CLUSTER = 600 * time.Second

const WAIT_INTERVAL = 20 * time.Second

const STATUS_ACTIVE = "active"
