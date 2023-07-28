package handler

import (
	//built in package

	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	//user defined package
	"todo/helper"
	logs "todo/log"
	"todo/models"
	"todo/repository"

	//third party package

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Database struct {
	Database *gorm.DB
}

// SignUp API
func (Db Database) Signup(c *fiber.Ctx) error {
	var user models.Information
	log := logs.Logs()
	log.Info("Signup api called successfully")

	if err := c.BodyParser(&user); err != nil {
		log.Error("error:'Invalid Format' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Invalid Format",
			"status": 400,
		})
	}
	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		log.Error("error:'Invalid Email Format' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Invalid Email Format",
			"status": 400,
		})
	}
	//make sure username field should not be empty
	if user.Username == "" {
		log.Error("error:'Username field should not be empty' status:400")
		return c.Status(http.StatusForbidden).JSON(map[string]interface{}{
			"Error":  "Username field should not be empty",
			"status": 403,
		})
	}
	//password should have minimum 8 character
	if len(user.Password) < 8 {
		log.Error("error:'Password should be more than 8 characters' status:400")
		return c.Status(http.StatusForbidden).JSON(map[string]interface{}{
			"Error":  "Password should be more than 8 characters",
			"status": 403,
		})

	}
	//passwords are stored in hashing method in the database
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil
	}
	user.Password = string(password)

	// Validate phone number
	phoneNumber := strings.TrimSpace(user.PhoneNumber)
	// Use regular expression to validate numeric characters and length
	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		log.Error("error:'Invalid phone number format' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Invalid phone number format",
			"status": 400,
		})
	}
	//checks the user already exist or not
	_, err = repository.ReadUserByEmail(Db.Database, user)
	if err == nil {
		log.Error("error:'user already exist' status:400")
		return c.Status(http.StatusForbidden).JSON(map[string]interface{}{
			"error":  "user already exist",
			"status": 403,
		})
	}
	repository.CreateUser(Db.Database, user)
	log.Info("message:'sign up successfull' status:200")
	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "sign up successfull",
		"status":  200,
	})
}

// Login API
func (Db Database) Login(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("login api called successfully")
	var login models.Information
	if err := c.BodyParser(&login); err != nil {
		log.Error("error:'Invalid Format' status:400")
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"error":  "Invalid Format",
			"status": 500,
		})
	}
	//verify the email whether its already registered in the SignUp API or not
	verify, err := repository.ReadUserByEmail(Db.Database, login)
	if err == nil {
		//checks whether the given password matches with the email
		if err := bcrypt.CompareHashAndPassword([]byte(verify.Password), []byte(login.Password)); err != nil {
			log.Error("error:'Password Not Matching' status:400")
			return c.Status(http.StatusForbidden).JSON(map[string]interface{}{
				"Error":  " Password Not Matching",
				"status": 403,
			})
		}
		userid := strconv.Itoa(int(verify.User_id))
		//generates token when email and password matches
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  userid,
			"email":    verify.Email,
			"password": login.Password,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
		err := helper.Configure(".env")
		if err != nil {
			fmt.Println("error is loading env file ")
		}
		secretkey := os.Getenv("SIGNINGKEY")
		tokenString, err := token.SignedString([]byte(secretkey))
		if err != nil {
			log.Error("error:'Failed To Generate Token' status:400")
			return c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
				"Error":  "Failed To Generate Token",
				"status": 401,
			})
		}
		log.Info("message:'Login Successful' status:200")
		return c.Status(http.StatusOK).JSON(map[string]interface{}{
			"message": "Login Successful",
			"token":   tokenString,
			"status":  200,
		})
	}
	log.Error("error:'login failed' status:400")
	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"Error":  "email not registered",
		"status": 400,
	})
}

// Task Posting API
func (Db Database) TaskPosting(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info(" TaskRemainder api called successfully")
	var post models.TaskDetails
	if err := c.BodyParser(&post); err != nil {
		log.Error("error:'invalid format' status:400")
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"Error":  "invalid format",
			"status": 500,
		})
	}

	tokenStr := c.GetReqHeaders()
	tokenString := tokenStr["Authorization"]
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing token")
	}
	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	userid, _ := strconv.Atoi(claims["user_id"].(string))
	post.User_id = uint(userid)

	if post.Status != "active" && post.Status != "completed" {
		log.Error("error:'' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Invalid value for status field.Only 'active' and 'completed' are allowed.",
			"status": 400,
		})
	}
	err := repository.TaskPosting(Db.Database, post)
	if err != nil {
		log.Error("error:'error in adding task details' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "error in adding task details",
			"status": 400,
		})
	}
	log.Info("error:'Task added Successfully' status:200")
	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "Task added Successfully",
		"status":  200,
	})
}

