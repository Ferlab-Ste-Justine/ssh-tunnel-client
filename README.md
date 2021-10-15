# About

This project provides a convenience for users to open an ssh-tunnel to a bastion sitting in front of a private services.

It allows for some configurability in ssh connection parameters and the bindings to the services offered on the bastion's network.

# Usage

## Requisite Files

The client requires an **auth_secret** file, either embedded in the binary at build time or present in the running directory. This file should contain either the private ssh key or login password of the user to be sshed as.

The client also requires a **tunnel_config.json** file, either embbeded in the binary at build time or present in the running directory. The file should have the following format:

```
{
    "host_md5_fingerprint": "<md5 fingerprint of the bastions public key>",
    "host_url": "<bastion ip:bastion ssh port>",
    "host_user": "<bastion user to ssh as>",
    "auth_method": "key"|"password",
    "bindings": [
        {
            "local": "<local ip to bind service 1 on>:<local port to bind service 1 on>",
            "remote": "<Remote ip on the bastion's network where service 1 is reachable>:<Remote port that service 1 can be reached at>"
        },
        {
            "local": "<local ip to bind service 2 on>:<local port to bind service 2 on>",
            "remote": "<Remote ip on the bastion's network where service 2 is reachable>:<Remote port that service 2 can be reached at>"
        },
        ...
    ]
}
```

## Embedding Files

The **auth_secret** and/or the **tunnel_config.json** files can be included in this directory prior to building the binary in which case they'll be embedded in the binary.

If both files are embedded in the binary, end-users will not have to define them (though the binary should be considered a secret in this case).

Those files can also be put in the running directory of the program after the binary is built.

If the files were embedded in the binary **AND** are also provided in the running directory of the program, the files in the running directory of the program will take precedence.