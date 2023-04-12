package terraform

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/databricks/bricks/bundle"
	"github.com/databricks/bricks/libs/filer"
	"github.com/databricks/bricks/libs/log"
	"github.com/databricks/databricks-sdk-go/apierr"
)

type statePull struct{}

func (l *statePull) Name() string {
	return "terraform:state-pull"
}

func (l *statePull) Apply(ctx context.Context, b *bundle.Bundle) ([]bundle.Mutator, error) {
	f, err := filer.NewWorkspaceFilesClient(b.WorkspaceClient(), b.Config.Workspace.StatePath)
	if err != nil {
		return nil, err
	}

	dir, err := Dir(b)
	if err != nil {
		return nil, err
	}

	// Download state file from filer to local cache directory.
	log.Infof(ctx, "Opening remote state file")
	remote, err := f.Read(ctx, TerraformStateFileName)
	if err != nil {
		// On first deploy this state file doesn't yet exist.
		if apierr.IsMissing(err) {
			log.Infof(ctx, "Remote state file does not exist")
			return nil, nil
		}
		return nil, err
	}

	// Expect the state file to live under dir.
	local, err := os.OpenFile(filepath.Join(dir, TerraformStateFileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	// Write file to disk.
	log.Infof(ctx, "Writing remote state file to local cache directory")
	_, err = io.Copy(local, remote)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func StatePull() bundle.Mutator {
	return &statePull{}
}