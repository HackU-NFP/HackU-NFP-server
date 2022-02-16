package ipresenter

// ILinePresenter LINEBOTプレゼンタ
type ILinePresenter interface {
	Parrot(token, msg string)
}
