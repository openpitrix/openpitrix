// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

const ClusterSchema = `
{
  "additionalProperties": false,
  "$schema": "http://json-schema.org/draft-04/schema#",
  "required": [
    "name",
    "description",
    "subnet",
    "nodes"
  ],
  "type": "object",
  "properties": {
    "description": {
      "type": "string"
    },
    "advanced_actions": {
      "uniqueItems": true,
      "items": {
        "enum": [
          "change_subnet",
          "scale_horizontal",
          "add_nodes",
          "associate_eip"
        ],
        "type": "string"
      },
      "type": "array"
    },
    "name": {
      "type": "string"
    },
    "subnet": {
      "type": "string"
    },
    "env": {
      "patternProperties": {
        "^.*$": {}
      },
      "type": "object"
    },
    "nodes": {
      "uniqueItems": true,
      "items": {
        "additionalProperties": false,
        "required": [
          "container",
          "count"
        ],
        "type": "object",
        "properties": {
          "vertical_scaling_policy": {
            "enum": [
              "sequential",
              "parallel"
            ],
            "type": "string"
          },
          "replica": {
            "minimum": 0,
            "type": "integer",
            "maximum": 20
          },
          "container": {
            "additionalProperties": false,
            "required": [
              "image",
              "type"
            ],
            "type": "object",
            "properties": {
              "image": {
                "type": "string",
                "maxLength": 2048
              },
              "type": {
                "enum": [
                  "kvm",
                  "docker",
                  "lxc",
                  "bm"
                ],
                "type": "string"
              },
              "zone": {
                "type": "string"
              }
            }
          },
          "role": {
            "type": "string"
          },
          "env": {
            "patternProperties": {
              "^.*$": {}
            },
            "type": "object"
          },
          "memory": {
            "type": "integer"
          },
          "gpu": {
            "type": "integer"
          },
          "server_id_upper_bound": {
            "minimum": 1,
            "type": "integer"
          },
          "volume": {
            "additionalProperties": false,
            "required": [],
            "type": "object",
            "properties": {
              "size": {
                "minimum": 0,
                "type": "integer"
              },
              "mount_point": {
                "anyOf": [
                  {
                    "uniqueItems": true,
                    "items": {
                      "pattern": "^/.*?$",
                      "type": "string"
                    },
                    "type": "array"
                  },
                  {
                    "uniqueItems": true,
                    "items": {
                      "pattern": "^[d-z]:$",
                      "type": "string"
                    },
                    "type": "array"
                  },
                  {
                    "pattern": "^/.*?$",
                    "type": "string"
                  },
                  {
                    "pattern": "^[d-z]:$",
                    "type": "string"
                  }
                ]
              },
              "filesystem": {
                "enum": [
                  "ext4",
                  "xfs",
                  "ntfs"
                ],
                "type": "string"
              },
              "mount_options": {
                "type": "string",
                "maxLength": 100
              }
            }
          },
          "services": {
            "additionalProperties": false,
            "maxProperties": 21,
            "patternProperties": {
              "^restart$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              },
              "^start$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              },
              "^stop$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              },
              "^backup$": {
                "additionalProperties": false,
                "type": "object",
                "properties": {
                  "service_params": {
                    "patternProperties": {
                      "^.*$": {}
                    },
                    "type": "object"
                  },
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              },
              "^upgrade$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "post_start_service": {
                    "type": "boolean"
                  },
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  },
                  "service_params": {
                    "patternProperties": {
                      "^.*$": {}
                    },
                    "type": "object"
                  },
                  "order": {
                    "type": "integer"
                  }
                }
              },
              "^scale_out$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  },
                  "pre_check": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  }
                }
              },
              "^init$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  },
                  "post_start_service": {
                    "type": "boolean"
                  }
                }
              },
              "^delete_snapshot$": {
                "additionalProperties": false,
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              },
              "^scale_in$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  },
                  "pre_check": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  }
                }
              },
              "^destroy$": {
                "additionalProperties": false,
                "required": [
                  "cmd"
                ],
                "type": "object",
                "properties": {
                  "allow_force": {
                    "type": "boolean"
                  },
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  },
                  "post_stop_service": {
                    "type": "boolean"
                  },
                  "order": {
                    "type": "integer"
                  }
                }
              },
              "^restore$": {
                "additionalProperties": false,
                "type": "object",
                "properties": {
                  "service_params": {
                    "patternProperties": {
                      "^.*$": {}
                    },
                    "type": "object"
                  },
                  "cmd": {
                    "pattern": "^.*[^\\s]+.*$",
                    "type": "string",
                    "maxLength": 1000
                  },
                  "nodes_to_execute_on": {
                    "type": "integer"
                  },
                  "order": {
                    "type": "integer"
                  },
                  "timeout": {
                    "type": "integer",
                    "maximum": 86400
                  }
                }
              }
            },
            "type": "object"
          },
          "count": {
            "minimum": 0,
            "type": "integer",
            "maximum": 200
          },
          "instance_type": {
            "type": "string"
          },
          "cpu": {
            "type": "integer"
          },
          "loadbalancer": {
            "type": "array"
          }
        }
      },
      "type": "array",
      "maxItems": 100
    },
    "incremental_backup_supported": {
      "type": "boolean"
    }
  }
}
`
