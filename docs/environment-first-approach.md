# Implementing "Environment First" Configuration

To prioritize Environment Variables over files or flags, we can introduce a `LoadEnv` function. This function uses Go's `reflect` package to inspect your configuration struct, look for matching environment variables, and inject the values directly.

### 1. Update the Config Struct (`config/config.go`)

First, we need to tell our code *which* environment variable maps to *which* field. We do this by adding an `env` tag to the struct fields.

```go
package config

// Config contains application settings.
type Config struct {
    // We add the `env` tag to specify the environment variable name
    General struct {
        ApplicationName string `env:"GOALERT_APP_NAME" info:"Name used in messaging."`
        PublicURL       string `env:"GOALERT_PUBLIC_URL" info:"Publicly routable URL."`
        DisableSMSLinks bool   `env:"GOALERT_DISABLE_SMS_LINKS"`
    }
    
    Database struct {
        URL string `env:"GOALERT_DB_URL"`
    }
    
    // ... other sections
}

```

# GoAlert Environment Variable Configuration

Below is the comprehensive list of environment variables used to configure GoAlert, extracted from the documentation.

### ⚠️ Critical / Required Configuration
These variables are essential for the application to start and function correctly in a production environment.

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_DB_URL`** | [cite_start]**Required.** Connection string for the PostgreSQL database[cite: 1046].<br>Example: `postgres://user:pass@localhost/goalert` |
| **`GOALERT_PUBLIC_URL`** | [cite_start]**Required.** The externally routable URL to the application[cite: 1051].<br>Used for link generation, validation, and auth callbacks. |
| **`GOALERT_DATA_ENCRYPTION_KEY`** | [cite_start]**Recommended/Required.** Used to encrypt sensitive data (like API keys) in the database[cite: 1053]. |

---

### 🌐 Network & Connectivity
Variables controlling how GoAlert listens for traffic and handles security headers.

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_LISTEN`** | Address and port to listen on. [cite_start]Default is `localhost:8081`[cite: 1087]. |
| **`GOALERT_LISTEN_TLS`** | Address and port for HTTPS. [cite_start]Requires setting TLS cert/key variables below[cite: 1095]. |
| **`GOALERT_TLS_CERT_FILE`** | [cite_start]Path to PEM-encoded certificate file for HTTPS[cite: 686]. |
| **`GOALERT_TLS_KEY_FILE`** | [cite_start]Path to PEM-encoded private key file for HTTPS[cite: 688]. |
| **`GOALERT_TLS_CERT_DATA`** | [cite_start]String containing PEM-encoded certificate data (alternative to file)[cite: 685]. |
| **`GOALERT_TLS_KEY_DATA`** | [cite_start]String containing PEM-encoded private key data (alternative to file)[cite: 687]. |
| **`GOALERT_DISABLE_HTTPS_REDIRECT`** | [cite_start]Disable automatic redirection from HTTP to HTTPS[cite: 651]. |
| **`GOALERT_ENABLE_SECURE_HEADERS`** | [cite_start]Enable security headers like X-Frame-Options, CSP, etc.[cite: 652]. |
| **`GOALERT_STATUS_ADDR`** | [cite_start]Open a specific port to emit status updates (useful for container health checks)[cite: 1134]. |

---

### 🗄️ Database Tuning & Switchover
Variables for connection pooling and database migration/maintenance.

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_DB_MAX_IDLE`** | [cite_start]Max idle DB connections (Default: 5)[cite: 1061]. |
| **`GOALERT_DB_MAX_OPEN`** | [cite_start]Max open DB connections (Default: 15)[cite: 1065]. |
| **`GOALERT_DB_URL_NEXT`** | Connection string for the *next* Postgres server. [cite_start]Used only during DB switchover/migration[cite: 1048]. |
| **`GOALERT_DATA_ENCRYPTION_KEY_OLD`** | [cite_start]Fallback key for decrypting old data when rotating the main encryption key[cite: 1055]. |

---

