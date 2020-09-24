package awsdata_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	. "github.com/manywho/awsinventory/internal/awsdata"
	"github.com/manywho/awsinventory/internal/inventory"
)

var testElastiCacheNodeRows = []inventory.Row{
	{
		UniqueAssetIdentifier:          "test-cluster-1-test-node-0",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-0.test-cluster-1.ameaqx.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-1",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.t2.medium",
		SoftwareDatabaseVendor:         "redis",
		SoftwareDatabaseNameAndVersion: "redis 5.2",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-1",
		VLANNetworkID:                  "vpc-12345678",
	},
	{
		UniqueAssetIdentifier:          "test-cluster-1-test-node-1",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-1.test-cluster-1.ameaqx.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-1",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.t2.medium",
		SoftwareDatabaseVendor:         "redis",
		SoftwareDatabaseNameAndVersion: "redis 5.2",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-1",
		VLANNetworkID:                  "vpc-12345678",
	},
	{
		UniqueAssetIdentifier:          "test-cluster-2-test-node-2",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-2.test-cluster-2.fnjyzo.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-2",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.t2.small",
		SoftwareDatabaseVendor:         "memcached",
		SoftwareDatabaseNameAndVersion: "memcached 3.2.10",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-2",
		VLANNetworkID:                  "vpc-12345678",
	},
	{
		UniqueAssetIdentifier:          "test-cluster-2-test-node-3",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-3.test-cluster-2.fnjyzo.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-2",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.t2.small",
		SoftwareDatabaseVendor:         "memcached",
		SoftwareDatabaseNameAndVersion: "memcached 3.2.10",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-2",
		VLANNetworkID:                  "vpc-12345678",
	},
	{
		UniqueAssetIdentifier:          "test-cluster-3-test-node-4",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-4.test-cluster-3.7wufxa.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-3",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.m4.large",
		SoftwareDatabaseVendor:         "redis",
		SoftwareDatabaseNameAndVersion: "redis 2.3",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-3",
		VLANNetworkID:                  "vpc-12345678",
	},
	{
		UniqueAssetIdentifier:          "test-cluster-3-test-node-5",
		Virtual:                        true,
		Public:                         false,
		DNSNameOrURL:                   "test-node-5.test-cluster-3.7wufxa.use1.cache.amazonaws.com",
		BaselineConfigurationName:      "cache-parameter-group-3",
		Location:                       DefaultRegion,
		AssetType:                      "ElastiCache Node",
		HardwareMakeModel:              "cache.m4.large",
		SoftwareDatabaseVendor:         "redis",
		SoftwareDatabaseNameAndVersion: "redis 2.3",
		SerialAssetTagNumber:           "arn:aws:elasticache:us-east-1:123456789012:cluster:test-cluster-3",
		VLANNetworkID:                  "vpc-12345678",
	},
}

