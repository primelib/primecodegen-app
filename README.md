# PrimeLib App

Application to automatically update openapi specs and the generated client libraries.

## Usage

```bash
primelib-app run
```

## Configuration

Projects are configured using a `primelib.yaml` file in the root of the repository.

**Java**

```yaml
modules:
  - name: osv-rest # module name
    dir: . # relative to repository root
    spec_file: openapi.json # local spec file
    spec_url: https://osv.dev/docs/osv_service_v1.swagger.json # update spec from url
    spec_script: | # patch openapi spec before generation
      jq '.host = "api.osv.dev"' "$1" | sponge "$1" # set api host
      jq 'walk(if type == "object" and has("operationId") then .operationId |= sub("^OSV_"; "") else . end)' "$1" | sponge "$1" # remove prefix from operationId
    config: # openapi generator config
      generatorName: prime-client-java-feign
      invokerPackage: io.github.primelib.osv4j
      apiPackage: io.github.primelib.osv4j.api
      modelPackage: io.github.primelib.osv4j.model
      enablePostProcessFile: true
      additionalProperties:
        projectArtifactGroupId: io.github.primelib
        projectArtifactId: osv4j
```

## Platform Configuration

You are *required* to have the environment variables for one platform set.

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

### GitLab User

Create a GitLab user and generate a personal access token with the following permissions:

- api

| Environment Variable  | Description                |
|-----------------------|----------------------------|
| `GITLAB_SERVER`       | The GitLab server URL.     |
| `GITLAB_ACCESS_TOKEN` | The personal access token. |

## License

Released under the [MIT license](./LICENSE).
