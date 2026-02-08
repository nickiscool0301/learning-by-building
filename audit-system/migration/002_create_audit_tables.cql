USE audit;

CREATE TABLE IF NOT EXISTS audit_events_by_entity {
    entity_id TEXT,
    time_bucket TEST,
    event_time TIMEUUID,
    event_id UUID,
    event_type TEXT,
    actior_id TEXT,
    actor_type TEXT,
    description TEXT,
    metadata MAP<TEXT, TEXT>,
    ip_address TEXT,
    PRIMARY KEY ((entity_id, time_bucket), event_time)
} WITH CLUSTERING ORDER BY (event_time DESC) AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 1
}

CREATE TABLE IF NOT EXISTS audit_events_by_actor {
    actor_id TEXT,
    time_bucket TEST,
    event_time TIMEUUID,
    event_id UUID,
    event_type TEXT,
    entity_id TEXT,
    description TEXT,
    metadata MAP<TEXT, TEXT>,
    ip_address TEXT,
    PRIMARY KEY ((actor_id, time_bucket), event_time)
} WITH CLUSTERING ORDER BY (event_time DESC) AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 1
}

CREATE TABLE IF NOT EXISTS audit_events_by_type {
    event_type TEXT,
    time_bucket TEST,
    event_time TIMEUUID,
    event_id UUID,
    entity_id TEXT,
    actor_id TEXT,
    description TEXT,
    metadata MAP<TEXT, TEXT>,
    ip_address TEXT,
    PRIMARY KEY ((event_type, time_bucket), event_time)
} WITH CLUSTERING ORDER BY (event_time DESC) AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 1
}
