package models

type Information struct {
	User_id     uint   `json:"-"                     gorm:"primaryKey;autoIncrement:true"`
	Email       string `json:"email"                 gorm:"column:email;type:varchar(35)"`
	Username    string `json:"username"              gorm:"column:user_name;type:varchar(100)"`
	Password    string `json:"password"              gorm:"column:password;type:varchar(100)"`
	PhoneNumber string `json:"phone_number"          gorm:"column:phone_number;type:varchar(15)"`
}

type TaskDetails struct {
	TASK_ID   uint   `json:"-"                   gorm:"primaryKey;autoIncrement:true"`
	// User_id   uint   `json:"-"                   gorm:"type:bigint references Information(user_id)"`
	TASK_NAME string `json:"task_name"           gorm:"column:task_name;type:varchar(50)"`
	Status    string `json:"status"              gorm:"column:status;type:varchar(50)"`
}
