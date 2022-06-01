package account

import (
	"context"

	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/logic/errcode"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/account"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/repo/globalrpc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/svc"
	"github.com/zecrey-labs/zecrey-legend/service/api/app/internal/types"
	"github.com/zecrey-labs/zecrey-legend/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountStatusByAccountNameLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	globalRPC globalrpc.GlobalRPC
	account   account.AccountModel
}

func NewGetAccountStatusByAccountNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountStatusByAccountNameLogic {
	return &GetAccountStatusByAccountNameLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		account:   account.New(svcCtx.Config),
		globalRPC: globalrpc.New(svcCtx.Config, ctx),
	}
}

func (l *GetAccountStatusByAccountNameLogic) GetAccountStatusByAccountName(req *types.ReqGetAccountStatusByAccountName) (resp *types.RespGetAccountStatusByAccountName, err error) {
	if utils.CheckAccountName(req.AccountName) {
		logx.Error("[CheckAccountIndex] param:%v", req.AccountName)
		return nil, errcode.ErrInvalidParam
	}
	account, err := l.account.GetAccountByAccountName(req.AccountName)
	if err != nil {
		logx.Error("[GetAccountByAccountName] err:%v", err)
		return nil, err
	}
	resp = &types.RespGetAccountStatusByAccountName{
		AccountStatus: uint32(account.Status),
		AccountPk:     account.PublicKey,
		AccountIndex:  uint32(account.AccountIndex),
	}
	return resp, nil

}