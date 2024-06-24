
# GGScan

Warning - this is an unofficial test project for learning purposes

<img alt="GGScan file demo with VHS" src="./demo.gif" width="800" />

# Installation

```
brew install ggscan
```

Add your GitGuardian API key as a ```GITGUARDIAN_API_KEY``` environment variable then check with the ```health```command
```
ggscan health
```

# Usage

### Health
Check if everything is ok before using the API
```
ggscan health
```

### quotas
Check your current quota
```
ggscan quotas
```
### scan
Select one file for scanning
```
ggscan scan
```
