package k8s_test

import (
	"testing"

	"zkctl/cmd/pkg/k8s"

	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
)

type KubeTaintTestSuite struct {
	suite.Suite
	TaintedNodes []*k8s.IncompatibleNode
}

func (suite *KubeTaintTestSuite) SetupSuite() {
	suite.TaintedNodes = []*k8s.IncompatibleNode{
		{
			NodeSummary: &k8s.NodeSummary{
				Taints: []v1.Taint{
					{
						Key:    "test",
						Value:  "test",
						Effect: "NoSchedule",
					},
					{
						Key:    "good",
						Value:  "good",
						Effect: "NoSchedule",
					},
				},
			},
		},
		{
			NodeSummary: &k8s.NodeSummary{
				Taints: []v1.Taint{
					{
						Key:    "bad",
						Value:  "bad",
						Effect: "NoSchedule",
					},
				},
			},
		},
		{
			NodeSummary: &k8s.NodeSummary{
				Taints: []v1.Taint{
					{
						Key:    "bad",
						Value:  "bad",
						Effect: "NoSchedule",
					},
				},
			},
		},
	}
}

func (suite *KubeTaintTestSuite) TearDownSuite() {}

func TestKubeTaintTestSuite(t *testing.T) {
	suite.Run(t, &KubeTaintTestSuite{})
}

func (suite *KubeTaintTestSuite) TestGetTaintsSuccess() {
	// prepare
	tolerationManager := &k8s.TolerationManager{
		TaintedNodes: suite.TaintedNodes,
	}

	// act
	taints, err := tolerationManager.GetTaints()
	suite.NoError(err)

	// assert

	expected := []string{
		"{\"key\":\"test\",\"value\":\"test\",\"effect\":\"NoSchedule\"}",
		"{\"key\":\"good\",\"value\":\"good\",\"effect\":\"NoSchedule\"}",
		"{\"key\":\"bad\",\"value\":\"bad\",\"effect\":\"NoSchedule\"}",
	}

	suite.ElementsMatch(expected, taints)
}

func (suite *KubeTaintTestSuite) TestGetTolerationsSuccess() {
	// prepare
	tolerationManager := &k8s.TolerationManager{
		TaintedNodes: suite.TaintedNodes,
	}

	// act
	tolerations, err := tolerationManager.GetTolerationsMap([]string{"{\"key\":\"test\",\"value\":\"test\",\"effect\":\"NoSchedule\"}"})
	suite.NoError(err)

	// assert
	var tolerationSeconds *int64
	expected := []map[string]interface{}{
		{
			"key":               "test",
			"value":             "test",
			"operator":          v1.TolerationOpEqual,
			"effect":            v1.TaintEffectNoSchedule,
			"tolerationSeconds": tolerationSeconds,
		},
	}

	suite.Equal(expected, tolerations)
}

func (suite *KubeTaintTestSuite) TestGetTolerableNodesSuccess() {
	// prepare
	tolerationManager := &k8s.TolerationManager{
		TaintedNodes: suite.TaintedNodes,
	}

	// act
	allowedTaints := []string{
		"{\"key\":\"test\",\"value\":\"test\",\"effect\":\"NoSchedule\"}",
		"{\"key\":\"good\",\"value\":\"good\",\"effect\":\"NoSchedule\"}",
	}

	nodes, err := tolerationManager.GetTolerableNodes(allowedTaints)
	suite.NoError(err)

	// assert

	expected := []*k8s.NodeSummary{
		suite.TaintedNodes[0].NodeSummary,
	}

	suite.Equal(expected, nodes)
}
