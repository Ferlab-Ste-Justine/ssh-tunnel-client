# About

This project provides a convenience for users to open an ssh-tunnel to a bastion sitting in front of a private kubernetes cluster.

# Assumptions

Currently, mostly due to the prototypal status of this project, the following assumptions are made:

- You want to open the k8 api and a tls ingress ports and your bastion load balances this traffic for you (by running some reverse proxy) on ports 6443 and 443 locally.
- Similar ports are available to be opened on your localhost
- You will login to the bastion as user **ubuntu**
- The following files are local to your binary:
  - host-md5-fingerprint: Contains your bastion's ssh fingerprint
  - authorized-ssh-private-key: Contains the ssh key that you will use to authenticate against the bastion 
  - tunnel-server-url: Contains the ssh url of the bastion in the following format ```<ip>:<port>```

Not as customizable as it should be for general availability, but it fits our needs right at this moment.

# Usage

Compile it, customize your input files and launch it. When you are done, hit **cltr + c**.