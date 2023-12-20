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

func (d *StatsRepository) GetByLinkID(ctx context.Context, linkID string) (domain.Stats, error) {
	input := &dynamodb.ScanInput{
		TableName:        &d.tableName,
		FilterExpression: aws.String("link_id = :link_id"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":link_id": &ddbtypes.AttributeValueMemberS{Value: linkID},
		},
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to query table: %w", err)
	}

	if len(result.Items) == 0 {
		return domain.Stats{}, fmt.Errorf("no stats found for the given linkID")
	}

	var stats domain.Stats
	err = attributevalue.UnmarshalMap(result.Items[0], &stats)
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return stats, nil
}

func (d *StatsRepository) IncreaseClickCountByLinkID(ctx context.Context, linkID string) error {
	queryInput := &dynamodb.QueryInput{
		TableName:              &d.tableName,
		IndexName:              aws.String("linkID-index"),
		KeyConditionExpression: aws.String("link_id = :linkID"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":linkID": &ddbtypes.AttributeValueMemberS{Value: linkID},
		},
		Limit: aws.Int32(1),
	}

	queryResult, err := d.client.Query(ctx, queryInput)
	if err != nil {
		return fmt.Errorf("failed to query items from DynamoDB: %w", err)
	}

	if len(queryResult.Items) == 0 {
		return fmt.Errorf("no item found with link_id: %s", linkID)
	}

	primaryID := queryResult.Items[0]["id"].(*ddbtypes.AttributeValueMemberS).Value

	updateInput := &dynamodb.UpdateItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: primaryID},
		},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":inc": &ddbtypes.AttributeValueMemberN{Value: "1"},
		},
		UpdateExpression: aws.String("ADD clickCount :inc"),
	}

	_, err = d.client.UpdateItem(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update item from DynamoDB: %w", err)
	}

	return nil
}
