package test

import (
	"DIY2/mocks"
	"DIY2/models"
	"DIY2/services"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrderServiceImplTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	orderRepo        *mocks.MockIOrderRepo
	storeRepo        *mocks.MockIStoreRepo
	orderServiceImpl services.IOrders
}

func (suite *OrderServiceImplTestSuite) BeforeTest(suiteName, testName string) {
	ctrl := gomock.NewController(suite.T())
	suite.ctrl = ctrl
	suite.orderRepo = mocks.NewMockIOrderRepo(suite.ctrl)
	suite.storeRepo = mocks.NewMockIStoreRepo(suite.ctrl)
	suite.orderServiceImpl = services.NewOrder(suite.storeRepo, suite.orderRepo)
}

func (suite *OrderServiceImplTestSuite) AfterTest(suiteName, testName string) {
	suite.ctrl.Finish()
}

func (suite *OrderServiceImplTestSuite) TestGetTopProductsForStoreWithWrongStoreID() {
	var storeId int64
	storeId = 10
	expectedError := sql.ErrNoRows

	orderModel := models.OrderModel{StoreId: storeId}
	suite.orderRepo.EXPECT().GetTopOrders(&orderModel, 1, 5).Return(nil, expectedError)

	response, actualError := suite.orderServiceImpl.TopProductsInStore(storeId)

	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *OrderServiceImplTestSuite) TestGetTopProductsForStoreSuccess() {
	var storeId int64
	storeId = 1

	orderModel := models.OrderModel{StoreId: storeId}
	var orders []int64
	suite.orderRepo.EXPECT().GetTopOrders(&orderModel, 1, 5).Return(orders, nil)

	response, actualError := suite.orderServiceImpl.TopProductsInStore(storeId)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), response, orders)
}

func (suite *OrderServiceImplTestSuite) TestGetTopProductsForAllStoreSuccess() {

	orderModel1 := models.OrderModel{StoreId: 1}
	orderModel2 := models.OrderModel{StoreId: 2}

	var orders []int64

	var stores = []int64{1, 2}

	//expectedMap := map[string][]int64{"1": orders, "2": orders}

	var expectedResponse []services.TopProductsResponse
	storeModel := models.StoreModel{}
	suite.storeRepo.EXPECT().GetAllStores(&storeModel).Return(stores, nil)

	suite.orderRepo.EXPECT().GetTopOrders(&orderModel1, 1, 2).Return(orders, nil)
	suite.orderRepo.EXPECT().GetTopOrders(&orderModel2, 1, 2).Return(orders, nil)

	response, actualError := suite.orderServiceImpl.TopProductsForAllStores()

	assert.Nil(suite.T(), actualError)
	assert.NotEqual(suite.T(), response, expectedResponse)
}

func TestOrderFunctions(t *testing.T) {
	suite.Run(t, new(OrderServiceImplTestSuite))
}
