./gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1  -t v0.5.0 -m non-prod
./gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod
./gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod --dry-run