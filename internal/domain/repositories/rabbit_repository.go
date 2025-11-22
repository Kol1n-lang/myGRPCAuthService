package repositories

type RabbitRepository interface {
	CreateEmailMSG(email string) error
}
