
patterns = {
    -- Cloud Provider Keys (Critical - Most Common)
    {
        name = "AWS Access Key ID",
        regex = "AKIA[0-9A-Z]{16}",
        severity = "critical"
    },
    {
        name = "AWS Secret Key",
        regex = "aws.{0,20}['\"][0-9a-zA-Z/+]{40}['\"]",
        severity = "critical"
    },
    {
        name = "Google Cloud API Key",
        regex = "AIza[0-9A-Za-z\\-_]{35}",
        severity = "critical"
    },
    {
       name = "Azure Subscription Key",
       regex = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
       severity = "medium",
       validator = "azure_context"
    },

    -- Payment Processing (Critical)
    {
        name = "Stripe Live API Key",
        regex = "sk_live_[0-9a-zA-Z]{24,}",
        severity = "critical"
    },
    {
        name = "Stripe Restricted API Key",
        regex = "rk_live_[0-9a-zA-Z]{24,}",
        severity = "critical"
    },
    {
        name = "PayPal Braintree Access Token",
        regex = "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
        severity = "critical"
    },
    {
        name = "Square Access Token",
        regex = "sq0atp-[0-9A-Za-z\\-_]{22}",
        severity = "critical"
    },

    -- Authentication & Tokens (Critical)
    {
        name = "Generic API Key",
        regex = "[aA][pP][iI]_?[kK][eE][yY].*['\"][0-9a-zA-Z]{32,45}['\"]",
        severity = "high"
    },
    {
        name = "JSON Web Token (JWT)",
        regex = "eyJ[A-Za-z0-9-_=]+\\.eyJ[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_.+/=]*",
        severity = "high"
    },
    {
        name = "GitHub Personal Access Token",
        regex = "ghp_[0-9a-zA-Z]{36}",
        severity = "critical"
    },
    {
        name = "GitHub OAuth Token",
        regex = "gho_[0-9a-zA-Z]{36}",
        severity = "critical"
    },
    {
        name = "GitLab Personal Access Token",
        regex = "glpat-[0-9a-zA-Z\\-]{20}",
        severity = "critical"
    },

    -- Database Connections (Critical)
    {
        name = "Database Connection String",
        regex = "(mongodb|mysql|postgres|postgresql)://[^\\s]+:[^\\s]+@[^\\s]+",
        severity = "critical"
    },
    {
        name = "JDBC Connection String",
        regex = "jdbc:[^\\s]+:[^\\s]+://[^\\s]+",
        severity = "high"
    },

    -- Private Keys (Critical)
    {
        name = "RSA Private Key",
        regex = "-----BEGIN RSA PRIVATE KEY-----",
        severity = "critical"
    },
    {
        name = "SSH Private Key",
        regex = "-----BEGIN OPENSSH PRIVATE KEY-----",
        severity = "critical"
    },
    {
        name = "PGP Private Key",
        regex = "-----BEGIN PGP PRIVATE KEY BLOCK-----",
        severity = "critical"
    },
    {
        name = "RSA Key (Base64)",
        regex = "[A-Za-z0-9+/]{200,}={0,2}",
        severity = "high",
        validator = "base64_high_entropy"
    },

    -- Messaging & Communication (High)
    {
        name = "Slack Webhook",
        regex = "https://hooks\\.slack\\.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
        severity = "high"
    },
    {
        name = "Slack Bot Token",
        regex = "xoxb-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}",
        severity = "high"
    },
    {
        name = "Twilio API Key",
        regex = "SK[0-9a-fA-F]{32}",
        severity = "high"
    },
    {
        name = "SendGrid API Key",
        regex = "SG\\.[0-9A-Za-z\\-_]{22}\\.[0-9A-Za-z\\-_]{43}",
        severity = "high"
    },

    -- Generic Credentials (Medium - High False Positives)
    {
        name = "Password in Code",
        regex = "[pP][aA][sS][sS][wW][oO][rR][dD]\\s*[=:]\\s*['\"][^'\"\\s]{8,}['\"]",
        severity = "medium"
    },
    {
        name = "Base64 Encoded Secret",
        regex = "[sS][eE][cC][rR][eE][tT]\\s*[=:]\\s*['\"][A-Za-z0-9+/]{100,}={0,2}['\"]",
        severity = "high"
    },
    {
        name = "Generic Secret",
        regex = "[sS][eE][cC][rR][eE][tT]\\s*[=:]\\s*['\"][^'\"\\s]{16,}['\"]",
        severity = "medium"
    },
    {
        name = "Authorization Bearer Token",
        regex = "[Bb]earer [a-zA-Z0-9_\\-\\.=]{20,}",
        severity = "high"
    },
}