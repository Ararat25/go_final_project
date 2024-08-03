package customError

import "errors"

var (
	ErrRepeatNotSpecified    = errors.New("параметр repeat не задан")
	ErrNotValidRepeat        = errors.New("параметр repeat имеет неверный формат")
	ErrTaskTitleNotSpecified = errors.New("не указан заголовок задачи")
	ErrNotValidID            = errors.New("указан неверный идентификатор")
)
