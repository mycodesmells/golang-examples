package actions_test

func (as *ActionSuite) Test_ResumeHandler() {
	res := as.HTML("/resume").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "This is resume page!")
}