// Test Data
var testElastiCacheNodeOutput = &elasticache.DescribeCacheClustersOutput{
	CacheClusters: []*elasticache.CacheCluster{
		{
			ARN:                  aws.String(testElastiCacheNodeRows[0].SerialAssetTagNumber),
			CacheClusterId:       aws.String("test-cluster-1"),
			Engine:               aws.String("redis"),
			EngineVersion:        aws.String("5.2"),
			CacheNodeType:        aws.String(testElastiCacheNodeRows[0].HardwareMakeModel),
			CacheSubnetGroupName: aws.String("cache-subnet-group"),
			CacheParameterGroup: &elasticache.CacheParameterGroupStatus{
				CacheParameterGroupName: aws.String("cache-parameter-group-1"),
			},
			CacheNodes: []*elasticache.CacheNode{
				{
					CacheNodeId: aws.String("test-node-0"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[0].DNSNameOrURL),
					},
				},
				{
					CacheNodeId: aws.String("test-node-1"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[1].DNSNameOrURL),
					},
				},
			},
		},
		{
			ARN:                  aws.String(testElastiCacheNodeRows[2].SerialAssetTagNumber),
			CacheClusterId:       aws.String("test-cluster-2"),
			Engine:               aws.String("memcached"),
			EngineVersion:        aws.String("3.2.10"),
			CacheNodeType:        aws.String(testElastiCacheNodeRows[2].HardwareMakeModel),
			CacheSubnetGroupName: aws.String("cache-subnet-group"),
			CacheParameterGroup: &elasticache.CacheParameterGroupStatus{
				CacheParameterGroupName: aws.String("cache-parameter-group-2"),
			},
			CacheNodes: []*elasticache.CacheNode{
				{
					CacheNodeId: aws.String("test-node-2"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[2].DNSNameOrURL),
					},
				},
				{
					CacheNodeId: aws.String("test-node-3"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[3].DNSNameOrURL),
					},
				},
			},
		},
		{
			ARN:                  aws.String(testElastiCacheNodeRows[4].SerialAssetTagNumber),
			CacheClusterId:       aws.String("test-cluster-3"),
			Engine:               aws.String("redis"),
			EngineVersion:        aws.String("2.3"),
			CacheNodeType:        aws.String(testElastiCacheNodeRows[4].HardwareMakeModel),
			CacheSubnetGroupName: aws.String("cache-subnet-group"),
			CacheParameterGroup: &elasticache.CacheParameterGroupStatus{
				CacheParameterGroupName: aws.String("cache-parameter-group-3"),
			},
			CacheNodes: []*elasticache.CacheNode{
				{
					CacheNodeId: aws.String("test-node-4"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[4].DNSNameOrURL),
					},
				},
				{
					CacheNodeId: aws.String("test-node-5"),
					Endpoint: &elasticache.Endpoint{
						Address: aws.String(testElastiCacheNodeRows[5].DNSNameOrURL),
					},
				},
			},
		},
	},
}

var testElastiCacheSubnetGroupOutput = &elasticache.DescribeCacheSubnetGroupsOutput{
	CacheSubnetGroups: []*elasticache.CacheSubnetGroup{
		{
			CacheSubnetGroupName: aws.String("cache-subnet-group"),
			VpcId: aws.String(testElastiCacheNodeRows[0].VLANNetworkID),
		},
	},
}

// Mocks
type ElastiCacheMock struct {
	elasticacheiface.ElastiCacheAPI
}

func (e ElastiCacheMock) DescribeCacheClusters(cfg *elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error) {
	return testElastiCacheNodeOutput, nil
}

func (e ElastiCacheMock) DescribeCacheSubnetGroups(cfg *elasticache.DescribeCacheSubnetGroupsInput) (*elasticache.DescribeCacheSubnetGroupsOutput, error) {
	return testElastiCacheSubnetGroupOutput, nil
}

type ElastiCacheErrorMock struct {
	elasticacheiface.ElastiCacheAPI
}

func (e ElastiCacheErrorMock) DescribeCacheClusters(cfg *elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error) {
	return &elasticache.DescribeCacheClustersOutput{}, testError
}

func (e ElastiCacheErrorMock) DescribeCacheSubnetGroups(cfg *elasticache.DescribeCacheSubnetGroupsInput) (*elasticache.DescribeCacheSubnetGroupsOutput, error) {
	return &elasticache.DescribeCacheSubnetGroupsOutput{}, testError
}

// Tests
func TestCanLoadElastiCacheNodes(t *testing.T) {
	d := New(logrus.New(), TestClients{ElastiCache: ElastiCacheMock{}})

	d.Load([]string{DefaultRegion}, []string{ServiceElastiCache})

	var count int
	d.MapRows(func(row inventory.Row) error {
		require.Equal(t, testElastiCacheNodeRows[count], row)
		count++
		return nil
	})
	require.Equal(t, 6, count)
}

func TestLoadElastiCacheNodesLogsError(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{ElastiCache: ElastiCacheErrorMock{}})

	d.Load([]string{DefaultRegion}, []string{ServiceElastiCache})

	require.Contains(t, hook.LastEntry().Message, testError.Error())
}