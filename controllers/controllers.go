package controllers

import (
	"strconv"
	"time"

	"github.com/Epic55/go_project_jwt_oauth/db"
	"github.com/Epic55/go_project_jwt_oauth/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data1 map[string]string

	if err := c.BodyParser(&data1); err != nil {
		return err
	}
	pwd, _ := bcrypt.GenerateFromPassword([]byte(data1["1"]), 14)

	user1 := models.User{
		Name:     data1["name"],
		Email:    data1["email"],
		Password: pwd,
	}

	db.DB.Create(&user1) //Adds the data to the DB

	return c.JSON(data1)
}

const SecretKey = "secret"

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db.DB.Where("email = ?", data["email"]).First(&user) //Check the email is present in the DB

	if user.ID == 0 { //If the ID return is '0' then there is no such email present in the DB
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	} // If the email is present in the DB then compare the Passwords and if incorrect password then return error.

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),            //issuer contains the ID of the user.
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //Adds time to the token i.e. 24 hours.
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	} //Creates the cookie to be passed.

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
