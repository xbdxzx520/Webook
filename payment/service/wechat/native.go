package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/payment/domain"
	"gitee.com/geekbang/basic-go/webook/payment/events"
	"gitee.com/geekbang/basic-go/webook/payment/repository"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
)

var errUnknownTransactionState = errors.New("未知的微信事务状态")

type NativePaymentService struct {
	appID string
	mchID string
	// 支付通知回调 URL
	notifyURL string
	// 自己的支付记录
	repo repository.PaymentRepository
	// 作业使用的 repo
	// 注意初始化的时候要传入正确的 repo
	repov1  *repository.PaymentGORMRepository
	msgRepo repository.LocalMsgRepository

	svc      *native.NativeApiService
	producer events.Producer

	l logger.LoggerV1

	// 在微信 native 里面，分别是
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（付款码支付）
	// USERPAYING：用户支付中（付款码支付）
	// PAYERROR：支付失败(其他原因，如银行返回失败)
	nativeCBTypeToStatus map[string]domain.PaymentStatus
}

func NewNativePaymentService(appID string, mchID string,
	repo repository.PaymentRepository, svc *native.NativeApiService,
	l logger.LoggerV1) *NativePaymentService {
	return &NativePaymentService{appID: appID, mchID: mchID, notifyURL: "http://wechat.meoying.com/pay/callback",
		repo: repo, svc: svc, l: l,
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFailed,
			"NOTPAY":   domain.PaymentStatusInit,
			"CLOSED":   domain.PaymentStatusFailed,
			"REVOKED":  domain.PaymentStatusFailed,
			"REFUND":   domain.PaymentStatusRefund,
			// 其它状态你都可以加
		},
	}
}

func (n *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	pmt.Status = domain.PaymentStatusInit
	err := n.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	//sn := uuid.New().String()
	resp, _, err := n.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(n.appID),
		Mchid:       core.String(n.mchID),
		Description: core.String(pmt.Description),
		OutTradeNo:  core.String(pmt.BizTradeNO),
		// 最好这个要带上
		TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
		Amount: &native.Amount{
			Total:    core.Int64(pmt.Amt.Total),
			Currency: core.String(pmt.Amt.Currency),
		},
	})

	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, nil
}

func (n *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNO string) error {
	// 对账
	txn, _, err := n.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNO),
		Mchid:      core.String(n.mchID),
	})
	if err != nil {
		return err
	}
	return n.updateByTxn(ctx, txn)
}

func (n *NativePaymentService) FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]domain.Payment, error) {
	return n.repo.FindExpiredPayment(ctx, offset, limit, t)
}

func (n *NativePaymentService) GetPayment(ctx context.Context, bizTradeId string) (domain.Payment, error) {
	return n.repo.GetPayment(ctx, bizTradeId)
}

func (n *NativePaymentService) HandleCallback(ctx context.Context, txn *payments.Transaction) error {
	return n.updateByTxn(ctx, txn)
}

func (n *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, 微信的状态是 %s", errUnknownTransactionState, *txn.TradeState)
	}
	// 很显然，就是更新一下我们本地数据库里面 payment 的状态
	err := n.repo.UpdatePayment(ctx, domain.Payment{
		// 微信过来的 transaction id
		TxnID:      *txn.TransactionId,
		BizTradeNO: *txn.OutTradeNo,
		Status:     status,
	})
	if err != nil {
		return err
	}
	// 就要通知业务方了
	// 有些人的系统，会根据支付状态来决定要不要通知
	// 我要是发消息失败了怎么办？
	// 站在业务的角度，你是不是至少应该发成功一次
	err1 := n.producer.ProducePaymentEvent(ctx, events.PaymentEvent{
		BizTradeNO: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	})
	if err1 != nil {
		n.l.Error("发送支付事件失败", logger.Error(err),
			logger.String("biz_trade_no", *txn.OutTradeNo))
	}
	return nil
}

// updateByTxnV1 第十七周作业
// 使用本地消息表的关键就是要把更新支付状态和插入代发送消息在一个数据库事务内完成操作
// 这一步我们下沉到了 repository，来规避在 service 上操作本地数据库事务
func (n *NativePaymentService) updateByTxnV1(ctx context.Context, txn *payments.Transaction) error {
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, 微信的状态是 %s", errUnknownTransactionState, *txn.TradeState)
	}
	// 很显然，就是更新一下我们本地数据库里面 payment 的状态
	evt := events.PaymentEvent{
		BizTradeNO: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	}
	var msgId int64
	// 在这种情况下，你无可避免会和底层耦合在一起
	err := n.repov1.Transaction(ctx, func(pmt *repository.PaymentGORMRepository, msg *repository.LocalMsgGORMRepository) error {
		err1 := pmt.UpdatePayment(ctx, domain.Payment{
			// 微信过来的 transaction id
			TxnID:      *txn.TransactionId,
			BizTradeNO: *txn.OutTradeNo,
			Status:     status,
		})
		if err1 != nil {
			return err1
		}
		evtData, err1 := json.Marshal(evt)
		if err1 != nil {
			return err1
		}
		msgId, err1 = msg.AddMsg(ctx, string(evtData))
		return err1
	})
	if err != nil {
		return err
	}

	// 就要通知业务方了
	// 有些人的系统，会根据支付状态来决定要不要通知
	// 我要是发消息失败了怎么办？
	// 站在业务的角度，你是不是至少应该发成功一次
	err1 := n.producer.ProducePaymentEvent(ctx, evt)
	if err1 != nil {
		// 失败的时候，我们并没有将本地消息表标记为失败，是因为我们后面还想继续重试
		// msg 表里面有一条处于 StatusInit 的数据
		n.l.Error("发送支付事件失败", logger.Error(err1),
			logger.String("biz_trade_no", *txn.OutTradeNo))
		return nil
	}

	// 更新本地消息表状态
	// 这里我认为即便是本地消息表更新失败，也不是业务失败
	err1 = n.msgRepo.MarkSuccess(ctx, msgId)
	if err1 != nil {
		// 没有把 msg 标记为发送成功
		// 消息队列里面会有至少两条消息
		n.l.Error("将本地消息表标记为成功操作失败", logger.Error(err1),
			logger.String("biz_trade_no", *txn.OutTradeNo))
	}

	return nil
}