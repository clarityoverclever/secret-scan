patterns = {
    {
        name = "AWS Access Key"
        regex = "AKAI[0-9A-Z]{16}"
        severity = "high"
    },
    {
        name = "Stripe Key"
        regex = "sk_live_[0-9a-zA-Z]+"
        severity = "high"
    },
}