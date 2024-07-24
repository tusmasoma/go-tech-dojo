package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-tech-dojo/domain/repository/mock"
)

type CreateUserAndTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

func TestUserUseCase_CreateUserAndToken(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		arg     CreateUserAndTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m1.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			arg: CreateUserAndTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m2 *mock.MockTransactionRepository) {
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: CreateUserAndTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr)
			jwt, err := usecase.CreateUserAndToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}
