package sync

import (
	"bytes"
	"testing"

	"github.com/PlakarKorp/kloset/objects"
	"github.com/PlakarKorp/kloset/repository"
	ptesting "github.com/PlakarKorp/plakar/testing"
	"github.com/stretchr/testify/require"
)

func TestExecuteCmdSyncReportsSynchronizationFailure(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "destination maintenance lock"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := setupSync(t, nil, nil)
			snap := ptesting.GenerateSnapshot(t, fixture.localRepo, mockFiles)
			t.Cleanup(func() { require.NoError(t, snap.Close()) })

			lockID := objects.RandomMAC()
			lock := repository.NewExclusiveLock("maintenance")
			lockData := &bytes.Buffer{}
			require.NoError(t, lock.SerializeToStream(lockData))
			_, err := fixture.peerRepo.PutLock(lockID, lockData)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, fixture.peerRepo.DeleteLock(lockID)) })

			cmd := &Sync{}
			require.NoError(t, cmd.Parse(fixture.localCtx, []string{"to", fixture.peerArg}))

			status, err := cmd.Execute(fixture.localCtx, fixture.localRepo)
			require.Error(t, err)
			require.Equal(t, 1, status)
		})
	}
}
