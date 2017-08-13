package actions_test

func (as *ActionSuite) Test_ContactHandler() {
	res := as.HTML("/contact").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "This is contact page!")
}
