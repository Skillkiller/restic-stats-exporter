package snapshot

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCollector_Describe(t *testing.T) {
	c := &Collector{
		resticExecutablePath: "",
	}

	expectedDesc := map[string]bool{
		snapshotCountTotalDesc.String():              true,
		snapshotCountDesc.String():                   true,
		lastSnapshotTimeDesc.String():                true,
		lastSnapshotBackupStartDesc.String():         true,
		lastSnapshotBackupEndDesc.String():           true,
		lastSnapshotFilesNewDesc.String():            true,
		lastSnapshotFilesChangedDesc.String():        true,
		lastSnapshotFilesUnmodifiedDesc.String():     true,
		lastSnapshotDirsNewDesc.String():             true,
		lastSnapshotDirsChangedDesc.String():         true,
		lastSnapshotDirsUnmodifiedDesc.String():      true,
		lastSnapshotDataBlobsDesc.String():           true,
		lastSnapshotTreeBlobsDesc.String():           true,
		lastSnapshotDataAddedDesc.String():           true,
		lastSnapshotDataAddedPackedDesc.String():     true,
		lastSnapshotTotalFilesProcessedDesc.String(): true,
		lastSnapshotTotalBytesProcessedDesc.String(): true,
		snapshotExitCode.String():                    true,
	}

	expectedCount := len(expectedDesc)

	ch := make(chan *prometheus.Desc, expectedCount)
	done := make(chan struct{})

	go func() {
		c.Describe(ch)
		close(done)
	}()

	select {
	case <-done:
		// Describe has finished
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Describe blocked or did not return within timeout — possible extra sends")
	}

	close(ch)

	got := map[string]bool{}
	for d := range ch {
		if d == nil {
			t.Fatalf("received nil descriptor")
		}
		got[d.String()] = true
	}

	if len(got) != expectedCount {
		t.Fatalf("wrong number of descriptors: got %d, want %d", len(got), expectedCount)
	}

	for want := range expectedDesc {
		if !got[want] {
			t.Fatalf("descriptor not sent: %s", want)
		}
	}
}

func TestCollector_Collect_VariousResticExitCodes(t *testing.T) {
	type fields struct {
		exitCode int
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "exit code 1",
			fields: fields{
				exitCode: 1,
			},
		},
		{
			name: "exit code 10",
			fields: fields{
				exitCode: 1,
			},
		},
		{
			name: "exit code 12",
			fields: fields{
				exitCode: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fakeExec returns the specified exit code
			fakeExec := func(exe string, args ...string) ([]byte, error, int) {
				return []byte(``), errors.New("error for unit test"), tt.fields.exitCode
			}

			c := &Collector{
				resticExecutablePath: "",
				commandExecutor:      fakeExec,
			}

			reg := prometheus.NewRegistry()
			if err := reg.Register(c); err != nil {
				t.Fatalf("failed to register collector: %v", err)
			}

			expected := `
# HELP restic_snapshot_exit_code Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes
# TYPE restic_snapshot_exit_code gauge
restic_snapshot_exit_code %d
`

			err := testutil.CollectAndCompare(reg, strings.NewReader(fmt.Sprintf(expected, tt.fields.exitCode)))
			if err != nil {
				t.Fatalf("unexpected metrics output: %v", err)
			}
		})
	}
}

func TestCollector_Collect_InvalidJsonOutput(t *testing.T) {
	type fields struct {
		json string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "invalid 1",
			fields: fields{
				json: "{{",
			},
		},
		{
			name: "invalid 2",
			fields: fields{
				json: "[{",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fakeExec returns the specified JSON output as bytes
			fakeExec := func(exe string, args ...string) ([]byte, error, int) {
				return []byte(tt.fields.json), nil, 0
			}

			c := &Collector{
				resticExecutablePath: "",
				commandExecutor:      fakeExec,
			}

			reg := prometheus.NewRegistry()
			if err := reg.Register(c); err != nil {
				t.Fatalf("failed to register collector: %v", err)
			}

			expected := `
# HELP restic_snapshot_exit_code Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes
# TYPE restic_snapshot_exit_code gauge
restic_snapshot_exit_code 1684
`

			err := testutil.CollectAndCompare(reg, strings.NewReader(expected))
			if err != nil {
				t.Fatalf("unexpected metrics output: %v", err)
			}
		})
	}
}

