package scenarios

import (
	"context"
	"fmt"
	"io"
	"os"

	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"px.dev/pxapi/types"
)

func PixieClient(apiKey string, clusterId string, cloudAddress string, pxlPath string) error {
	var API_KEY = apiKey
	var CLUSTER_ID = clusterId
	var CLOUD_ADDR = cloudAddress

	// fmt.Println("PX API key : ", apiKey)
	// fmt.Println("PX cluster ID : ", clusterId)
	// fmt.Println("Cloud address : ", cloudAddress)

	dat, err := os.ReadFile(pxlPath)
	pxl := string(dat)
	// fmt.Print(pxl)
	if err != nil {
		return err
	}

	// Create a Pixie client.
	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(API_KEY), pxapi.WithCloudAddr(CLOUD_ADDR))
	if err != nil {
		return err
	}

	// Create a connection to the cluster.
	vz, err := client.NewVizierClient(ctx, CLUSTER_ID)
	if err != nil {
		return err
	}

	// Create TableMuxer to accept results table.
	tm := &tableMux{}

	// Execute the PxL script.
	resultSet, err := vz.ExecuteScript(ctx, pxl, tm)
	if err != nil && err != io.EOF {
		return err
	}

	// Receive the PxL script results.
	defer resultSet.Close()
	if err := resultSet.Stream(); err != nil {
		if errdefs.IsCompilationError(err) {
			fmt.Printf("Got compiler error: \n %s\n", err.Error())
		} else {
			println("Error")
			fmt.Printf("Got error : %+v, while streaming\n", err)
		}
	}

	// Get the execution stats for the script execution.
	stats := resultSet.Stats()
	fmt.Printf("Execution Time: %v\n", stats.ExecutionTime)
	fmt.Printf("Bytes received: %v\n", stats.TotalBytes)

	return nil
}

// Satisfies the TableRecordHandler interface.
type tablePrinter struct{}

func (t *tablePrinter) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *tablePrinter) HandleRecord(ctx context.Context, r *types.Record) error {
	for _, d := range r.Data {
		fmt.Printf("%s ", d.String())
	}
	fmt.Printf("\n")
	return nil
}

func (t *tablePrinter) HandleDone(ctx context.Context) error {
	return nil
}

// Satisfies the TableMuxer interface.
type tableMux struct {
}

func (s *tableMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &tablePrinter{}, nil
}
