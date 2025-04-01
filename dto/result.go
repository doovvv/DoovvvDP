package dto
type Result struct {
	Success bool `json:"success"`
	ErrorMsg string `json:"errorMsg"`
	Data any `json:"data"`
	Total int `json:"total"`
}
func (r *Result) Ok() *Result{
	r.Success = true
	return r
}
func (r *Result) OkWithData(data any)*Result{
	r.Success = true
	r.Data = data
	return r
}
func (r *Result) Fail(msg string) *Result{
	r.Success = false
	r.ErrorMsg = msg
	return r
}