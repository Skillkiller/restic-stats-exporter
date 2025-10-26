package snapshot

import "github.com/prometheus/client_golang/prometheus"

var (
	snapshotCountTotalDesc = prometheus.NewDesc(
		"restic_snapshot_count_total",
		"Total number of snapshots in the repository",
		nil, nil,
	)
	snapshotCountDesc = prometheus.NewDesc(
		"restic_snapshot_count",
		"Number of snapshots",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTimeDesc = prometheus.NewDesc(
		"restic_last_snapshot_time_seconds",
		"Unix timestamp of the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotBackupStartDesc = prometheus.NewDesc(
		"restic_last_snapshot_backup_start_seconds",
		"Unix timestamp: start time of the last backup",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotBackupEndDesc = prometheus.NewDesc(
		"restic_last_snapshot_backup_end_seconds",
		"Unix timestamp: end time of the last backup",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesNewDesc = prometheus.NewDesc(
		"restic_last_snapshot_files_new",
		"Number of newly added files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesChangedDesc = prometheus.NewDesc(
		"restic_last_snapshot_files_changed",
		"Number of changed files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotFilesUnmodifiedDesc = prometheus.NewDesc(
		"restic_last_snapshot_files_unmodified",
		"Number of unmodified files in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsNewDesc = prometheus.NewDesc(
		"restic_last_snapshot_dirs_new",
		"Number of newly added directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsChangedDesc = prometheus.NewDesc(
		"restic_last_snapshot_dirs_changed",
		"Number of changed directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDirsUnmodifiedDesc = prometheus.NewDesc(
		"restic_last_snapshot_dirs_unmodified",
		"Number of unmodified directories in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataBlobsDesc = prometheus.NewDesc(
		"restic_last_snapshot_data_blobs",
		"Number of data blobs in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTreeBlobsDesc = prometheus.NewDesc(
		"restic_last_snapshot_tree_blobs",
		"Number of tree blobs in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataAddedDesc = prometheus.NewDesc(
		"restic_last_snapshot_data_added_bytes",
		"Number of bytes added in the last snapshot (unpacked)",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotDataAddedPackedDesc = prometheus.NewDesc(
		"restic_last_snapshot_data_added_packed_bytes",
		"Number of bytes added in the last snapshot (packed)",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTotalFilesProcessedDesc = prometheus.NewDesc(
		"restic_last_snapshot_total_files_processed",
		"Total number of files processed in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	lastSnapshotTotalBytesProcessedDesc = prometheus.NewDesc(
		"restic_last_snapshot_total_bytes_processed",
		"Total number of bytes processed in the last snapshot",
		[]string{"restic_hostname", "restic_tags"}, nil,
	)

	snapshotExitCode = prometheus.NewDesc("restic_snapshot_exit_code",
		"Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: "+
			"https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes",
		nil, nil)
)
