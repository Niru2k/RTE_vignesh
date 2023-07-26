package handler

import (
	//built in package

	"net/http"
	"regexp"
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
)

// SignUp API
func Signup(c *fiber.Ctx) error {
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
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Username field should not be empty",
			"status": 400,
		})
	}
	//password should have minimum 8 character
	if len(user.Password) < 8 {
		log.Error("error:'Password should be more than 8 characters' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Password should be more than 8 characters",
			"status": 400,
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
	_, err = repository.ReadUserByEmail(user)
	if err == nil {
		log.Error("error:'user already exist' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"error":  "user already exist",
			"status": 400,
		})
	}
	repository.CreateUser(user)
	log.Info("message:'sign up successfull' status:200")
	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"message": "sign up successfull",
		"status":  200,
	})
}

// Login API
func Login(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("login api called successfully")
	var login models.Information
	if err := c.BodyParser(&login); err != nil {
		log.Error("error:'Invalid Format' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"error":  "Invalid Format",
			"status": 400,
		})
	}
	//verify the email whether its already registered in the SignUp API or not
	verify, err := repository.ReadUserByEmail(login)
	if err == nil {
		//checks whether the given password matches with the email
		if err := bcrypt.CompareHashAndPassword([]byte(verify.Password), []byte(login.Password)); err != nil {
			log.Error("error:'Password Not Matching' status:400")
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"Error":  " Password Not Matching",
				"status": 400,
			})
		}
		//generates token when email and password matches
		login.Email = verify.Email
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  login.User_id,
			"email":    login.Email,
			"password": login.Password,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(helper.SigningKey)
		if err != nil {
			log.Error("error:'Failed To Generate Token' status:400")
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"Error":  "Failed To Generate Token",
				"status": 400,
			})
		}
		log.Info("message:'Login Successful' status:200")
		return c.Status(http.StatusAccepted).JSON(map[string]interface{}{
			"message": "Login Successful",
			"token":   tokenString,
			"status":  200,
		})
	}
	log.Error("error:'login failed' status:400")
	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"Error":  "login failed",
		"status": 400,
	})
}

// Task Posting API
func TaskRemainder(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info(" TaskRemainder api called successfully")
	var post models.TaskDetails
	if err := c.BodyParser(&post); err != nil {
		log.Error("error:'invalid format' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "invalid format",
			"status": 400,
		})
	}

	if post.Status != "active" && post.Status != "completed" {
		log.Error("error:'' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "Invalid value for status field.Only 'active' and 'completed' are allowed.",
			"status": 400,
		})
	}
	// jobId := strconv.Itoa(int(post.User_id))
	// _,err:= repository.GetUserId(jobId)
	// if err != nil {
	// 	log.Error("Error:'User ID not found' status:404")
	// 	return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
	// 		"Error":  "User ID not found",
	// 		"status": 404,
	// 	})
	// }

	err := repository.TaskPosting(post)
	if err != nil {
		log.Error("error:'error in adding task details' status:400")
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"Error":  "error in adding task details",
			"status": 400,
		})
	}
	log.Info("error:'Task added Successfully' status:200")
	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"message": "Task added Successfully",
		"status":  200,
	})
}

// Get all task added details
func GetAllTaskDetails(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("GetAllTaskDetails API called successfully")
	creates, err := repository.GetAllTask()
	if err != nil {
		log.Error("error:'no record found' status:404")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "no record found",
			"status": 404,
		})
	}
	return c.JSON(map[string]interface{}{
		"status": fiber.StatusOK,
		"task":   creates,
	})
}

// get task detail by using task ID
func GetTaskDetailsByID(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("GetTaskDetailsByID API  called successfully")
	TaskID := c.Params("id")
	create, err := repository.GetTaskPostId(TaskID)
	if err != nil {
		log.Error("error:'Task details does not exist' status:404")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "task details does not exist",
			"status": 404,
		})
	}
	return c.JSON(map[string]interface{}{
		"status": fiber.StatusOK,
		"task":   create,
	})
}

// update task details by using task ID
func UpdateTask(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("UpdateTask API called Successfully")
	TaskID := c.Params("id")
	updatedjob, err := repository.ReadTaskPostById(TaskID)
	if err == nil {
		log.Error("error:'can't update status:400")
		if err := c.BodyParser(&updatedjob); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  "can't update status",
				"status": 400,
			})
		}

		err := repository.UpdateTask(TaskID, updatedjob)
		if err != nil {
			log.Error("error:'task id not found' status:404")
			return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
				"Error":  " task id not found",
				"status": 404,
			})
		}
		log.Info("message:'task updated successfully' status:200")
		return c.Status(http.StatusOK).JSON(map[string]interface{}{
			"message": "task updated successfully",
			"status":  200,
		})
	}
	log.Error("Error:'task post not found' status:404")
	return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
		"Error":  "task post not found",
		"status": 404,
	})
}

// Deletes the task details  using task id
func DeleteTask(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("Deleting task API called successfully")
	TaskID := c.Params("id")
	deletejob, err := repository.ReadTaskPostById(TaskID)
	if err == nil {

		repository.DeleteTaskPost(TaskID, deletejob)
		log.Info("message:'task details successfully deleted' status:200")
		return c.Status(http.StatusOK).JSON(map[string]interface{}{
			"message": " job post deleted successfully",
			"status":  200,
		})
	}
	log.Error("Error:'task id not found' status:404")
	return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
		"Error":  " task id not found",
		"status": 404,
	})
}

// get the active and completed task status
func GetTaskStatus(c *fiber.Ctx) error {
	log := logs.Logs()
	log.Info("GetTaskStatus API called successfully")
	company_jobs := c.Params("status")
	company_name, err := repository.GetTaskStatus(company_jobs)
	if err != nil || len(company_name) == 0 {
		log.Error("Error:'currently no status for this' status:404")
		return c.Status(http.StatusNotFound).JSON(map[string]interface{}{
			"error":  "currently no status for this",
			"status": 404,
		})
	}
	return c.JSON(map[string]interface{}{
		"status": fiber.StatusOK,
		"task":   company_name,
	})
}
