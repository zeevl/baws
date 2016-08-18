# BAWS (Better AWS CLI)

Utility to provide some functionality that's missing from the AWS CLI.

Currently supports the following commands:

### s3 ls

Improved version of s3 ls that supports `include` filemask, and returns
*all* results, not just the first 1000.

`
baws s3 ls s3://bucketname/prefix --include "*/*/*.png"
`
