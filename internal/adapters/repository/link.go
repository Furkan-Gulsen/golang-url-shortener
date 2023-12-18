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

type LinkRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewLinkRepository(ctx context.Context, tableName string) *LinkRepository {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &LinkRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *LinkRepository) All(ctx context.Context) ([]domain.Link, error) {
	var links []domain.Link

	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
		Limit:     aws.Int32(20),
	}

	result, err := d.client.Scan(ctx, input)

	if err != nil {
		return links, fmt.Errorf("failed to get items from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return links, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return links, nil
}

func (d *LinkRepository) Get(ctx context.Context, id string) (domain.Link, error) {
	link := domain.Link{}

	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return link, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalMap(result.Item, &link)
	if err != nil {
		return link, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return link, nil
}

func (d *LinkRepository) Create(ctx context.Context, link domain.Link) error {
	item, err := attributevalue.MarshalMap(link)
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

func (d *LinkRepository) Delete(ctx context.Context, id string) error {
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
