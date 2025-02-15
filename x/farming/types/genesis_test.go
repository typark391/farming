package types_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/farming/x/farming/types"
)

func TestValidateGenesis(t *testing.T) {
	validAcc := sdk.AccAddress(crypto.AddressHash([]byte("validAcc")))
	validStakingCoinDenom := "denom1"
	validPlan := types.NewRatioPlan(
		types.NewBasePlan(
			1,
			"planA",
			types.PlanTypePublic,
			validAcc.String(),
			validAcc.String(),
			sdk.NewDecCoins(
				sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
			),
			types.ParseTime("0001-01-01T00:00:00Z"),
			types.ParseTime("9999-12-31T00:00:00Z"),
		),
		sdk.NewDecWithPrec(5, 2),
	)
	validStaking := types.Staking{
		Amount:        sdk.NewInt(1000000),
		StartingEpoch: 1,
	}
	validQueuedStaking := types.QueuedStaking{
		Amount: sdk.NewInt(1000000),
	}
	validHistoricalRewards := types.HistoricalRewards{
		CumulativeUnitRewards: sdk.NewDecCoins(sdk.NewInt64DecCoin("denom3", 100000)),
	}
	validOutstandingRewards := types.OutstandingRewards{
		Rewards: sdk.NewDecCoins(sdk.NewInt64DecCoin("denom3", 1000000)),
	}

	testCases := []struct {
		name        string
		configure   func(*types.GenesisState)
		expectedErr string
	}{
		{
			"default case",
			func(genState *types.GenesisState) {
				params := types.DefaultParams()
				genState.Params = params
			},
			"",
		},
		{
			"invalid NextEpochDays case",
			func(genState *types.GenesisState) {
				params := types.DefaultParams()
				params.NextEpochDays = 0
				genState.Params = params
			},
			"next epoch days must be positive: 0",
		},
		{
			"invalid plan",
			func(genState *types.GenesisState) {
				plan := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypeNil,
						validAcc.String(),
						validAcc.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.NewDecWithPrec(5, 2),
				)
				planAny, _ := types.PackPlan(plan)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:             *planAny,
						FarmingPoolCoins: sdk.NewCoins(),
					},
				}
			},
			"unknown plan type: PLAN_TYPE_UNSPECIFIED: invalid plan type",
		},
		{
			"invalid plan records - empty type url",
			func(genState *types.GenesisState) {
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:             cdctypes.Any{},
						FarmingPoolCoins: sdk.NewCoins(),
					},
				}
			},
			"empty type url: invalid type",
		},
		{
			"invalid plan records - invalid farming pool coins",
			func(genState *types.GenesisState) {
				planAny, _ := types.PackPlan(validPlan)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:             *planAny,
						FarmingPoolCoins: sdk.Coins{sdk.NewInt64Coin("denom3", 0)},
					},
				}
			},
			"coin 0denom3 amount is not positive",
		},
		{
			"not sorted plan ids",
			func(genState *types.GenesisState) {
				planA := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypePublic,
						validAcc.String(),
						validAcc.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.NewDecWithPrec(5, 2),
				)
				planB := types.NewFixedAmountPlan(
					types.NewBasePlan(
						2,
						"planB",
						types.PlanTypePublic,
						validAcc.String(),
						validAcc.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.NewCoins(sdk.NewInt64Coin("denom3", 1000000)),
				)
				planAAny, _ := types.PackPlan(planA)
				planBAny, _ := types.PackPlan(planB)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:             *planBAny,
						FarmingPoolCoins: sdk.NewCoins(),
					},
					{
						Plan:             *planAAny,
						FarmingPoolCoins: sdk.NewCoins(),
					},
				}
			},
			"pool records must be sorted",
		},
		{
			"invalid plan records - invalid sum of epoch ratio",
			func(genState *types.GenesisState) {
				planA := types.NewRatioPlan(
					types.NewBasePlan(
						1,
						"planA",
						types.PlanTypePublic,
						validAcc.String(),
						validAcc.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planB := types.NewRatioPlan(
					types.NewBasePlan(
						2,
						"planB",
						types.PlanTypePublic,
						validAcc.String(),
						validAcc.String(),
						sdk.NewDecCoins(
							sdk.NewInt64DecCoin(validStakingCoinDenom, 1),
						),
						types.ParseTime("0001-01-01T00:00:00Z"),
						types.ParseTime("9999-12-31T00:00:00Z"),
					),
					sdk.OneDec(),
				)
				planAAny, _ := types.PackPlan(planA)
				planBAny, _ := types.PackPlan(planB)
				genState.PlanRecords = []types.PlanRecord{
					{
						Plan:             *planAAny,
						FarmingPoolCoins: sdk.NewCoins(),
					},
					{
						Plan:             *planBAny,
						FarmingPoolCoins: sdk.NewCoins(),
					},
				}
			},
			"total epoch ratio must be lower than 1: invalid request",
		},
		{
			"invalid staking records - invalid staking coin denom",
			func(genState *types.GenesisState) {
				genState.StakingRecords = []types.StakingRecord{
					{
						StakingCoinDenom: "!",
						Farmer:           validAcc.String(),
						Staking:          validStaking,
					},
				}
			},
			"invalid denom: !",
		},
		{
			"invalid staking records - invalid farmer addr",
			func(genState *types.GenesisState) {
				genState.StakingRecords = []types.StakingRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						Farmer:           "invalid",
						Staking:          validStaking,
					},
				}
			},
			"decoding bech32 failed: invalid bech32 string length 7",
		},
		{
			"invalid staking records - invalid staking amount",
			func(genState *types.GenesisState) {
				genState.StakingRecords = []types.StakingRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						Farmer:           validAcc.String(),
						Staking: types.Staking{
							Amount:        sdk.ZeroInt(),
							StartingEpoch: 0,
						},
					},
				}
			},
			"staking amount must be positive: 0",
		},
		{
			"invalid queued staking records - invalid staking coin denom",
			func(genState *types.GenesisState) {
				genState.QueuedStakingRecords = []types.QueuedStakingRecord{
					{
						StakingCoinDenom: "!",
						Farmer:           validAcc.String(),
						QueuedStaking:    validQueuedStaking,
					},
				}
			},
			"invalid denom: !",
		},
		{
			"invalid queued staking records - invalid farmer addr",
			func(genState *types.GenesisState) {
				genState.QueuedStakingRecords = []types.QueuedStakingRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						Farmer:           "invalid",
						QueuedStaking:    validQueuedStaking,
					},
				}
			},
			"decoding bech32 failed: invalid bech32 string length 7",
		},
		{
			"invalid queued staking records - invalid queued staking amount",
			func(genState *types.GenesisState) {
				genState.QueuedStakingRecords = []types.QueuedStakingRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						Farmer:           validAcc.String(),
						QueuedStaking: types.QueuedStaking{
							Amount: sdk.ZeroInt(),
						},
					},
				}
			},
			"queued staking amount must be positive: 0",
		},
		{
			"invalid historical rewards records - invalid staking coin denom",
			func(genState *types.GenesisState) {
				genState.HistoricalRewardsRecords = []types.HistoricalRewardsRecord{
					{
						StakingCoinDenom:  "!",
						Epoch:             0,
						HistoricalRewards: validHistoricalRewards,
					},
				}
			},
			"invalid denom: !",
		},
		{
			"invalid historical rewards records - invalid historical rewards",
			func(genState *types.GenesisState) {
				genState.HistoricalRewardsRecords = []types.HistoricalRewardsRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						Epoch:            0,
						HistoricalRewards: types.HistoricalRewards{
							CumulativeUnitRewards: sdk.DecCoins{sdk.NewInt64DecCoin("denom3", 0)},
						},
					},
				}
			},
			"coin 0.000000000000000000denom3 amount is not positive",
		},
		{
			"invalid outstanding rewards records - invalid staking coin denom",
			func(genState *types.GenesisState) {
				genState.OutstandingRewardsRecords = []types.OutstandingRewardsRecord{
					{
						StakingCoinDenom:   "!",
						OutstandingRewards: validOutstandingRewards,
					},
				}
			},
			"invalid denom: !",
		},
		{
			"invalid outstanding rewards records - invalid outstanding rewards",
			func(genState *types.GenesisState) {
				genState.OutstandingRewardsRecords = []types.OutstandingRewardsRecord{
					{
						StakingCoinDenom: validStakingCoinDenom,
						OutstandingRewards: types.OutstandingRewards{
							Rewards: sdk.DecCoins{sdk.NewInt64DecCoin("denom3", 0)},
						},
					},
				}
			},
			"coin 0.000000000000000000denom3 amount is not positive",
		},
		{
			"invalid current epoch records - invalid staking coin denom",
			func(genState *types.GenesisState) {
				genState.CurrentEpochRecords = []types.CurrentEpochRecord{
					{
						StakingCoinDenom: "!",
						CurrentEpoch:     0,
					},
				}
			},
			"invalid denom: !",
		},
		{
			"invalid staking reserve coins",
			func(genState *types.GenesisState) {
				genState.StakingReserveCoins = sdk.Coins{sdk.NewInt64Coin(validStakingCoinDenom, 0)}
			},
			"coin 0denom1 amount is not positive",
		},
		{
			"invalid reward pool coins",
			func(genState *types.GenesisState) {
				genState.RewardPoolCoins = sdk.Coins{sdk.NewInt64Coin(validStakingCoinDenom, 0)}
			},
			"coin 0denom1 amount is not positive",
		},
		//{
		//	"invalid current epoch days",
		//	func(genState *types.GenesisState) {
		//		genState.CurrentEpochDays = 0
		//	},
		//	"not implemented",
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genState := types.DefaultGenesisState()
			tc.configure(genState)

			err := types.ValidateGenesis(*genState)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
