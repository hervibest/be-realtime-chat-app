How to run cassandra migration : 

CREATE KEYSPACE messaging_service
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};
