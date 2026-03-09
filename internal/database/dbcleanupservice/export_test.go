package dbcleanupservice

// SetDeleteOldInvocationsBatchSizeForTest overrides the deletion batch size for tests.
func SetDeleteOldInvocationsBatchSizeForTest(batchSize int64) {
	deleteOldInvocationsBatchSize = batchSize
}

// DeleteOldInvocationsBatchSizeForTest returns the current deletion batch size for tests.
func DeleteOldInvocationsBatchSizeForTest() int64 {
	return deleteOldInvocationsBatchSize
}
