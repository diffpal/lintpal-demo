# LintPal demo

This repository demonstrates [LintPal](https://github.com/diffpal/lintpal) on a small Go HTTP service. The default branch is a passing baseline. The canonical pull request intentionally introduces authorization and input-validation violations so LintPal can publish inline findings and a deterministic gate result.

## Run the service

```bash
export DEMO_API_KEY='local-demo-key'
go run ./cmd/server
```

Create an order:

```bash
curl -i http://localhost:8080/orders \
  -H 'Authorization: Bearer local-demo-key' \
  -H 'Content-Type: application/json' \
  --data '{"customer":"Ada","item":"keyboard","quantity":2}'
```

## Repository rules

The committed `.lintpal/rules/` files come from the versioned [LintPal rule catalog](https://github.com/diffpal/lintpal-rules):

```bash
lintpal rule import github:diffpal/lintpal-rules//general@v1.0.0
lintpal rule import github:diffpal/lintpal-rules//go@v1.0.0
lintpal rule validate
```

This demo lowers the authorization and input-validation thresholds to `0.50` so the intentionally broken pull request remains a reliable feedback example. Production repositories should choose thresholds for their own gate policy.

## Enable pull-request feedback

1. Open **Settings → Secrets and variables → Actions**.
2. Add a repository secret named `TYPESAFE_API_KEY` containing the TypeSafe/Jev provider key.
3. Re-run the LintPal workflow on the canonical pull request.

The workflow uses GitHub's generated `GITHUB_TOKEN` only to publish the check result and inline comments. Provider credentials are environment variables, never Action inputs. A job-level same-repository guard prevents fork pull requests from reaching `TYPESAFE_API_KEY`.

LintPal classifies changed lines against fixed rules through Jev. It does not generate a narrative code review. The summary is deterministic: `No blocking findings`, `1 blocking finding`, or `N blocking findings`.

## Ecosystem

- [LintPal CLI](https://github.com/diffpal/lintpal)
- [Reusable GitHub Action](https://github.com/diffpal/lintpal-action)
- [Versioned rule catalog](https://github.com/diffpal/lintpal-rules)

MIT licensed. See [LICENSE](LICENSE).
