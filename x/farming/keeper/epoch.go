package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/x/farming/types"
)

func (k Keeper) GetLastEpochTime(ctx sdk.Context) (time.Time, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastEpochTimeKey)
	if bz == nil {
		return time.Time{}, false
	}
	var ts gogotypes.Timestamp
	k.cdc.MustUnmarshal(bz, &ts)
	t, err := gogotypes.TimestampFromProto(&ts)
	if err != nil {
		panic(err)
	}
	return t, true
}

func (k Keeper) SetLastEpochTime(ctx sdk.Context, t time.Time) {
	store := ctx.KVStore(k.storeKey)
	ts, err := gogotypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	bz := k.cdc.MustMarshal(ts)
	store.Set(types.LastEpochTimeKey, bz)
}

func (k Keeper) AdvanceEpoch(ctx sdk.Context) error {
	if err := k.AllocateRewards(ctx); err != nil {
		return err
	}
	k.ProcessQueuedCoins(ctx)
	k.SetLastEpochTime(ctx, ctx.BlockTime())

	return nil
}

// GetCurrentEpochDays returns the  current epoch days.
func (k Keeper) GetCurrentEpochDays(ctx sdk.Context) uint32 {
	var epochDays uint32
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentEpochDaysKey)
	if bz == nil {
		// initialize with next epoch days
		epochDays = k.GetParams(ctx).NextEpochDays
	} else {
		val := gogotypes.UInt32Value{}
		if err := k.cdc.Unmarshal(bz, &val); err != nil {
			panic(err)
		}
		epochDays = val.GetValue()
	}
	return epochDays
}

// SetCurrentEpochDays sets the  current epoch days.
func (k Keeper) SetCurrentEpochDays(ctx sdk.Context, epochDays uint32) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt32Value{Value: epochDays})
	store.Set(types.CurrentEpochDaysKey, bz)
}
