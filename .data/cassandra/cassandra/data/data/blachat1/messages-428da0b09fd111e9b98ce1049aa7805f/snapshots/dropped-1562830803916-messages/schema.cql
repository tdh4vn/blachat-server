CREATE TABLE IF NOT EXISTS blachat1.messages (
	id uuid,
	channel_id uuid,
	created_at timestamp,
	author_id uuid,
	content text,
	system boolean,
	type int,
	PRIMARY KEY ((id, channel_id), created_at))
	WITH ID = 428da0b0-9fd1-11e9-b98c-e1049aa7805f
	AND CLUSTERING ORDER BY (created_at DESC)
	AND bloom_filter_fp_chance = 0.01
	AND dclocal_read_repair_chance = 0.1
	AND crc_check_chance = 1.0
	AND default_time_to_live = 0
	AND gc_grace_seconds = 864000
	AND min_index_interval = 128
	AND max_index_interval = 2048
	AND memtable_flush_period_in_ms = 0
	AND read_repair_chance = 0.0
	AND speculative_retry = '99PERCENTILE'
	AND comment = ''
	AND caching = { 'keys': 'ALL', 'rows_per_partition': 'NONE' }
	AND compaction = { 'max_threshold': '32', 'min_threshold': '4', 'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy' }
	AND compression = { 'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor' }
	AND cdc = false
	AND extensions = {  };
CREATE INDEX messages_channel_id_idx ON blachat1.messages (channel_id);
