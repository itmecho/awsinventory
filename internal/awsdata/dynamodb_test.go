package awsdata_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	. "github.com/manywho/awsinventory/internal/awsdata"
	"github.com/manywho/awsinventory/internal/inventory"
)

var testDynamoDBTableRows = []inventory.Row{
	{
		UniqueAssetIdentifier:  "TestTable1",
		Virtual:                true,
		Location:               DefaultRegion,
		AssetType:              AssetTypeDynamoDBTable,
		SoftwareDatabaseVendor: "Amazon",
		Comments:               "100 B",
		SerialAssetTagNumber:   "arn:aws:dynamodb:us-east-1:123456789012:table/TestTable1",
	},
	{
		UniqueAssetIdentifier:  "TestTable2",
		Virtual:                true,
		Location:               DefaultRegion,
		AssetType:              AssetTypeDynamoDBTable,
		SoftwareDatabaseVendor: "Amazon",
		Comments:               "50.0 kB",
		SerialAssetTagNumber:   "arn:aws:dynamodb:us-east-1:123456789012:table/TestTable2",
	},
	{
		UniqueAssetIdentifier:  "TestTable3",
		Virtual:                true,
		Location:               DefaultRegion,
		AssetType:              AssetTypeDynamoDBTable,
		SoftwareDatabaseVendor: "Amazon",
		Comments:               "20.0 MB",
		SerialAssetTagNumber:   "arn:aws:dynamodb:us-east-1:123456789012:table/TestTable3",
	},
}

// Test Data
var testDynamoDBListTablesOutputPage1 = &dynamodb.ListTablesOutput{
    LastEvaluatedTableName: aws.String(testDynamoDBTableRows[1].UniqueAssetIdentifier),
	TableNames: []*string{
		aws.String(testDynamoDBTableRows[0].UniqueAssetIdentifier),
		aws.String(testDynamoDBTableRows[1].UniqueAssetIdentifier),
	},
}


var testDynamoDBListTablesOutputPage2 = &dynamodb.ListTablesOutput{
    LastEvaluatedTableName: nil,
	TableNames: []*string{
		aws.String(testDynamoDBTableRows[2].UniqueAssetIdentifier),
	},
}

// Mocks
type DynamoDBMock struct {
	dynamodbiface.DynamoDBAPI
}

func (e DynamoDBMock) ListTables(cfg *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
    if cfg.ExclusiveStartTableName == nil {
        return testDynamoDBListTablesOutputPage1, nil
    }

	return testDynamoDBListTablesOutputPage2, nil
}

func (e DynamoDBMock) DescribeTable(cfg *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	var row int
	var bytes int64
	switch aws.StringValue(cfg.TableName) {
	case testDynamoDBTableRows[0].UniqueAssetIdentifier:
		row = 0
		bytes = 100
	case testDynamoDBTableRows[1].UniqueAssetIdentifier:
		row = 1
		bytes = 51200
	case testDynamoDBTableRows[2].UniqueAssetIdentifier:
		row = 2
		bytes = 20971520
	}
	return &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableArn:       aws.String(testDynamoDBTableRows[row].SerialAssetTagNumber),
			TableName:      aws.String(testDynamoDBTableRows[row].UniqueAssetIdentifier),
			TableSizeBytes: aws.Int64(bytes),
		},
	}, nil
}

type DynamoDBErrorMock struct {
	dynamodbiface.DynamoDBAPI
}

func (e DynamoDBErrorMock) ListTables(cfg *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	return &dynamodb.ListTablesOutput{}, testError
}

func (e DynamoDBErrorMock) DescribeTable(cfg *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return &dynamodb.DescribeTableOutput{}, testError
}

// Tests
func TestCanLoadDynamoDBTables(t *testing.T) {
	d := New(logrus.New(), TestClients{DynamoDB: DynamoDBMock{}})

	d.Load([]string{DefaultRegion}, []string{ServiceDynamoDB})

	var count int
	d.MapRows(func(row inventory.Row) error {
		require.Equal(t, testDynamoDBTableRows[count], row)
		count++
		return nil
	})
	require.Equal(t, 3, count)
}

func TestLoadDynamoDBTablesLogsError(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{DynamoDB: DynamoDBErrorMock{}})

	d.Load([]string{DefaultRegion}, []string{ServiceDynamoDB})

	require.Contains(t, hook.LastEntry().Message, testError.Error())
}
