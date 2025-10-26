package snapshot

import "github.com/prometheus/client_golang/prometheus"

var (
	snapshotCountTotal = prometheus.NewDesc(
		"restic_snapshot_count_total",
		"Total number of snapshots in the repository",
		nil, nil,
	)
	snapshotCount = prometheus.NewDesc(
		"restic_snapshot_count",
		"Number of snapshots",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTime = prometheus.NewDesc(
		"restic_last_snapshot_time_seconds",
		"Unix timestamp of the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotBackupStart = prometheus.NewDesc(
		"restic_last_snapshot_backup_start_seconds",
		"Unix timestamp: start time of the last backup",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotBackupEnd = prometheus.NewDesc(
		"restic_last_snapshot_backup_end_seconds",
		"Unix timestamp: end time of the last backup",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesNew = prometheus.NewDesc(
		"restic_last_snapshot_files_new",
		"Number of newly added files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesChanged = prometheus.NewDesc(
		"restic_last_snapshot_files_changed",
		"Number of changed files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesUnmodified = prometheus.NewDesc(
		"restic_last_snapshot_files_unmodified",
		"Number of unmodified files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsNew = prometheus.NewDesc(
		"restic_last_snapshot_dirs_new",
		"Number of newly added directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsChanged = prometheus.NewDesc(
		"restic_last_snapshot_dirs_changed",
		"Number of changed directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsUnmodified = prometheus.NewDesc(
		"restic_last_snapshot_dirs_unmodified",
		"Number of unmodified directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataBlobs = prometheus.NewDesc(
		"restic_last_snapshot_data_blobs",
		"Number of data blobs in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTreeBlobs = prometheus.NewDesc(
		"restic_last_snapshot_tree_blobs",
		"Number of tree blobs in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataAdded = prometheus.NewDesc(
		"restic_last_snapshot_data_added_bytes",
		"Number of bytes added in the last snapshot (unpacked)",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataAddedPacked = prometheus.NewDesc(
		"restic_last_snapshot_data_added_packed_bytes",
		"Number of bytes added in the last snapshot (packed)",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTotalFilesProcessed = prometheus.NewDesc(
		"restic_last_snapshot_total_files_processed",
		"Total number of files processed in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTotalBytesProcessed = prometheus.NewDesc(
		"restic_last_snapshot_total_bytes_processed",
		"Total number of bytes processed in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)
)
