package transaction

import (
	"context"
	"strconv"

	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/accounthistory"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/block"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/globalrpc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/mempool"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/tx"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/svc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTxsByPubKeyLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	account   accounthistory.AccountHistory
	globalRpc globalrpc.GlobalRPC
	tx        tx.Tx
	mempool   mempool.Mempool
	block     block.Block
}

func NewGetTxsByPubKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTxsByPubKeyLogic {
	return &GetTxsByPubKeyLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		account:   accounthistory.New(svcCtx.Config),
		globalRpc: globalrpc.New(svcCtx.Config, ctx),
		tx:        tx.New(svcCtx.Config),
		mempool:   mempool.New(svcCtx.Config),
		block:     block.New(svcCtx.Config),
	}
}

func (l *GetTxsByPubKeyLogic) GetTxsByPubKey(req *types.ReqGetTxsByPubKey) (resp *types.RespGetTxsByPubKey, err error) {
	//err = utils.CheckRequestParam(utils.TypeAccountPk, reflect.ValueOf(req.Pk))

	//err = utils.CheckRequestParam(utils.TypeLimit, reflect.ValueOf(req.Limit))

	//err = utils.CheckRequestParam(utils.TypeOffset, reflect.ValueOf(req.Offset))

	account, err := l.account.GetAccountByPk(req.Pk)
	if err != nil {
		logx.Error("[transaction.GetTxsByPubKey] err:%v", err)
		return &types.RespGetTxsByPubKey{}, err
	}

	txList, _, err := l.globalRpc.GetLatestTxsListByAccountIndex(uint32(account.AccountIndex), req.Limit)
	if err != nil {
		logx.Error("[transaction.GetTxsByPubKey] err:%v", err)
		return &types.RespGetTxsByPubKey{}, err
	}
	txCount, err := l.tx.GetTxsTotalCountByAccountIndex(account.AccountIndex)
	if err != nil {
		logx.Error("[transaction.GetTxsByPubKey] err:%v", err)
		return &types.RespGetTxsByPubKey{}, err
	}

	mempoolTxCount, err := l.mempool.GetMempoolTxsTotalCountByAccountIndex(account.AccountIndex)
	if err != nil {
		logx.Error("[transaction.GetTxsByPubKey] err:%v", err)
		return &types.RespGetTxsByPubKey{}, err
	}

	results := make([]*types.Tx, 0)
	for _, tx := range txList {
		txDetails := make([]*types.TxDetail, 0)
		for _, txDetail := range tx.MempoolDetails {
			txDetails = append(txDetails, &types.TxDetail{
				AssetId:        uint32(txDetail.AssetId),
				AssetType:      uint32(txDetail.AssetType),
				AccountIndex:   int32(txDetail.AccountIndex),
				AccountName:    txDetail.AccountName,
				AccountBalance: txDetail.BalanceDelta,
			})
		}
		txAmount, _ := strconv.ParseInt(tx.TxAmount, 10, 64)
		blockInfo, err := l.block.GetBlockByBlockHeight(tx.L2BlockHeight)
		if err != nil {
			logx.Error("[transaction.GetTxsByPubKey] err:%v", err)
			return &types.RespGetTxsByPubKey{}, err
		}
		results = append(results, &types.Tx{
			TxHash:        tx.TxHash,
			TxType:        uint32(tx.TxType),
			GasFeeAssetId: uint32(tx.GasFeeAssetId),
			GasFee:        tx.GasFee,
			TxStatus:      uint32(tx.Status),
			BlockHeight:   uint32(tx.L2BlockHeight),
			BlockStatus:   uint32(blockInfo.BlockStatus),
			BlockId:       uint32(blockInfo.ID),
			//Todo: do we still need AssetAId and AssetBId
			AssetAId:      uint32(tx.AssetId),
			AssetBId:      uint32(tx.AssetId),
			TxAmount:      uint32(txAmount),
			TxDetails:     txDetails,
			NativeAddress: tx.NativeAddress,
			CreatedAt:     tx.CreatedAt.UnixNano() / 1e6,
			Memo:          tx.Memo,
		})
	}
	return &types.RespGetTxsByPubKey{
		Total: uint32(txCount + mempoolTxCount),
		Txs:   results,
	}, nil
}
