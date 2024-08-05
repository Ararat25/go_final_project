package customError

import "errors"

var (
	ErrRepeatNotSpecified      = errors.New("параметр repeat не задан")
	ErrNotValidRepeat          = errors.New("параметр repeat имеет неверный формат")
	ErrTaskTitleNotSpecified   = errors.New("не указан заголовок задачи")
	ErrNotValidID              = errors.New("указан неверный идентификатор")
	ErrIdNotSpecified          = errors.New("не указан идентификатор")
	ErrInvalidIdFormat         = errors.New("указан не верный формат идентификатора")
	ErrEnvPasswordNotSpecified = errors.New("в переменных среды не указана переменная пароля")
	ErrPasswordNotSpecified    = errors.New("пароль не указан")
	ErrPasswordNotValid        = errors.New("неверный пароль")
	ErrCannotParseJwtClaims    = errors.New("не удалось спрасить claims из jwt токена")
	ErrNotValidJwtToken        = errors.New("не валидный jwt токен")
)
