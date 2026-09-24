```markdown
# Tirle  Deployment  Task :xxxxx

## Symbol | สัญลักษณ์

```html
✅ Success  ❌ Failed   ✅ Rolled back 
```

--- 
### Deployment Checklist

```markdown

## Deployment Checklist
### Pre-Deployment (Staging)

#### Code Readiness
- [ ] All tests passing (`go test ./...`)
- [ ] Race detector passed (`go test -race ./...`)
- [ ] Linter passed (`golangci-lint run ./...`)
- [ ] Build successful (`go build ./...`)
- [ ] All PRs merged and approved

#### Configuration
- [ ] Environment variables verified
- [ ] Configuration files updated for staging
- [ ] Feature flags configured
- [ ] Third-party service credentials verified

#### Database
- [ ] Migration scripts reviewed
- [ ] Migrations tested in staging environment
- [ ] Rollback plan documented
- [ ] Backup created before migration

#### Infrastructure
- [ ] Container images built and tagged
- [ ] Kubernetes/ deployment files updated
- [ ] Resource limits configured
- [ ] Health check endpoints configured
- [ ] Monitoring and alerting configured

#### Security
- [ ] Security scan passed
- [ ] No secrets in code or config
- [ ] TLS certificates valid

---
  

```html
**Deployment Status:** [ ] Success | [ ] Failed | [ ] Rolled back

**Deployed by:** _________________________

**Date:** _________________________

**Version:** _________________________

```
 