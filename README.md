# PrimeLib App

Application to automatically update openapi specs and the generated client libraries.

## Configuration

### GitHub App

Create a GitHub App and configure it with the following permissions:

- Repository contents: Read & write
- Pull requests: Read & write
- Commit statuses: Read & write
- Checks: Read & write
- Metadata: Read-only

Create a private key and store it in a file.

| Environment Variable          | Description                       |
|-------------------------------|-----------------------------------|
| `GITHUB_APP_ID`               | The ID of the GitHub App.         |
| `GITHUB_APP_PRIVATE_KEY_FILE` | The path to the private key file. |

## Usage

You are *required* to have the environment variables for one platform set.

```bash
primelib-app run
```

## Configuration

Projects are configured using a `primelib.yaml` file in the root of the repository.

```yaml
# TODO: add configuration example
```

## License

Released under the [MIT license](./LICENSE).
