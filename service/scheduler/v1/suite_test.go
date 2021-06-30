package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/itiky/charge_scheduler/service/scheduler/testutil"
	"github.com/itiky/charge_scheduler/storage/sqlite_base"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx    context.Context
	baseSt *sqlite_base.SQLiteBase
	r      *testutil.SchedulerServiceTestResource
}

func (s *ServiceTestSuite) SetupSuite() {
	baseSt, err := sqlite_base.SetupTempSQLiteBase(s.T().TempDir())
	if err != nil {
		panic(fmt.Errorf("base storage init: %w", err))
	}

	r, err := NewTestResource(baseSt)
	if err != nil {
		panic(fmt.Errorf("resource init: %w", err))
	}

	s.ctx = context.TODO()
	s.baseSt = baseSt
	s.r = r
}

// nolint:errcheck
func (s *ServiceTestSuite) TearDownSuite() {
	if s.baseSt != nil {
		s.baseSt.Close()
	}
}

func TestSuite_SchedulerService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
