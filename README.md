# migrator

A tool to migrate a [BOSH vars-store](https://bosh.io/docs/cli-int.html#vars-store) to [CredHub](https://github.com/cloudfoundry-incubator/credhub).

### Installation

```sh
go get github.com/ishustava/migrator
```

### Usage

Migrator has only one command - `migrate`.

```
Usage:
  migrator [OPTIONS] migrate [migrate-OPTIONS]

Help Options:
  -h, --help                 Show this help message

[migrate command options]
      -v, --vars-store=      Path to vars store file.
      -u, --credhub-server=  URL of the CredHub server, e.g. https://example.com:8844 [$CREDHUB_SERVER]
      -c, --credhub-client=  UAA client for the CredHub Server [$CREDHUB_CLIENT]
      -s, --credhub-secret=  UAA client secret for the CredHub Server [$CREDHUB_SECRET]
      -e, --director-name=   Name of the BOSH director
      -d, --deployment-name= Name of the BOSH deployment with which vars store is used
          --ca-cert=         Trusted CA for API and UAA TLS connections. Multiple flags may be provided.
                             [$CREDHUB_CA_CERT]
```
