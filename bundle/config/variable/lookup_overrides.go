package variable

import (
	"context"
	"fmt"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/compute"
)

var lookupOverrides = map[string]resolverFunc{
	"Cluster": resolveCluster,
}

// We added a custom resolver for the cluster to add filtering for the cluster source when we list all clusters.
// Without the filtering listing could take a very long time (5-10 mins) which leads to lookup timeouts.
func resolveCluster(ctx context.Context, w *databricks.WorkspaceClient, name string) (string, error) {
	result, err := w.Clusters.ListAll(ctx, compute.ListClustersRequest{
		FilterBy: &compute.ListClustersFilterBy{
			ClusterSources: []compute.ClusterSource{compute.ClusterSourceApi, compute.ClusterSourceUi},
		},
	})

	if err != nil {
		return "", err
	}

	tmp := map[string][]compute.ClusterDetails{}
	for _, v := range result {
		key := v.ClusterName
		tmp[key] = append(tmp[key], v)
	}
	alternatives, ok := tmp[name]
	if !ok || len(alternatives) == 0 {
		return "", fmt.Errorf("cluster named '%s' does not exist", name)
	}
	if len(alternatives) > 1 {
		return "", fmt.Errorf("there are %d instances of clusters named '%s'", len(alternatives), name)
	}
	return alternatives[0].ClusterId, nil
}