func TestCollector_Collect_ExecutablePath(t *testing.T) {
	type fields struct {
		resticExecutablePath string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "restic",
			fields: fields{
				resticExecutablePath: "restic",
			},
		},
		{
			name: "/usr/bin/restic",
			fields: fields{
				resticExecutablePath: "/usr/bin/restic",
			},
		},
		{
			name: "./resticv123",
			fields: fields{
				resticExecutablePath: "./resticv123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fakeExec check if the executable path is as expected and return empty JSON array output
			fakeExec := func(exe string, args ...string) ([]byte, error, int) {
				if exe != tt.fields.resticExecutablePath {
					t.Fatalf("unexpected executable: %s", exe)
				}
				return []byte(`[]`), nil, 0
			}

			c := &Collector{
				resticExecutablePath: tt.fields.resticExecutablePath,
				commandExecutor:      fakeExec,
			}

			reg := prometheus.NewRegistry()
			if err := reg.Register(c); err != nil {
				t.Fatalf("failed to register collector: %v", err)
			}

			expected := `
# HELP restic_snapshot_exit_code Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes
# TYPE restic_snapshot_exit_code gauge
restic_snapshot_exit_code 0
# HELP restic_snapshot_count_total Total number of snapshots in the repository
# TYPE restic_snapshot_count_total gauge
restic_snapshot_count_total 0
`

			err := testutil.CollectAndCompare(reg, strings.NewReader(expected))
			if err != nil {
				t.Fatalf("unexpected metrics output: %v", err)
			}

		})
	}
}

func TestCollector_Collect_No_Snapshot(t *testing.T) {
	// fakeExec returns the specified JSON output with zero snapshots
	fakeExec := func(exe string, args ...string) ([]byte, error, int) {
		return []byte(`[]`), nil, 0
	}

	c := &Collector{
		resticExecutablePath: "restic",
		commandExecutor:      fakeExec,
	}

	reg := prometheus.NewRegistry()
	if err := reg.Register(c); err != nil {
		t.Fatalf("failed to register collector: %v", err)
	}

	expected := `
# HELP restic_snapshot_exit_code Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes
# TYPE restic_snapshot_exit_code gauge
restic_snapshot_exit_code 0
# HELP restic_snapshot_count_total Total number of snapshots in the repository
# TYPE restic_snapshot_count_total gauge
restic_snapshot_count_total 0
`

	err := testutil.CollectAndCompare(reg, strings.NewReader(expected))
	if err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}

