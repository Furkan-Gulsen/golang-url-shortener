package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StatisticsRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewStatisticsRepository(ctx context.Context, tableName string) *StatisticsRepository {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &StatisticsRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *StatisticsRepository) Get(ctx context.Context, id string) (*domain.Statistics, error) {
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	var statistics domain.Statistics
	err = attributevalue.UnmarshalMap(result.Item, &statistics)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &statistics, nil
}

func (d *StatisticsRepository) Create(ctx context.Context, statistics domain.Statistics) error {
	item, err := attributevalue.MarshalMap(statistics)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	}

	_, err = d.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item to DynamoDB: %w", err)
	}

	return nil
}

func (d *StatisticsRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	_, err := d.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %w", err)
	}
	return nil
}
