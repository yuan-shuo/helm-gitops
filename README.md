[中文版](https://github.com/yuan-shuo/helm-gitops/blob/main/README.zh-CN.md)

# Helm-GitOps

A Helm extension that provides GitOps-related auxiliary functions for Helm

## Usage

### Git-related Features

Compared to the conventional approach of manually creating Helm charts, pasting `.gitignore` files internally, and performing a series of tedious operations such as creating branches, committing, and managing version numbers, this extension provides a more comfortable simplified solution:

```bash
# 1. Create: Add GitOps skeleton on top of the original Helm chart
helm gitops create test
helm gitops create test --actions        # Also generate .github/workflows/<>.yaml


# 2. Checkout: Automatically switch to development branch (optional sync with main branch)
helm gitops checkout feature/foo
helm gitops checkout feature/foo -s       # Pull origin/main first, then create branch


# 3. Commit: add + commit + optional push & automatic PR
helm gitops commit -m "fix: foo"                    # Local commit
helm gitops commit -m "feat: xxx" --push            # Commit and push
helm gitops commit -m "ci: update" --pr --push      # Commit + push + automatic PR (includes [create-pr] commit marker)


# 4. Local check: helm lint + helm unittest
helm gitops lint


# 5. Push code: lint → push (protected branch interception)
helm gitops push                                      # Push to origin/current branch


# 6. Version management
helm gitops version # Query current version only
helm gitops version -m pr -l patch # Traditional PR mode (create release branch first → submit PR → CI auto tag)
helm gitops version -m main -l patch # Quick main branch mode (direct commit + tag + push simultaneously)
# --mode=main|pr
# --level=patch|minor|major
```

### Environment Repository Features

To be developed

### ArgoCD Features

To be developed

## Installation

### Using helm plugin install

```bash
helm plugin install https://github.com/yuan-shuo/helm-gitops
```

### Using Binary Files

- Go to: [Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases) to download the binary file for your operating system

- Place the extracted `gitops` binary file in the `$HELM_PLUGIN_DIR/bin/` directory

- Grant execute permissions to the gitops binary file

```bash
chmod +x $HELM_PLUGIN_DIR/bin/gitops
```

## Requirements

* git (version >= 2.23)
* helm
* helm-unitest (optional, can be installed by running **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`**)
