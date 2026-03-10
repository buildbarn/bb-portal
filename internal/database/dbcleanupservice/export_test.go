package dbcleanupservice

// SetCompactLogsBatchSizeForTest overrides the compaction batch size for tests.
func SetCompactLogsBatchSizeForTest(batchSize int) {
	compactLogsBatchSize = batchSize
}

// CompactLogsBatchSizeForTest returns the current compaction batch size for tests.
func CompactLogsBatchSizeForTest() int {
	return compactLogsBatchSize
}
