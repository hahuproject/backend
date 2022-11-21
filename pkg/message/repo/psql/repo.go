package message_psql_repo

import (
	"database/sql"
	"log"

	message_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/domain"
	message_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/repo"
)

type MessagePsqlRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewMessagePsqlRepoAdapter(log *log.Logger, db *sql.DB) (message_repo.MessageRepoPort, error) {
	err := db.Ping()
	if err != nil {
		return &MessagePsqlRepoAdapter{}, err
	}
	return &MessagePsqlRepoAdapter{log: log, db: db}, nil
}

func populateMessage(db *sql.DB, message *message_domain.Message) {

	//Populate sender
	err := db.QueryRow(`SELECT
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type, created_at,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.users
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.users.user_id = $1`, message.From.ID).Scan(
		&message.From.ID, &message.From.FirstName, &message.From.LastName, &message.From.Email, &message.From.Phone, &message.From.Username, &message.From.ProfilePic, &message.From.Verified, &message.From.Type, &message.From.CreatedAt,
		&message.From.Address.ID, &message.From.Address.Country, &message.From.Address.Region, &message.From.Address.City, &message.From.Address.SubCity, &message.From.Address.Woreda, &message.From.Address.HouseNo)

	if err != nil {
		log.Println("Error in populating sender")
		log.Println(err)
	}

	//Populate reciever
	err = db.QueryRow(`SELECT
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type, created_at,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.users
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.users.user_id = $1`, message.To.ID).Scan(
		&message.To.ID, &message.To.FirstName, &message.To.LastName, &message.To.Email, &message.To.Phone, &message.To.Username, &message.To.ProfilePic, &message.To.Verified, &message.To.Type, &message.From.CreatedAt,
		&message.To.Address.ID, &message.To.Address.Country, &message.To.Address.Region, &message.To.Address.City, &message.To.Address.SubCity, &message.To.Address.Woreda, &message.To.Address.HouseNo)

	if err != nil {
		log.Println("Error in populating reciever")
		log.Println(err)
	}
}

func (messageRepo MessagePsqlRepoAdapter) StoreMessage(message message_domain.Message) (message_domain.Message, error) {
	var addedMessage message_domain.Message
	messageRepo.log.Println(message.From.ID)
	messageRepo.log.Println(message.To.ID)
	messageRepo.log.Println(message.Content)
	err := messageRepo.db.QueryRow(`INSERT INTO public.messages ("from", "to" , content) VALUES ($1,$2,$3) RETURNING public.messages.message_id`, message.From.ID, message.To.ID, message.Content).Scan(&addedMessage.ID)
	if err != nil {
		return addedMessage, err
	}
	return messageRepo.FindMessageById(addedMessage.ID)
}

func (messageRepo MessagePsqlRepoAdapter) FindMessageById(messageId string) (message_domain.Message, error) {
	var message message_domain.Message

	err := messageRepo.db.QueryRow(`SELECT message_id, "from", "to", content, created_at, read FROM public.messages WHERE message_id = $1`, messageId).Scan(
		&message.ID, &message.From.ID, &message.To.ID, &message.Content, &message.CreatedAt, &message.Read,
	)
	if err != nil {
		return message, err
	}

	populateMessage(messageRepo.db, &message)

	return message, nil
}

func (messageRepo MessagePsqlRepoAdapter) FindUserMessages(userId string) ([]message_domain.Message, error) {
	var userMessages []message_domain.Message = make([]message_domain.Message, 0)

	rows, err := messageRepo.db.Query(`SELECT message_id, public.messages.from, public.messages.to, content, created_at, read FROM public.messages WHERE "from" = $1 OR "to" = $1 ORDER BY created_at`, userId)
	if err != nil {
		return userMessages, err
	}

	for rows.Next() {
		var message message_domain.Message
		rows.Scan(&message.ID, &message.From.ID, &message.To.ID, &message.Content, &message.CreatedAt, &message.Read)

		populateMessage(messageRepo.db, &message)

		userMessages = append(userMessages, message)
	}

	return userMessages, nil
}

func (messageRepo MessagePsqlRepoAdapter) ReadMessage(messageId string, userId string) (message_domain.Message, error) {
	var updatedMessage message_domain.Message

	err := messageRepo.db.QueryRow(
		`UPDATE 
		public.messages 
		SET read = $1
		WHERE message_id = $2 AND "to" = $3 
		RETURNING message_id, public.messages.from, public.messages.to, content, created_at, read `, true, messageId, userId).Scan(
		&updatedMessage.ID, &updatedMessage.From.ID, &updatedMessage.To.ID, &updatedMessage.Content, &updatedMessage.CreatedAt, &updatedMessage.Read,
	)
	if err != nil {
		return updatedMessage, err
	}

	populateMessage(messageRepo.db, &updatedMessage)

	return updatedMessage, nil
}
