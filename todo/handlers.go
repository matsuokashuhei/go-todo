package todo

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type TodoHandler struct {
	repository *TodoRepository
}

func (handler *TodoHandler) GetAll(c *fiber.Ctx) error {
	var todos []Todo = handler.repository.FindAll()
	return c.JSON(todos)
}

func (handler *TodoHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	todo, err := handler.repository.Find(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": 404, "error": err})
	}
	return c.JSON(todo)
}

func (handler *TodoHandler) Create(c *fiber.Ctx) error {
	data := new(Todo)
	if err := c.BodyParser(data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"error":   err,
		})
	}
	item, err := handler.repository.Create(*data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  404,
			"message": "Failed creating item",
			"error":   err,
		})
	}
	return c.JSON(item)
}

func (handler *TodoHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Item not found",
			"error":   err,
		})
	}
	todo, err := handler.repository.Find(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Item not found",
		})
	}
	todoData := new(Todo)
	if err := c.BodyParser(todoData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err,
		})
	}
	todo.Name = todoData.Name
	todo.Description = todoData.Description
	todo.Status = todoData.Status
	item, err := handler.repository.Save(todo)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"messaage": "Error updating todo",
			"error":    err,
		})
	}
	return c.JSON(item)
}

func (handler *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Item not found",
			"error":   err,
		})
	}
	RowsAffected := handler.repository.Delete(id)
	statusCode := 204
	if RowsAffected == 0 {
		statusCode = 400
	}
	return c.Status(statusCode).JSON(nil)
}

func NewTodoHandler(repoistory *TodoRepository) *TodoHandler {
	return &TodoHandler{repository: repoistory}
}

func Register(router fiber.Router, database *gorm.DB) {
	database.AutoMigrate(&Todo{})
	todoRepository := NewTodoRepository(database)
	todoHandler := NewTodoHandler(todoRepository)

	todoRouter := router.Group("/todo")
	todoRouter.Get("/", todoHandler.GetAll)
	todoRouter.Get("/:id", todoHandler.Get)
	todoRouter.Put("/:id", todoHandler.Update)
	todoRouter.Post("/", todoHandler.Create)
	todoRouter.Delete("/:id", todoHandler.Delete)

}
