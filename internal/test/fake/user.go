package fake

type UserInput struct {
	Name     string
	Email    string
	Password string
}

func User() UserInput {
	return UserInput{
		Name:     "john",
		Email:    "john@test.com",
		Password: "abcABC123",
	}
}

func Admin() UserInput {
	return UserInput{
		Name:     "admin",
		Email:    "admin@test.com",
		Password: "abcABC123",
	}
}
