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

### Purpose

Modern BOSH deployment manifests support `variables` declarations which allow for automatic generation of credentials, including hierarchies of certificates, so that the operator doesn't have to bear that toil.  If your BOSH Director is not integrated with CredHub, the only way to leverage this capability is to have the BOSH CLI generate a local file (called a "vars-store") full of credentials, which you then need to pass to your `bosh deploy` command.  A better, more secure way to do this is to connect your Director to CredHub so that it can securely generate and store all those credentials.  A problem arises if you have already deployed something with locally generated credentials, and want to migrate to a CredHub based deployment.  This is where `migrator` comes in!

`migrator` takes a "vars-store" file, and moves all the credentials therein into CredHub such that you could then throw away that vars-store and have everything managed by CredHub.  There are some finicky details that you need to get right when moving credentials from a vars-store into CredHub so that your deployments continue to work, which `migrator` takes care of.  This includes:

1. ensuring that credentials are moved to the correctly-namespaced path within the CredHub key-value store (this is what `migrator` uses the `--director-name` and `--deployment-name` flags for)
1. determining and setting the `ca_name` for every certificate so that certificates can be easily rotated at any point in the future using the [`credhub bulk-regenerate` command](https://github.com/cloudfoundry-incubator/credhub-cli/releases/tag/1.7.0)