func TestCollector_Collect_Multiple_Groups(t *testing.T) {
	// fakeExec returns the specified JSON output
	fakeExec := func(exe string, args ...string) ([]byte, error, int) {
		data, err := os.ReadFile("testdata/multiple_groups.json")
		if err != nil {
			t.Fatalf("readJson test file: %v", err)
		}

		return data, nil, 0
	}

	c := &Collector{
		resticExecutablePath: "restic",
		commandExecutor:      fakeExec,
	}

	reg := prometheus.NewRegistry()
	if err := reg.Register(c); err != nil {
		t.Fatalf("failed to register collector: %v", err)
	}

	expected := `
# HELP restic_last_snapshot_backup_end_seconds Unix timestamp: end time of the last backup
# TYPE restic_last_snapshot_backup_end_seconds gauge
restic_last_snapshot_backup_end_seconds{restic_hostname="DPC1",restic_tags="minebase"} 1.761495631e+09
restic_last_snapshot_backup_end_seconds{restic_hostname="DPC1",restic_tags="papermc"} 1.761506824e+09
# HELP restic_last_snapshot_backup_start_seconds Unix timestamp: start time of the last backup
# TYPE restic_last_snapshot_backup_start_seconds gauge
restic_last_snapshot_backup_start_seconds{restic_hostname="DPC1",restic_tags="minebase"} 1.761495627e+09
restic_last_snapshot_backup_start_seconds{restic_hostname="DPC1",restic_tags="papermc"} 1.76150682e+09
# HELP restic_last_snapshot_data_added_bytes Number of bytes added in the last snapshot (unpacked)
# TYPE restic_last_snapshot_data_added_bytes gauge
restic_last_snapshot_data_added_bytes{restic_hostname="DPC1",restic_tags="minebase"} 276436
restic_last_snapshot_data_added_bytes{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_data_added_packed_bytes Number of bytes added in the last snapshot (packed)
# TYPE restic_last_snapshot_data_added_packed_bytes gauge
restic_last_snapshot_data_added_packed_bytes{restic_hostname="DPC1",restic_tags="minebase"} 157995
restic_last_snapshot_data_added_packed_bytes{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_data_blobs Number of data blobs in the last snapshot
# TYPE restic_last_snapshot_data_blobs gauge
restic_last_snapshot_data_blobs{restic_hostname="DPC1",restic_tags="minebase"} 126
restic_last_snapshot_data_blobs{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_dirs_changed Number of changed directories in the last snapshot
# TYPE restic_last_snapshot_dirs_changed gauge
restic_last_snapshot_dirs_changed{restic_hostname="DPC1",restic_tags="minebase"} 0
restic_last_snapshot_dirs_changed{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_dirs_new Number of newly added directories in the last snapshot
# TYPE restic_last_snapshot_dirs_new gauge
restic_last_snapshot_dirs_new{restic_hostname="DPC1",restic_tags="minebase"} 92
restic_last_snapshot_dirs_new{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_dirs_unmodified Number of unmodified directories in the last snapshot
# TYPE restic_last_snapshot_dirs_unmodified gauge
restic_last_snapshot_dirs_unmodified{restic_hostname="DPC1",restic_tags="minebase"} 0
restic_last_snapshot_dirs_unmodified{restic_hostname="DPC1",restic_tags="papermc"} 391
# HELP restic_last_snapshot_files_changed Number of changed files in the last snapshot
# TYPE restic_last_snapshot_files_changed gauge
restic_last_snapshot_files_changed{restic_hostname="DPC1",restic_tags="minebase"} 0
restic_last_snapshot_files_changed{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_files_new Number of newly added files in the last snapshot
# TYPE restic_last_snapshot_files_new gauge
restic_last_snapshot_files_new{restic_hostname="DPC1",restic_tags="minebase"} 129
restic_last_snapshot_files_new{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_last_snapshot_files_unmodified Number of unmodified files in the last snapshot
# TYPE restic_last_snapshot_files_unmodified gauge
restic_last_snapshot_files_unmodified{restic_hostname="DPC1",restic_tags="minebase"} 0
restic_last_snapshot_files_unmodified{restic_hostname="DPC1",restic_tags="papermc"} 288
# HELP restic_last_snapshot_time_seconds Unix timestamp of the last snapshot
# TYPE restic_last_snapshot_time_seconds gauge
restic_last_snapshot_time_seconds{restic_hostname="DPC1",restic_tags="minebase"} 1.761495627e+09
restic_last_snapshot_time_seconds{restic_hostname="DPC1",restic_tags="papermc"} 1.76150682e+09
# HELP restic_last_snapshot_total_bytes_processed Total number of bytes processed in the last snapshot
# TYPE restic_last_snapshot_total_bytes_processed gauge
restic_last_snapshot_total_bytes_processed{restic_hostname="DPC1",restic_tags="minebase"} 126189
restic_last_snapshot_total_bytes_processed{restic_hostname="DPC1",restic_tags="papermc"} 2.02676503e+08
# HELP restic_last_snapshot_total_files_processed Total number of files processed in the last snapshot
# TYPE restic_last_snapshot_total_files_processed gauge
restic_last_snapshot_total_files_processed{restic_hostname="DPC1",restic_tags="minebase"} 129
restic_last_snapshot_total_files_processed{restic_hostname="DPC1",restic_tags="papermc"} 288
# HELP restic_last_snapshot_tree_blobs Number of tree blobs in the last snapshot
# TYPE restic_last_snapshot_tree_blobs gauge
restic_last_snapshot_tree_blobs{restic_hostname="DPC1",restic_tags="minebase"} 92
restic_last_snapshot_tree_blobs{restic_hostname="DPC1",restic_tags="papermc"} 0
# HELP restic_snapshot_count Number of snapshots
# TYPE restic_snapshot_count gauge
restic_snapshot_count{restic_hostname="DPC1",restic_tags="minebase"} 1
restic_snapshot_count{restic_hostname="DPC1",restic_tags="papermc"} 3
# HELP restic_snapshot_count_total Total number of snapshots in the repository
# TYPE restic_snapshot_count_total gauge
restic_snapshot_count_total 4
# HELP restic_snapshot_exit_code Exit code of the list snapshots command. See restic exit codes, except 1684 for json output parsing errors: https://restic.readthedocs.io/en/stable/075_scripting.html#exit-codes
# TYPE restic_snapshot_exit_code gauge
restic_snapshot_exit_code 0
`

	err := testutil.CollectAndCompare(reg, strings.NewReader(expected))
	if err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}
