CREATE TABLE analytic
(
    ID               UInt64,
    RequestUUID      UUID,
    UserUUID         UUID,
    Timestamp        DateTime,
    Method           LowCardinality(String),
    Host             LowCardinality(String),
    Path             String,
    Query            String,
    Status           UInt16,
    Latency          UInt32,
    IP               IPv4,
    Country          LowCardinality(Nullable(String)),
    City             Nullable(String),
    Browser          LowCardinality(String),
    Device           LowCardinality(String),
    OS               LowCardinality(String),
    Language         LowCardinality(String),
    Referrer         String,
    ResolutionWidth  Nullable(UInt16),
    ResolutionHeight Nullable(UInt16),
    RequestSize      UInt32,
    ResponseSize     UInt32,
    UTMCampaign      LowCardinality(String),
    UTMSource        LowCardinality(String),
    UTMMedium        LowCardinality(String),

    INDEX            ip_idx IP TYPE bloom_filter GRANULARITY 8192,
    INDEX            user_uuid_idx UserUUID TYPE bloom_filter GRANULARITY 8192
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(Timestamp)
ORDER BY (Timestamp, RequestUUID)
TTL Timestamp + INTERVAL 1 YEAR;
