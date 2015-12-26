# GOTS - Go One-Time-Secret 

## Description

This is a Go client for the onetimesecret.com service.  Please note you can also checkout the opensource code base for the service and self-host it.


## Dependencies

- Go 1.5+ (only tested with this version)
- github.com/codegangsta/cli
- github.com/franela/goreq
- An account at https://onetimesecret.com (if not self-hosting)

## Building

- git clone https://github.com/hartfordfive/gots.git
- cd gots && go build


## Usage & Parameters

You must have a ".gots" file in your home directory. Just use the .gots_sample file and replace values with the appropriate settings

The available commands are as follows:

Share a secret:
```bash
gots share [secret] [passphrase] [ttl] [recipient_email]
```

Generate a random a secret:
```bash
gots generate [passphrase] [ttl] [metadata_ttl] [secret_ttl] [recipient_email]
```

Get a current secret:
```bash
gots get [secret_key] [passphrase]
```

Get meta data for a given secret:
```bash
gots getmeta [metadata_key]
```

Get recent meta data:
```bash
gots recentmeta
```

View API status:
```bash
gots status
```

## Bugs & Feature Requests

Please open an issue for any bugs or feature requests.


## Author

Alain Lefebvre  (hartfordfive@gmail.com)


## License

Covered under the MIT License
