# terraform-provider-openrouter

Terraform provider for [OpenRouter](https://openrouter.ai) — unified API access to 300+ AI models.

## Quick Reference

| | |
|---|---|
| **Registry** | `echohello-dev/openrouter` |
| **SDK** | `terraform-plugin-sdk/v2` |
| **Go** | 1.25 |
| **Tools** | mise, goreleaser, golangci-lint, tfplugindocs |
| **Release** | release-please + goreleaser + GPG signing |

## Architecture

```
main.go                              # plugin entrypoint
openrouter/
  provider.go                        # provider schema + resource/data source registration
  client.go                          # HTTP client, request/response types
  resource_chat_completion.go        # the only resource
  data_source_*.go                   # data sources
  provider_test.go                   # schema validation + JSON parsing tests
.github/workflows/
  ci.yml                             # build, test, lint, docs-check via mise
  release.yml                        # release-please + goreleaser
examples/
  resources/<name>/main.tf           # Registry per-resource examples
  data-sources/<name>/main.tf        # Registry per-data-source examples
scripts/
  docs-check.sh                      # assert docs are up to date
  gpg-verify.sh                      # assert GPG key fingerprint
mise.toml                            # tool versions + tasks
.goreleaser.yaml                     # build matrix, archives, checksums, GPG signing
```

## Tooling

All tools managed by [mise](https://mise.jdx.dev). Install once, then:

```bash
mise install          # install all tools from mise.toml
mise run ci           # build + vet + test + lint
mise run validate     # fmt + vet + lint + test
mise run doc          # regenerate tfplugindocs
mise run doc-check    # assert docs are up to date
mise run release      # goreleaser release --clean
```

## Adding a Data Source

1. Add request/response types to `openrouter/client.go`
2. Create `openrouter/data_source_<name>.go` with `ReadContext` func
3. Register in `openrouter/provider.go` `DataSourcesMap`
4. Add schema test in `openrouter/provider_test.go`
5. Add `examples/data-sources/<name>/main.tf`
6. Run `mise run doc` to regenerate docs
7. Run `mise run ci` before committing

## Adding a Resource

Same pattern as data source, but use `CreateContext`/`ReadContext`/`DeleteContext`.

Current resource (`openrouter_chat_completion`) is create-only:
- `ReadContext` is a no-op (OpenRouter has no GET endpoint for past completions)
- `DeleteContext` clears the ID
- `Update` is not implemented; all mutable fields use `ForceNew: true`

## Schema Conventions

- Required fields: `Required: true` + `ValidateFunc: validation.NoZeroValues`
- Optional numeric params that default to API value: no `Default`, use `GetRawConfig` + `IsNull()` check so explicit `0` is sent
- Optional params with validation: `ValidateFunc` for ranges/enums
- Computed fields: `Computed: true`, never `Required` or `Optional`
- Sensitive: `Sensitive: true` for API keys, tokens, credentials

## Testing

- Schema validation tests: `TestProvider`, `Test*Schema` functions
- JSON parsing tests: unmarshal real API responses into Go structs
- No integration tests (would need live API key)

## Release Process

1. Push conventional commits to `main`
2. release-please opens a "chore(main): release X.Y.Z" PR
3. Merge the release PR
4. release-please creates tag + GitHub Release
5. goreleaser builds binaries, generates SHASUMS, signs with GPG
6. Terraform Registry ingests from GitHub Release automatically

Required repo secrets for releases:
- `GPG_PRIVATE_KEY` — ASCII-armored secret key
- `GPG_PASSPHRASE` — key passphrase (if any)

## Commit Convention

Conventional Commits, enforced by release-please changelog grouping:

| Type | Section |
|---|---|
| `feat:` | Features |
| `fix:` | Bug Fixes |
| `perf:` | Performance |
| `refactor:` | Refactoring |
| `docs:` | Documentation |
| `chore:` / `ci:` | Miscellaneous |

Always include `Co-authored-by: opencode-agent <noreply@opencode.ai>` when an AI agent assists.