// update the task details by using task ID
func (Db Database) UpdateTask(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("UpdateTask API called successfully")

	taskID := c.Params("id")
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing token")
	}

	tokenString := ""
	for index, char := range tokenStr {
		if char == ' ' {
			tokenString = tokenStr[index+1:]
		}
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		log.Error("error: invalid token")
		return c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"error":  "Invalid token",
			"status": http.StatusUnauthorized,
		})
	}

	userID := claims["user_id"].(string)
	task, err := repository.GetTaskById(Db.Database, taskID, userID)
	if err != nil {
		log.Error("error: task not found or user does not have access")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "Task not found or user does not have access",
			"status": http.StatusNotFound,
		})
	}
	var updatedTask models.TaskDetails
	if err := c.BodyParser(&updatedTask); err != nil {
		log.Error("error: failed to parse request body")
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"error":  "Failed to parse request body",
			"status": http.StatusInternalServerError,
		})
	}

	// Update the task object with the
	task.TASK_NAME = updatedTask.TASK_NAME
	task.Status = updatedTask.Status

	//update operation
	err = repository.UpdateTask(Db.Database, task)
	if err != nil {
		log.Error("error: 'Task not found or user does not have access'")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "Task not found or user does not have access",
			"status": http.StatusNotFound,
		})
	}

	log.Info("Task updated successfully")
	return c.JSON(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Task updated successfully",
	})
}

// delete the task details by using task id
func (Db Database) DeleteTask(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("DeleteTask API called successfully")

	taskID := c.Params("id")
	tokenStr := c.GetReqHeaders()
	tokenString := tokenStr["Authorization"]
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing token")
	}
	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		log.Error("error: invalid token")
		return c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"error":  "Invalid token",
			"status": http.StatusUnauthorized,
		})
	}

	userID := claims["user_id"].(string)
	task, err := repository.GetTaskById(Db.Database, taskID, userID)
	if err != nil {
		log.Error("error: task not found or user does not have access")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "Task not found or user does not have access",
			"status": http.StatusNotFound,
		})
	}

	// Perform the delete operation when the user have the permission
	err = repository.DeleteTask(Db.Database, task)
	if err != nil {
		log.Error("error: failed to delete task")
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"error":  "Failed to delete task",
			"status": http.StatusInternalServerError,
		})
	}
	log.Error("error:task deleted")
	return c.JSON(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Task deleted successfully",
	})
}

// Get tasks by user and date
func (Db Database) GetTasksByUserAndDate(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("GetTasksByUserAndDate API called successfully")

	// Get the user ID from the JWT token
	tokenStr := c.GetReqHeaders()
	tokenString := tokenStr["Authorization"]
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).SendString("Missing token")
	}
	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	userID := claims["user_id"].(string)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	status := c.Query("status")
	t1, _ := time.Parse("02-01-2006", startDate)
	t2, _ := time.Parse("02-01-2006", endDate)

	if status == "" {
		// Get tasks by user and date
		tasks, err := repository.GetTasksByUserAndDate(Db.Database, userID, t1, t2)
		if err != nil {
			log.Error("error: failed to get tasks by user and date")
			return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
				"error":  "Failed to get tasks by user and date",
				"status": http.StatusInternalServerError,
			})
		}
		if len(tasks) == 0 {
			log.Info("No tasks assigned on the given date")
			return c.JSON(map[string]interface{}{
				"status":  http.StatusOK,
				"message": "No tasks assigned on the given date",
			})
		}
		return c.JSON(map[string]interface{}{
			"status": http.StatusOK,
			"tasks":  tasks,
		})
	} else if startDate == "" && endDate == ""{
		tasks, err := repository.GetTasksByStatus(Db.Database, userID, status)
		if err != nil {
			log.Error("error: no status ")
			return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
				"error":  "no status",
				"status": http.StatusInternalServerError,
			})
		}
		if len(tasks) == 0 {
			log.Info("no status found ")
			return c.JSON(map[string]interface{}{
				"status":  http.StatusOK,
				"message": "no status found ",
			})
		}
		return c.JSON(map[string]interface{}{
			"status": http.StatusOK,
			"tasks":  tasks,
		})
	}else{
		tasks, err := repository.GetTasksByUserAndDateAndStatus(Db.Database, userID, t1, t2, status)
		if err != nil {
			log.Error("error: failed to get tasks by user and date")
			return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
				"error":  "Failed to get tasks by user and date",
				"status": http.StatusInternalServerError,
			})
		}
		if len(tasks) == 0 {
			log.Info("No status found")
			return c.JSON(map[string]interface{}{
				"status":  http.StatusOK,
				"message": "No status found",
			})
		}
		return c.JSON(map[string]interface{}{
			"status": http.StatusOK,
			"tasks":  tasks,
		})
	}
}
