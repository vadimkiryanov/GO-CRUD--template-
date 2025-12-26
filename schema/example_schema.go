package schema

import "time"

type ExampleSchema struct {
	// json:"-" - означает, что поле будет пропущено при JSON сериализации
	// db:"id" указывает, что поле Id в структуре соответствует колонке с именем "id" в таблице базы данных
	Id        int       `json:"-" db:"id"`
	Name      string    `json:"name" binding:"required"` // Когда поле помечено как binding:"required",
	Email     string    `json:"email" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
