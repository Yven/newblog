package validate

type Admin struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
