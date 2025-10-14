# cert-manager-webhook-joker

[![Go Report Card](https://goreportcard.com/badge/github.com/nauski/cert-manager-webhook-joker)](https://goreportcard.com/report/github.com/nauski/cert-manager-webhook-joker)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A cert-manager webhook for Joker.com DNS-01 ACME challenges that enables automatic wildcard certificate issuance and renewal.

## Overview

This webhook integrates [cert-manager](https://cert-manager.io/) with [Joker.com](https://joker.com/) DNS services to automatically manage DNS TXT records for ACME DNS-01 challenges. It enables secure, automated issuance and renewal of wildcard SSL/TLS certificates using Let's Encrypt or other ACME certificate authorities.

**Why this webhook?** The existing `4nxio/cert-manager-webhook-joker` contains nil pointer dereference bugs that prevent proper operation. This implementation provides a clean, reliable alternative with comprehensive error handling and testing.

## Features

- ✅ **DNS-01 Challenge Support**: Automated TXT record management for ACME challenges
- ✅ **Wildcard Certificates**: Full support for `*.domain.com` certificates
- ✅ **Secure Credentials**: Kubernetes-native secret management
- ✅ **Production Ready**: Comprehensive error handling and logging
- ✅ **Kubernetes Native**: Seamless cert-manager integration
- ✅ **Multi-Domain Support**: Handles complex domain hierarchies

## Prerequisites

Before installing this webhook, ensure you have:

- **Kubernetes cluster** (v1.20+)
- **cert-manager** installed (v1.12+) - [Installation Guide](https://cert-manager.io/docs/installation/)
- **Domain hosted with Joker.com** with API access enabled
- **Joker.com credentials** (username and password for DNS management)

### Verify cert-manager Installation

```bash
kubectl get pods -n cert-manager
# Should show cert-manager, cert-manager-cainjector, and cert-manager-webhook pods
```

## Installation

### Step 1: Create Credentials Secret

Create a Kubernetes secret with your Joker.com credentials:

```bash
kubectl create secret generic joker-credentials \
  --namespace=cert-manager \
  --from-literal=username="YOUR_JOKER_USERNAME" \
  --from-literal=password="YOUR_JOKER_PASSWORD"
```

<details>
<summary>Alternative: Using a YAML file</summary>

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: joker-credentials
  namespace: cert-manager
type: Opaque
data:
  username: <base64-encoded-username>
  password: <base64-encoded-password>
```

```bash
# Encode your credentials
echo -n "your-username" | base64
echo -n "your-password" | base64
```
</details>

### Step 2: Deploy the Webhook

```bash
# Clone the repository (if not already done)
git clone https://github.com/nauski/cert-manager-webhook-joker.git
cd cert-manager-webhook-joker

# Deploy all components
kubectl apply -f deploy/
```

### Step 3: Verify Deployment

```bash
# Check webhook pod status
kubectl get pods -n cert-manager -l app=cert-manager-webhook-joker

# Check webhook logs
kubectl logs -n cert-manager deployment/cert-manager-webhook-joker

# Verify API service registration
kubectl get apiservice v1alpha1.acme.joker.com
```

### Step 4: Create a ClusterIssuer

Create a ClusterIssuer for Let's Encrypt with the Joker.com webhook:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-joker-prod
spec:
  acme:
    # Use Let's Encrypt production server
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com  # Replace with your email
    privateKeySecretRef:
      name: letsencrypt-joker-prod
    solvers:
    - dns01:
        webhook:
          groupName: acme.joker.com
          solverName: joker
          config:
            username:
              secretKeyRef:
                name: joker-credentials
                key: username
            password:
              secretKeyRef:
                name: joker-credentials
                key: password
```

<details>
<summary>Staging Environment ClusterIssuer</summary>

For testing, use the Let's Encrypt staging environment:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-joker-staging
spec:
  acme:
    # Use Let's Encrypt staging server (higher rate limits)
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-joker-staging
    solvers:
    - dns01:
        webhook:
          groupName: acme.joker.com
          solverName: joker
          config:
            username:
              secretKeyRef:
                name: joker-credentials
                key: username
            password:
              secretKeyRef:
                name: joker-credentials
                key: password
```
</details>

### Step 5: Request a Certificate

Create a certificate resource:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard-example-com
  namespace: default  # Change to your target namespace
spec:
  secretName: wildcard-example-com-tls
  issuerRef:
    name: letsencrypt-joker-prod  # Use staging for testing
    kind: ClusterIssuer
  dnsNames:
  - "example.com"      # Replace with your domain
  - "*.example.com"    # Wildcard certificate
```

## Configuration

### Webhook Configuration Parameters

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `username` | SecretKeyRef | Joker.com account username | ✅ Yes |
| `password` | SecretKeyRef | Joker.com account password | ✅ Yes |

### Example Configuration with Multiple Domains

```yaml
solvers:
- dns01:
    webhook:
      groupName: acme.joker.com
      solverName: joker
      config:
        username:
          secretKeyRef:
            name: joker-credentials
            key: username
        password:
          secretKeyRef:
            name: joker-credentials
            key: password
  selector:
    dnsNames:
    - "*.example.com"
    - "*.subdomain.example.com"
```

## DNS Setup Requirements

### Domain Hosting
- Your domain **must** be hosted with Joker.com nameservers
- Verify with: `dig NS yourdomain.com` should show `ns.joker.com` and `ns2.joker.com`

### DNS Propagation
- TXT record changes typically propagate within 2-5 minutes
- cert-manager will retry failed challenges automatically
- Monitor with: `dig TXT _acme-challenge.yourdomain.com`

## Monitoring and Observability

### Certificate Status

```bash
# Check certificate status
kubectl get certificates -A

# Describe certificate for details
kubectl describe certificate wildcard-example-com

# Check certificate events
kubectl get events --field-selector involvedObject.name=wildcard-example-com
```

### Challenge Debugging

```bash
# List active challenges
kubectl get challenges -A

# Describe challenge for details
kubectl describe challenge <challenge-name>

# Check order status
kubectl get orders -A
```

## Development

### Local Development Setup

```bash
# Clone repository
git clone https://github.com/nauski/cert-manager-webhook-joker.git
cd cert-manager-webhook-joker

# Install dependencies
go mod download

# Run tests
make test

# Build binary
make build

# Build Docker image
make docker-build
```

### Project Structure

```
cert-manager-webhook-joker/
├── cmd/webhook/           # Application entry point
├── internal/
│   ├── joker/            # Joker.com API client
│   ├── solver/           # cert-manager webhook implementation
│   └── util/             # DNS parsing utilities
├── deploy/               # Kubernetes manifests
├── docs/                 # Documentation
└── Makefile             # Build automation
```

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/util -run TestParseChallengeDomain
```

## API Reference

### Joker.com DNS API

The webhook integrates with the Joker.com DNS management API:

- **Endpoint**: `https://svc.joker.com/nic/replace`
- **Method**: POST
- **Authentication**: Form-based (username/password)
- **Rate Limits**: Reasonable usage expected

### Domain Parsing Logic

The webhook automatically parses domains to determine the correct DNS zone and record label:

| Input Domain | DNS Zone | TXT Record Label |
|--------------|----------|------------------|
| `example.com` | `example.com` | `_acme-challenge` |
| `sub.example.com` | `example.com` | `_acme-challenge.sub` |
| `*.app.example.com` | `example.com` | `_acme-challenge.app` |

## Troubleshooting

### Common Issues

#### 1. Authentication Failures

**Symptoms**: `Authentication failed: invalid credentials`

**Solutions**:
- Verify Joker.com credentials are correct
- Check secret exists and is in correct namespace
- Ensure API access is enabled in Joker.com dashboard

```bash
# Verify secret
kubectl get secret joker-credentials -n cert-manager -o yaml

# Check webhook logs
kubectl logs -n cert-manager deployment/cert-manager-webhook-joker
```

#### 2. DNS Propagation Issues

**Symptoms**: Challenge validation timeouts

**Solutions**:
- Wait 5-10 minutes for DNS propagation
- Verify domain uses Joker.com nameservers
- Check for conflicting DNS records

```bash
# Check DNS propagation
dig TXT _acme-challenge.yourdomain.com @8.8.8.8

# Verify nameservers
dig NS yourdomain.com
```

#### 3. Domain Parsing Errors

**Symptoms**: `Failed to parse domain` errors

**Solutions**:
- Ensure domain format is correct (no spaces, valid TLD)
- Check domain is actually hosted with Joker.com
- Verify certificate covers intended domains

#### 4. Webhook Not Responding

**Symptoms**: `connection refused` or timeout errors

**Solutions**:

```bash
# Check webhook pod status
kubectl get pods -n cert-manager -l app=cert-manager-webhook-joker

# Check service endpoints
kubectl get endpoints -n cert-manager cert-manager-webhook-joker

# Verify API service
kubectl get apiservice v1alpha1.acme.joker.com
```

### Debug Mode

Enable verbose logging by setting log level:

```yaml
# In webhook deployment
args:
- --v=4  # Verbose logging
```

### Getting Help

If you encounter issues:

1. **Check logs** with `kubectl logs`
2. **Search existing issues** on GitHub
3. **Create detailed bug report** with:
   - Kubernetes version
   - cert-manager version
   - Webhook logs
   - Certificate/Challenge YAML

## Security Considerations

- **Credentials**: Store only in Kubernetes secrets, never in code
- **RBAC**: Webhook runs with minimal required permissions
- **Network**: Consider network policies for additional isolation
- **Updates**: Keep webhook updated for security patches

## Compatibility

| Component | Version | Support |
|-----------|---------|---------|
| Kubernetes | 1.20+ | ✅ Tested |
| cert-manager | 1.12+ | ✅ Tested |
| cert-manager | 1.13+ | ✅ Recommended |
| Go | 1.21+ | ✅ Required |

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Reporting Issues

Please use GitHub Issues for:
- Bug reports
- Feature requests
- Documentation improvements

Include:
- Environment details (K8s, cert-manager versions)
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history and breaking changes.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [cert-manager](https://cert-manager.io/) team for the excellent ACME automation framework
- [Joker.com](https://joker.com/) for providing reliable DNS services and API access
- Community contributors and issue reporters

---

**Questions?** Check our [FAQ](docs/FAQ.md) or open an issue on GitHub.