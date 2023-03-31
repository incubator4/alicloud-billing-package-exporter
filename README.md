AliCloud Billing Package Exporter
------
This repo is used to trans AliCloud billing package data to Prometheus metric-like data.


Before running it, you should create a pair of `access-key` and `scret-key` with policy
```json
{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "bssapi:QueryResourcePackageInstances"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
```
I suggest to use RAM account with minimal permission instead of directly create it.

### ENV

- ACCESS_KEY
- SECRET_KEY
- REGION (default: us-east-1)