package keeper_test

// func (s *KeeperTestSuite) TestSoftwareUpgrade() {
// 	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

// 	testCases := []struct {
// 		name      string
// 		req       *types.MsgSoftwareUpgrade
// 		expectErr bool
// 		errMsg    string
// 	}{
// 		{
// 			"unauthorized authority address",
// 			&types.MsgSoftwareUpgrade{
// 				Authority: s.addrs[0].String(),
// 				Plan: types.Plan{
// 					Name:   "all-good",
// 					Info:   "some text here",
// 					Height: 123450000,
// 				},
// 			},
// 			true,
// 			"expected gov account as only signer for proposal message",
// 		},
// 		{
// 			"invalid plan",
// 			&types.MsgSoftwareUpgrade{
// 				Authority: govAccAddr,
// 				Plan: types.Plan{
// 					Height: 123450000,
// 				},
// 			},
// 			true,
// 			"name cannot be empty: invalid request",
// 		},
// 		{
// 			"successful upgrade scheduled",
// 			&types.MsgSoftwareUpgrade{
// 				Authority: govAccAddr,
// 				Plan: types.Plan{
// 					Name:   "all-good",
// 					Info:   "some text here",
// 					Height: 123450000,
// 				},
// 			},
// 			false,
// 			"",
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			_, err := s.msgSrvr.SoftwareUpgrade(s.ctx, tc.req)
// 			if tc.expectErr {
// 				s.Require().Error(err)
// 				s.Require().Contains(err.Error(), tc.errMsg)
// 			} else {
// 				s.Require().NoError(err)
// 				plan, found := s.app.UpgradeKeeper.GetUpgradePlan(s.ctx)
// 				s.Require().Equal(true, found)
// 				s.Require().Equal(tc.req.Plan, plan)
// 			}
// 		})
// 	}
// }
