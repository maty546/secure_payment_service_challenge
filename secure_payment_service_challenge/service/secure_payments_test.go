package service

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/maty546/secure_payment_service_challenge/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSecurePaymentsService_StartTransfer(t *testing.T) {
	testCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	mockCtrl := gomock.NewController(t)
	mockAccRepo := mocks.NewMockIAccountRepository(mockCtrl)
	mockTrRepo := mocks.NewMockITransferRepository(mockCtrl)

	service := NewService(mockAccRepo, mockTrRepo, "", "")

	type mock struct {
		times    int
		response interface{}
		err      error
	}

	type want struct {
		expectedTransfer models.Transfer
		expectedErr      error
	}
	nowTime := time.Now()
	testTransfer := models.Transfer{ID: 1, Amount: 10, FromAccountID: "ext_a", ToAccountID: "b", CreatedAt: nowTime, UpdatedAt: nowTime}
	testAcc := models.Account{ID: "b", Balance: 10}
	expTransfer := models.Transfer{ID: 1, Amount: 10, FromAccountID: "ext_a", ToAccountID: "b", Status: models.TRANSFER_STATUS_PENDING, CreatedAt: nowTime, UpdatedAt: nowTime}

	tests := []struct {
		description    string
		transfer       models.Transfer
		transferToSave models.Transfer
		accMock        mock
		trMock         mock
		want           want
	}{

		{description: "happy path external to internal",
			transfer:       testTransfer,
			transferToSave: expTransfer,
			accMock:        mock{times: 1, response: testAcc},
			trMock:         mock{times: 1, response: expTransfer},
			want:           want{expectedTransfer: expTransfer, expectedErr: nil}},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mockAccRepo.EXPECT().GetByID(gomock.Any(), tc.transfer.ToAccountID).Return(tc.accMock.response, tc.accMock.err).Times(tc.accMock.times)
			mockTrRepo.EXPECT().Save(gomock.Any(), tc.transferToSave).Return(tc.trMock.response, tc.trMock.err).Times(tc.trMock.times)

			resultTransfer, err := service.StartTransfer(testCtx, tc.transfer)

			assert.Equal(t, tc.want.expectedTransfer, resultTransfer)
			assert.Equal(t, tc.want.expectedErr, err)
		})
	}
}
