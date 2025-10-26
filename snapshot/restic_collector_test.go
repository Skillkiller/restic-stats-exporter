package snapshot

import (
	"errors"
	"fmt"
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
