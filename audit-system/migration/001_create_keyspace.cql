CREATE KEYSPACE IF NOT EXISTS audit WITH replication = {
'class': 'NetworkTopologyStrategy',
'dc1': '3'
};
