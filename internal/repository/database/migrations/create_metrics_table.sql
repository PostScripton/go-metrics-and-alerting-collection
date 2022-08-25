CREATE TABLE metrics_and_alerting.metrics
(
    id    VARCHAR(255)     NOT NULL,
    type  VARCHAR(30)      NOT NULL,
    delta BIGINT           NULL,
    value DOUBLE PRECISION NULL,
    PRIMARY KEY (id, type)
);
