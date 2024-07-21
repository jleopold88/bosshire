// handlers/handlers.go
package handlers

import (
	"database/sql"
	"time"

	db "bosshire.com/db"
	"bosshire.com/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
	}
	user.Password = string(hash)

	_, err = db.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, user.Password, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
}

func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	user := models.User{}
	err := db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username=$1", input.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect username or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create token"})
	}

	return c.JSON(fiber.Map{"token": t})
}

func ViewJobs(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, title, description, requirements, employer_id FROM jobs")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch jobs"})
	}
	defer rows.Close()

	jobs := []models.Job{}
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Requirements, &job.EmployerID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot scan job"})
		}
		jobs = append(jobs, job)
	}

	return c.JSON(jobs)
}

func ProcessApplication(c *fiber.Ctx) error {
	applicationID := c.Params("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	_, err := db.DB.Exec("UPDATE applications SET status = $1 WHERE id = $2", input.Status, applicationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update application"})
	}

	return c.JSON(fiber.Map{"message": "application updated"})
}

func ApplyJob(c *fiber.Ctx) error {
	jobID := c.Params("id")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	talentID := int(claims["id"].(float64))

	if claims["role"] != "talent" {
		return c.SendStatus(fiber.ErrUnauthorized.Code)
	}

	_, err := db.DB.Exec("INSERT INTO applications (job_id, talent_id, status) VALUES ($1, $2, $3)", jobID, talentID, "applied")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot apply for job"})
	}

	return c.JSON(fiber.Map{"message": "applied for job"})
}

func PostJob(c *fiber.Ctx) error {
	var job models.Job
	if err := c.BodyParser(&job); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	employerID := int(claims["id"].(float64))

	_, err := db.DB.Exec("INSERT INTO jobs (title, description, requirements, employer_id) VALUES ($1, $2, $3, $4)", job.Title, job.Description, job.Requirements, employerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot post job"})
	}

	return c.JSON(fiber.Map{"message": "job posted"})
}

func ReviewApplications(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))
	role := claims["role"].(string)

	var rows *sql.Rows
	var err error

	if role == "employer" {
		rows, err = db.DB.Query("SELECT id, job_id, talent_id, status FROM applications WHERE job_id IN (SELECT id FROM jobs WHERE employer_id = $1)", userID)
	} else {
		rows, err = db.DB.Query("SELECT id, job_id, talent_id, status FROM applications WHERE talent_id = $1", userID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch applications"})
	}
	defer rows.Close()

	applications := []models.Application{}
	for rows.Next() {
		var application models.Application
		if err := rows.Scan(&application.ID, &application.JobID, &application.TalentID, &application.Status); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot scan application"})
		}
		applications = append(applications, application)
	}

	return c.JSON(applications)
}