### 📧 SMTP / Email Server (Ingress)
Configuration for the built-in SMTP server (receiving emails to trigger alerts).

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_SMTP_LISTEN`** | [cite_start]Address:port for the internal SMTP server[cite: 1117]. |
| **`GOALERT_SMTP_LISTEN_TLS`** | Address:port for SMTPS (secure). [cite_start]Requires TLS cert/key variables[cite: 1119]. |
| **`GOALERT_EMAIL_INTEGRATION_DOMAIN`** | [cite_start]**Required if SMTP is enabled.** The domain used for generating alert email addresses[cite: 1070]. |
| **`GOALERT_SMTP_ADDITIONAL_DOMAINS`** | [cite_start]Comma-separated list of extra allowed domains for incoming email[cite: 1114]. |
| **`GOALERT_SMTP_MAX_RECIPIENTS`** | [cite_start]Max recipients per message (Default: 1)[cite: 1121]. |
| **`GOALERT_SMTP_TLS_CERT_FILE`** | [cite_start]Path to PEM cert for SMTPS[cite: 1125]. |
| **`GOALERT_SMTP_TLS_KEY_FILE`** | [cite_start]Path to PEM key for SMTPS[cite: 1129]. |
| **`GOALERT_SMTP_TLS_CERT_DATA`** | [cite_start]PEM cert data string for SMTPS[cite: 1123]. |
| **`GOALERT_SMTP_TLS_KEY_DATA`** | [cite_start]PEM key data string for SMTPS[cite: 1127]. |

---

### 🛠️ Logging, Debugging & Dev
Variables for observability and development modes.

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_LOG_REQUESTS`** | [cite_start]Log all HTTP requests (Default: logs only debug/trace)[cite: 1102]. |
| **`GOALERT_LOG_ERRORS_ONLY`** | [cite_start]Only log errors (supersedes other log flags)[cite: 1100]. |
| **`GOALERT_JSON`** | [cite_start]Output logs in JSON format[cite: 1083]. |
| **`GOALERT_VERBOSE`** | [cite_start]Enable verbose logging[cite: 690]. |
| **`GOALERT_STACK_TRACES`** | [cite_start]Enable stack traces with all error logs[cite: 1132]. |
| **`GOALERT_LOG_ENGINE_CYCLES`** | [cite_start]Log the start and end of every engine cycle[cite: 1097]. |
| **`GOALERT_STUB_NOTIFIERS`** | Replace real senders (Twilio, etc.) with stubs that always succeed. [cite_start]Useful for staging[cite: 1139]. |
| **`GOALERT_UI_DIR`** | [cite_start]Serve UI assets from a local directory instead of embedded memory[cite: 689]. |

---

### ⚙️ Advanced & System
Variables for system tuning, limits, and experimental features.

| Environment Variable | Description |
| :--- | :--- |
| **`GOALERT_API_ONLY`** | Starts in API-only mode (no engine/processing). [cite_start]Useful for scaling in clusters[cite: 1058]. |
| **`GOALERT_REGION_NAME`** | [cite_start]Name of the region for message processing (Default: "default")[cite: 1109]. |
| **`GOALERT_ENGINE_CYCLE_TIME`** | [cite_start]Time between engine processing cycles (Default: 5s)[cite: 1071]. |
| **`GOALERT_MAX_REQUEST_BODY_BYTES`** | [cite_start]Max body size for incoming requests (Default: 262144)[cite: 1103]. |
| **`GOALERT_MAX_REQUEST_HEADER_BYTES`** | [cite_start]Max header size for incoming requests (Default: 4096)[cite: 1106]. |
| **`GOALERT_EXPERIMENTAL`** | [cite_start]Enable specific experimental features[cite: 1075]. |
| **`GOALERT_STRICT_EXPERIMENTAL`** | [cite_start]Fail to start if unknown experimental features are specified[cite: 1137]. |
| **`GOALERT_LIST_EXPERIMENTAL`** | [cite_start]List available experimental features[cite: 1085]. |
| **`GOALERT_LISTEN_PROMETHEUS`** | [cite_start]Bind address for Prometheus metrics[cite: 1090]. |
| **`GOALERT_LISTEN_SYSAPI`** | (Experimental) [cite_start]Listen address for system API (gRPC)[cite: 1092]. |
| **`GOALERT_GITHUB_BASE_URL`** | [cite_start]Base URL for GitHub auth/API (for Enterprise GitHub)[cite: 1077]. |
| **`GOALERT_SLACK_BASE_URL`** | [cite_start]Override Slack base URL[cite: 1113]. |
| **`GOALERT_TWILIO_BASE_URL`** | [cite_start]Override Twilio API URL[cite: 689]. |