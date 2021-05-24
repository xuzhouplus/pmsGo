package gateway

type Weibo struct {
}

func NewWeibo() (*Weibo, error) {
	gateway := &Weibo{}

	return gateway, nil
}
func (gateway Weibo) AuthorizeUrl() {

}
