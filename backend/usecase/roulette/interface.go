package roulette

import "context"

type Usecase interface {
	// スピン実行
	Spin(ctx context.Context, in SpinInput) (*SpinOutput, error)

	// 履歴取得
	GetSpinHistory(ctx context.Context, in HistoryInput) (*HistoryOutput, error)
}
