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
	"github.com/aws/aws-sdk-go/aws"
)

type StatsRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewStatsRepository(ctx context.Context, tableName string) *StatsRepository {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &StatsRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *StatsRepository) Get(ctx context.Context, id string) (domain.Stats, error) {
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	stats := domain.Stats{}
	err = attributevalue.UnmarshalMap(result.Item, &stats)
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return stats, nil
}

func (d *StatsRepository) All(ctx context.Context) ([]domain.Stats, error) {
	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	stats := []domain.Stats{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return stats, nil
}

func (d *StatsRepository) Create(ctx context.Context, stats domain.Stats) error {
	item, err := attributevalue.MarshalMap(stats)
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

func (d *StatsRepository) Delete(ctx context.Context, id string) error {
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

func (d *StatsRepository) GetStatsByLinkID(ctx context.Context, linkID string) ([]domain.Stats, error) {
	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":linkID": &ddbtypes.AttributeValueMemberS{Value: linkID},
		},
		FilterExpression: aws.String("link_id = :linkID"),
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	stats := []domain.Stats{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return stats, nil
}
