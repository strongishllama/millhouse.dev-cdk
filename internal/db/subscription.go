package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type Subscription struct {
	EmailAddress string    `json:"emailAddress"`
	ID           string    `json:"id"`
	IsConfirmed  bool      `json:"isConfirmed"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Create creates a new subscription.
func (s *Subscription) Create(ctx context.Context) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	if err := putItem(ctx, s); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// DeleteSubscription deletes a subscription via its email address and ID.
func DeleteSubscription(ctx context.Context, emailAddress, id string) error {
	if err := deleteItem(ctx, itemTypeSubscription, fmt.Sprintf("%s#%s", itemTypeSubscription, emailAddress), fmt.Sprintf("%s#%s", itemTypeSubscription, id)); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// GetSubscription fetches a subscription via its email address.
func GetSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	var subscription *Subscription
	if err := getItem(ctx, fmt.Sprintf("%s#%s", itemTypeSubscription, emailAddress), "", &subscription); err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscription, nil
}

// GetSubscriptions fetches a slice of subscriptions.
func GetSubscriptions(ctx context.Context) ([]*Subscription, error) {
	subscriptions := []*Subscription{}
	if err := getItems(ctx, itemTypeSubscription, &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscriptions, nil
}

// Update updates an existing subscription.
func (s *Subscription) Update(ctx context.Context) error {
	if err := updateItem(ctx, s); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (s *Subscription) pk() string {
	return fmt.Sprintf("%s#%s", itemTypeSubscription, s.EmailAddress)
}

func (s *Subscription) sk() string {
	return fmt.Sprintf("%s#%s", itemTypeSubscription, s.ID)
}

func (s *Subscription) countPK() string {
	return string(itemTypeCount)
}

func (s *Subscription) countSK() string {
	return fmt.Sprintf("%s#%s", itemTypeCount, s.itemType())
}

func (s *Subscription) itemType() itemType {
	return itemTypeSubscription
}

func (s *Subscription) updateExpression() (expression.Expression, error) {
	return expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("isConfirmed"),
			expression.Value(s.IsConfirmed),
		),
	).Build()
}

func (s *Subscription) validate() error {
	if len(s.EmailAddress) == 0 {
		return errors.New("email address cannot be empty")
	}
	if s.CreatedAt == (time.Time{}) {
		return errors.New("created at cannot be empty")
	}
	if s.UpdatedAt == (time.Time{}) {
		return errors.New("updated at cannot be empty")
	}

	return nil
}
