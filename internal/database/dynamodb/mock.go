package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDBClientMock struct {
	dynamodbiface.DynamoDBAPI
	GetItemFn         func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	GetItemInvoked    bool
	ScanFn            func(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	ScanInvoked       bool
	PutItemFn         func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	PutItemInvoked    bool
	DeleteItemFn      func(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	DeleteItemInvoked bool
}

func (m *DynamoDBClientMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	m.GetItemInvoked = true
	return m.GetItemFn(input)
}

func (m *DynamoDBClientMock) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	m.ScanInvoked = true
	return m.ScanFn(input)
}

func (m *DynamoDBClientMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.PutItemInvoked = true
	return m.PutItemFn(input)
}

func (m *DynamoDBClientMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	m.DeleteItemInvoked = true
	return m.DeleteItemFn(input)
}
